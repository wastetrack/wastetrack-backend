package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/model/converter"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"github.com/wastetrack/wastetrack-backend/internal/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB                            *gorm.DB
	Log                           *logrus.Logger
	Validate                      *validator.Validate
	UserRepository                *repository.UserRepository
	CustomerRepository            *repository.CustomerRepository
	WasteBankRepository           *repository.WasteBankRepository
	WasteCollectorRepository      *repository.WasteCollectorRepository
	IndustryRepository            *repository.IndustryRepository
	CollectorManagementRepository *repository.CollectorManagementRepository
	StorageRepository             *repository.StorageRepository
	JWTHelper                     *helper.JWTHelper
	EmailHelper                   *helper.EmailHelper
	BaseURL                       string
}

func NewUserUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	userRepository *repository.UserRepository,
	customerRepository *repository.CustomerRepository,
	wasteBankRepository *repository.WasteBankRepository,
	wasteCollectorRepository *repository.WasteCollectorRepository,
	industryRepository *repository.IndustryRepository,
	collectorManagementRepository *repository.CollectorManagementRepository,
	storageRepository *repository.StorageRepository,
	jwtHelper *helper.JWTHelper,
	emailHelper *helper.EmailHelper,
	baseURL string,
) *UserUseCase {
	return &UserUseCase{
		DB:                            db,
		Log:                           log,
		Validate:                      validate,
		UserRepository:                userRepository,
		CustomerRepository:            customerRepository,
		WasteBankRepository:           wasteBankRepository,
		WasteCollectorRepository:      wasteCollectorRepository,
		IndustryRepository:            industryRepository,
		CollectorManagementRepository: collectorManagementRepository,
		StorageRepository:             storageRepository,
		JWTHelper:                     jwtHelper,
		EmailHelper:                   emailHelper,
		BaseURL:                       baseURL,
	}
}

// TODO: Create Government profile upon registering
func getIsAcceptingCustomer(ptr *bool) bool {
	if ptr == nil {
		return true
	}
	return *ptr
}

func (c *UserUseCase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("validation error: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if email already exists
	total, err := c.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		c.Log.Warnf("Failed to count by email: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	if total > 0 {
		c.Log.Warnf("email already exist")
		return nil, fiber.NewError(fiber.StatusConflict, "email already exist")
	}

	// Hash password
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to hash password: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate email verification token
	verificationToken, err := c.EmailHelper.GenerateVerificationToken()
	if err != nil {
		c.Log.Warnf("Failed to generate verification token: %v", err)
		return nil, fiber.ErrInternalServerError
	}
	var location *types.Point
	if request.Location != nil {
		location = &types.Point{
			Lat: request.Location.Latitude,
			Lng: request.Location.Longitude,
		}
	}
	if !request.IsAgreedToTerms {
		return nil, fiber.NewError(fiber.StatusBadRequest, "You must agree to the terms and conditions")
	}
	isAcceptingCustomer := getIsAcceptingCustomer(request.IsAcceptingCustomer)
	user := &entity.User{
		Username:               request.Username,
		Email:                  request.Email,
		Password:               string(password),
		Role:                   request.Role,
		PhoneNumber:            request.PhoneNumber,
		Institution:            request.Institution,
		Address:                request.Address,
		City:                   request.City,
		Province:               request.Province,
		IsEmailVerified:        false,
		IsAcceptingCustomer:    isAcceptingCustomer,
		EmailVerificationToken: verificationToken,
		Location:               location,
	}
	c.Log.Infof("Creating user with IsAcceptingCustomer: %v", user.IsAcceptingCustomer)
	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed to create user to database: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Check user role
	if user.Role == "customer" {
		// Create customer profile
		customer := &entity.CustomerProfile{
			UserID: user.ID,
		}

		if err := c.CustomerRepository.Create(tx, customer); err != nil {
			c.Log.Warnf("Failed to create customer profile: %v", err)
			return nil, fiber.ErrInternalServerError
		}
	}
	if user.Role == "waste_bank_unit" || user.Role == "waste_bank_central" {
		// Create waste bank profile
		wasteBank := &entity.WasteBankProfile{
			UserID: user.ID,
		}

		if err := c.WasteBankRepository.Create(tx, wasteBank); err != nil {
			c.Log.Warnf("Failed to create waste bank profile: %v", err)
			return nil, fiber.ErrInternalServerError
		}
		storage := &entity.Storage{
			UserID: user.ID,
			Length: 0,
			Width:  0,
			Height: 0,
		}

		if err := c.StorageRepository.Create(tx, storage); err != nil {
			c.Log.Warnf("Failed to create storage: %v", err)
			return nil, fiber.ErrInternalServerError
		}
	}
	if user.Role == "waste_collector_unit" || user.Role == "waste_collector_central" {
		// Create waste collector profile
		wasteCollector := &entity.WasteCollectorProfile{
			UserID: user.ID,
		}

		// Only create collector management if InstitutionID is provided
		if request.InstitutionID != "" {
			// Validate that the institution exists
			wasteBank := new(entity.User)
			if err := c.UserRepository.FindById(tx, wasteBank, request.InstitutionID); err != nil {
				c.Log.Warnf("Failed to find waste bank by id: %v", err)
				return nil, fiber.NewError(fiber.StatusBadRequest, "Institution not found")
			}

			collectorManagement := &entity.CollectorManagement{
				WasteBankID: uuid.MustParse(request.InstitutionID),
				CollectorID: user.ID,
				Status:      "active",
			}
			if err := c.CollectorManagementRepository.Create(tx, collectorManagement); err != nil {
				c.Log.Warnf("Failed to create collector management: %v", err)
				return nil, fiber.ErrInternalServerError
			}
		}

		if err := c.WasteCollectorRepository.Create(tx, wasteCollector); err != nil {
			c.Log.Warnf("Failed to create waste collector profile: %v", err)
			return nil, fiber.ErrInternalServerError
		}
	}
	if user.Role == "industry" {
		// Create industry profile
		industry := &entity.IndustryProfile{
			UserID: user.ID,
		}

		if err := c.IndustryRepository.Create(tx, industry); err != nil {
			c.Log.Warnf("Failed to create industry profile: %v", err)
			return nil, fiber.ErrInternalServerError
		}
		storage := &entity.Storage{
			UserID:                user.ID,
			Length:                0,
			Width:                 0,
			Height:                0,
			IsForRecycledMaterial: false,
		}

		if err := c.StorageRepository.Create(tx, storage); err != nil {
			c.Log.Warnf("Failed to create storage: %v", err)
			return nil, fiber.ErrInternalServerError
		}
		recycleStorage := &entity.Storage{
			UserID:                user.ID,
			Length:                0,
			Width:                 0,
			Height:                0,
			IsForRecycledMaterial: true,
		}

		if err := c.StorageRepository.Create(tx, recycleStorage); err != nil {
			c.Log.Warnf("Failed to create storage: %v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Send verification email
	go func() {
		if err := c.EmailHelper.SendVerificationEmail(user.Email, user.Username, verificationToken, c.BaseURL); err != nil {
			c.Log.Errorf("Failed to send verification email: %v", err)
		}
	}()

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		c.Log.Warnf("Failed compare password: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	// Enforce session limit by automatically revoking oldest tokens if needed
	if err := c.JWTHelper.EnforceSessionLimit(tx, user.ID, 5); err != nil {
		c.Log.Warnf("Failed to enforce session limit: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate JWT tokens
	accessToken, err := c.JWTHelper.GenerateAccessToken(user.ID.String(), user.Role, user.IsEmailVerified)
	if err != nil {
		c.Log.Warnf("Failed to generate access token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate and store refresh token
	refreshToken, err := c.JWTHelper.GenerateRefreshToken(tx, user.ID)
	if err != nil {
		c.Log.Warnf("Failed to generate refresh token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := converter.UserToResponse(user)
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	return response, nil
}

func (c *UserUseCase) VerifyEmail(ctx context.Context, request *model.VerifyEmailRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmailVerificationToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("Failed find user by verification token: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if user.IsEmailVerified {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Email already verified")
	}

	user.IsEmailVerified = true
	user.EmailVerificationToken = ""

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed update user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate new tokens with updated email verification status
	accessToken, err := c.JWTHelper.GenerateAccessToken(user.ID.String(), user.Role, user.IsEmailVerified)
	if err != nil {
		c.Log.Warnf("Failed to generate access token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshToken, err := c.JWTHelper.GenerateRefreshToken(tx, user.ID)
	if err != nil {
		c.Log.Warnf("Failed to generate refresh token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := converter.UserToResponse(user)
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken

	return response, nil
}

func (c *UserUseCase) ResendVerification(ctx context.Context, request *model.ResendVerificationRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email: %+v", err)
		return fiber.ErrNotFound
	}

	if user.IsEmailVerified {
		return fiber.NewError(fiber.StatusBadRequest, "Email already verified")
	}

	// Generate new verification token
	verificationToken, err := c.EmailHelper.GenerateVerificationToken()
	if err != nil {
		c.Log.Warnf("Failed to generate verification token: %v", err)
		return fiber.ErrInternalServerError
	}

	user.EmailVerificationToken = verificationToken

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed update user: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return fiber.ErrInternalServerError
	}

	// Send verification email
	go func() {
		if err := c.EmailHelper.SendVerificationEmail(user.Email, user.Username, verificationToken, c.BaseURL); err != nil {
			c.Log.Errorf("Failed to send verification email: %v", err)
		}
	}()

	return nil
}

func (c *UserUseCase) RefreshToken(ctx context.Context, request *model.RefreshTokenRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Validate refresh token from database
	refreshToken, err := c.JWTHelper.ValidateRefreshToken(tx, request.RefreshToken)
	if err != nil {
		c.Log.Warnf("Invalid refresh token: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, refreshToken.UserID.String()); err != nil {
		c.Log.Warnf("Failed find user by id: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Revoke the used refresh token
	if err := c.JWTHelper.RevokeRefreshToken(tx, request.RefreshToken); err != nil {
		c.Log.Warnf("Failed to revoke refresh token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate new tokens
	accessToken, err := c.JWTHelper.GenerateAccessToken(user.ID.String(), user.Role, user.IsEmailVerified)
	if err != nil {
		c.Log.Warnf("Failed to generate access token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newRefreshToken, err := c.JWTHelper.GenerateRefreshToken(tx, user.ID)
	if err != nil {
		c.Log.Warnf("Failed to generate refresh token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	response := converter.UserToResponse(user)
	response.AccessToken = accessToken
	response.RefreshToken = newRefreshToken

	return response, nil
}

func (c *UserUseCase) ForgotPassword(ctx context.Context, request *model.ForgotPasswordRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByEmail(tx, user, request.Email); err != nil {
		c.Log.Warnf("Failed find user by email: %+v", err)

		return fiber.NewError(fiber.StatusNotFound, "Email not found")
	}

	// Generate reset token
	resetToken, err := c.EmailHelper.GenerateVerificationToken()
	if err != nil {
		c.Log.Warnf("Failed to generate reset token: %v", err)
		return fiber.ErrInternalServerError
	}

	// Set token and expiry (1 hour)
	expiry := time.Now().Add(time.Hour)
	user.ResetPasswordToken = resetToken
	user.ResetPasswordExpiry = &expiry

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed update user: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return fiber.ErrInternalServerError
	}

	// Send reset email
	go func() {
		if err := c.EmailHelper.SendPasswordResetEmail(user.Email, user.Username, resetToken, c.BaseURL); err != nil {
			c.Log.Errorf("Failed to send reset email: %v", err)
		}
	}()

	return nil
}

func (c *UserUseCase) ResetPassword(ctx context.Context, request *model.ResetPasswordRequest) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindByResetPasswordToken(tx, user, request.Token); err != nil {
		c.Log.Warnf("Failed find user by reset token: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid or expired reset token")
	}

	// Hash new password
	password, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.Log.Warnf("Failed to hash password: %v", err)
		return fiber.ErrInternalServerError
	}

	// Update password and clear reset token
	user.Password = string(password)
	user.ResetPasswordToken = ""
	user.ResetPasswordExpiry = nil

	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed update user: %+v", err)
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

// Add missing methods for user management
func (c *UserUseCase) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	return c.Register(ctx, request)
}

func (c *UserUseCase) Current(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Get(ctx context.Context, request *model.GetUserRequest) (*model.UserResponse, error) {
	return c.Current(ctx, request)
}

func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find the user to be updated
	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed to find user by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Authorization check: ensure the requesting user is the owner
	if user.ID.String() != request.UserID {
		c.Log.Warnf("Unauthorized update attempt: user %s trying to update user %s", request.UserID, request.ID)
		return nil, fiber.ErrForbidden
	}

	// Update user fields (only update non-empty fields)
	if request.Username != "" {
		user.Username = request.Username
	}
	if request.PhoneNumber != "" {
		user.PhoneNumber = request.PhoneNumber
	}
	if request.Address != "" {
		user.Address = request.Address
	}
	if request.City != "" {
		user.City = request.City
	}
	if request.Province != "" {
		user.Province = request.Province
	}
	if request.IsAcceptingCustomer != nil {
		user.IsAcceptingCustomer = *request.IsAcceptingCustomer
	}

	// Update the user in database
	if err := c.UserRepository.Update(tx, user); err != nil {
		c.Log.Warnf("Failed to update user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Convert to response and return
	response := converter.UserToResponse(user)
	return response, nil
}

func (c *UserUseCase) Search(ctx context.Context, request *model.SearchUserRequest) ([]model.UserListResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, 0, fiber.ErrBadRequest
	}

	if request.Page < 1 {
		request.Page = 1
	}
	if request.Size < 1 {
		request.Size = 10
	}

	users, customerProfiles, wasteBankProfiles, industryProfiles, governmentProfiles, total, err := c.UserRepository.Search(tx, request)
	if err != nil {
		c.Log.Warnf("Failed search users with profiles: %+v", err)
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.UserListResponse, len(users))
	for i, user := range users {
		response := converter.UserToListResponse(&user)
		userID := user.ID.String()

		// Attach profile data based on role
		switch user.Role {
		case "customer":
			if profile, exists := customerProfiles[userID]; exists {
				response.CustomerProfile = converter.CustomerToResponse(profile)
			}
		case "waste_bank_unit", "waste_bank_central":
			if profile, exists := wasteBankProfiles[userID]; exists {
				response.WasteBankProfile = converter.WasteBankToResponse(profile)
			}
		case "industry":
			if profile, exists := industryProfiles[userID]; exists {
				response.IndustryProfile = converter.IndustryToResponse(profile)
			}
		case "government":
			if profile, exists := governmentProfiles[userID]; exists {
				response.GovernmentProfile = converter.GovernmentToResponse(profile)
			}
			// Add other role cases as needed (waste_collector_unit, waste_collector_central, etc.)
		}

		responses[i] = *response
	}

	return responses, total, nil
}

func (c *UserUseCase) Delete(ctx context.Context, request *model.DeleteUserRequest) (*model.UserResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
		c.Log.Warnf("Failed find user by id: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := c.UserRepository.Delete(tx, user); err != nil {
		c.Log.Warnf("Failed delete user: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToResponse(user), nil
}

func (c *UserUseCase) Logout(ctx context.Context, request *model.LogoutUserRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return false, fiber.ErrBadRequest
	}

	// If refresh token provided, revoke it
	if request.RefreshToken != "" {
		if err := c.JWTHelper.RevokeRefreshToken(tx, request.RefreshToken); err != nil {
			c.Log.Warnf("Failed to revoke refresh token: %+v", err)
			// Don't fail logout if token doesn't exist
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func (c *UserUseCase) LogoutAllDevices(ctx context.Context, userID string) error {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fiber.ErrBadRequest
	}

	if err := c.JWTHelper.RevokeAllUserTokens(tx, userUUID); err != nil {
		c.Log.Warnf("Failed to revoke all user tokens: %+v", err)
		return fiber.ErrInternalServerError
	}

	return tx.Commit().Error
}
