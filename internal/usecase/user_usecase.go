package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/helper"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/model/converter"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
	JWTHelper      *helper.JWTHelper
	EmailHelper    *helper.EmailHelper
	BaseURL        string
}

func NewUserUseCase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	userRepository *repository.UserRepository,
	jwtHelper *helper.JWTHelper,
	emailHelper *helper.EmailHelper,
	baseURL string,
) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		Validate:       validate,
		UserRepository: userRepository,
		JWTHelper:      jwtHelper,
		EmailHelper:    emailHelper,
		BaseURL:        baseURL,
	}
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
		return nil, fiber.ErrConflict
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

	user := &entity.User{
		Email:                  request.Email,
		Username:               request.Username,
		Password:               string(password),
		Role:                   request.Role,
		IsEmailVerified:        false,
		EmailVerificationToken: verificationToken,
	}

	if err := c.UserRepository.Create(tx, user); err != nil {
		c.Log.Warnf("Failed to create user to database: %v", err)
		return nil, fiber.ErrInternalServerError
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

	// Generate JWT tokens
	accessToken, err := c.JWTHelper.GenerateAccessToken(user.ID.String(), user.Role, user.IsEmailVerified)
	if err != nil {
		c.Log.Warnf("Failed to generate access token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshToken, err := c.JWTHelper.GenerateRefreshToken(user.ID.String())
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

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Generate new tokens with updated email verification status
	accessToken, err := c.JWTHelper.GenerateAccessToken(user.ID.String(), user.Role, user.IsEmailVerified)
	if err != nil {
		c.Log.Warnf("Failed to generate access token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshToken, err := c.JWTHelper.GenerateRefreshToken(user.ID.String())
	if err != nil {
		c.Log.Warnf("Failed to generate refresh token: %+v", err)
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
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Validate refresh token
	userID, err := c.JWTHelper.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		c.Log.Warnf("Invalid refresh token: %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	user := new(entity.User)
	if err := c.UserRepository.FindById(tx, user, userID); err != nil {
		c.Log.Warnf("Failed find user by id: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Generate new tokens
	accessToken, err := c.JWTHelper.GenerateAccessToken(user.ID.String(), user.Role, user.IsEmailVerified)
	if err != nil {
		c.Log.Warnf("Failed to generate access token: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	refreshToken, err := c.JWTHelper.GenerateRefreshToken(user.ID.String())
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
		// Don't reveal that email doesn't exist
		return nil
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

// func (c *UserUseCase) Search(ctx context.Context, request *model.SearchUserRequest) ([]model.UserResponse, int64, error) {
// 	tx := c.DB.WithContext(ctx).Begin()
// 	defer tx.Rollback()

// 	if err := c.Validate.Struct(request); err != nil {
// 		c.Log.Warnf("Invalid request body: %+v", err)
// 		return nil, 0, fiber.ErrBadRequest
// 	}

// 	if request.Page < 1 {
// 		request.Page = 1
// 	}
// 	if request.Size < 1 {
// 		request.Size = 10
// 	}

// 	users, total, err := c.UserRepository.Search(tx, request)
// 	if err != nil {
// 		c.Log.Warnf("Failed search users: %+v", err)
// 		return nil, 0, fiber.ErrInternalServerError
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		c.Log.Warnf("Failed commit transaction: %+v", err)
// 		return nil, 0, fiber.ErrInternalServerError
// 	}

// 	responses := make([]model.UserResponse, len(users))
// 	for i, user := range users {
// 		responses[i] = *converter.UserToResponse(&user)
// 	}

// 	return responses, total, nil
// }

// func (c *UserUseCase) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
// 	tx := c.DB.WithContext(ctx).Begin()
// 	defer tx.Rollback()

// 	if err := c.Validate.Struct(request); err != nil {
// 		c.Log.Warnf("Invalid request body: %+v", err)
// 		return nil, fiber.ErrBadRequest
// 	}

// 	user := new(entity.User)
// 	if err := c.UserRepository.FindById(tx, user, request.ID); err != nil {
// 		c.Log.Warnf("Failed find user by id: %+v", err)
// 		return nil, fiber.ErrNotFound
// 	}

// 	// Update fields if provided
// 	if request.Username != "" {
// 		user.Username = request.Username
// 	}
// 	if request.PhoneNumber != "" {
// 		user.PhoneNumber = request.PhoneNumber
// 	}
// 	if request.AvatarUrl != "" {
// 		user.AvatarUrl = request.AvatarUrl
// 	}
// 	if request.BirthDate != nil {
// 		user.BirthDate = *request.BirthDate
// 	}
// 	if request.GradeLevel != 0 {
// 		user.GradeLevel = request.GradeLevel
// 	}

// 	if err := c.UserRepository.Update(tx, user); err != nil {
// 		c.Log.Warnf("Failed update user: %+v", err)
// 		return nil, fiber.ErrInternalServerError
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		c.Log.Warnf("Failed commit transaction: %+v", err)
// 		return nil, fiber.ErrInternalServerError
// 	}

// 	return converter.UserToResponse(user), nil
// }

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
	// In a stateless JWT system, logout is typically handled client-side
	// by removing the tokens. However, you could implement token blacklisting
	// if needed for additional security.

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return false, fiber.ErrBadRequest
	}

	// For now, we'll just return true as the client should remove the tokens
	// In a production system, you might want to:
	// 1. Add tokens to a blacklist
	// 2. Store logout events
	// 3. Invalidate refresh tokens in database

	return true, nil
}
