package usecase

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/model/converter"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"gorm.io/gorm"
)

type PointConversionUsecase struct {
	DB                        *gorm.DB
	Log                       *logrus.Logger
	Validate                  *validator.Validate
	PointConversionRepository *repository.PointConversionRepository
	UserRepository            *repository.UserRepository
}

func NewPointConversionUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, pointConversionRepository *repository.PointConversionRepository, userRepository *repository.UserRepository) *PointConversionUsecase {
	return &PointConversionUsecase{
		DB:                        db,
		Log:                       log,
		Validate:                  validate,
		PointConversionRepository: pointConversionRepository,
		UserRepository:            userRepository,
	}
}

func (u *PointConversionUsecase) Create(ctx context.Context, request *model.PointConversionRequest) (*model.PointConversionSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Parse UUID
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		u.Log.Warnf("Invalid user ID: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if user exists
	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.UserID); err != nil {
		u.Log.Warnf("User not found: %v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if user.Points < request.Amount {
		u.Log.Warnf("User points are not enough: %v", err)
		return nil, fiber.NewError(fiber.StatusPaymentRequired, "User points are not enough")
	}

	pointConversion := &entity.PointConversion{
		UserID: userID,
		Amount: request.Amount,
		Status: request.Status,
	}

	if err := u.PointConversionRepository.Create(tx, pointConversion); err != nil {
		u.Log.Warnf("Failed to create point conversion: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PointConversionToSimpleResponse(pointConversion), nil
}

func (u *PointConversionUsecase) Get(ctx context.Context, id string) (*model.PointConversionResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	pointConversion := new(entity.PointConversion)
	if err := u.PointConversionRepository.FindByIdWithRelations(tx, pointConversion, id); err != nil {
		u.Log.Warnf("Point conversion not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PointConversionToResponse(pointConversion), nil
}

func (u *PointConversionUsecase) Update(ctx context.Context, request *model.UpdatePointConversionRequest) (*model.PointConversionSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	pointConversion := new(entity.PointConversion)
	if err := u.PointConversionRepository.FindById(tx, pointConversion, request.ID); err != nil {
		u.Log.Warnf("Point conversion not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	// Update UserID if provided
	if request.UserID != "" {
		userID, err := uuid.Parse(request.UserID)
		if err != nil {
			u.Log.Warnf("Invalid user ID: %v", err)
			return nil, fiber.ErrBadRequest
		}

		// Check if user exists
		user := new(entity.User)
		if err := u.UserRepository.FindById(tx, user, request.UserID); err != nil {
			u.Log.Warnf("User not found: %v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		pointConversion.UserID = userID
	}

	// Update status if provided
	if request.Status != "" {
		pointConversion.Status = request.Status
	}

	// Update IsDeleted if provided
	if request.IsDeleted != nil {
		pointConversion.IsDeleted = *request.IsDeleted
	}

	if err := u.PointConversionRepository.Update(tx, pointConversion); err != nil {
		u.Log.Warnf("Update failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PointConversionToSimpleResponse(pointConversion), nil
}

func (u *PointConversionUsecase) Search(ctx context.Context, request *model.SearchPointConversionRequest) ([]model.PointConversionSimpleResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	// Use search without preloads for simple response
	pointConversions, total, err := u.PointConversionRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Warn("Search failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Commit failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.PointConversionSimpleResponse, len(pointConversions))
	for i, pc := range pointConversions {
		responses[i] = *converter.PointConversionToSimpleResponse(&pc)
	}

	return responses, total, nil
}

func (u *PointConversionUsecase) Delete(ctx context.Context, id string) (*model.PointConversionSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	pointConversion := new(entity.PointConversion)
	if err := u.PointConversionRepository.FindById(tx, pointConversion, id); err != nil {
		u.Log.Warnf("Point conversion not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := u.PointConversionRepository.SoftDelete(tx, pointConversion); err != nil {
		u.Log.Warnf("Delete failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.PointConversionToSimpleResponse(pointConversion), nil
}
