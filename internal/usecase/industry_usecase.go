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

type IndustryUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	IndustryRepository *repository.IndustryRepository
}

func NewIndustryUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, industryRepository *repository.IndustryRepository) *IndustryUseCase {
	return &IndustryUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		IndustryRepository: industryRepository,
	}
}

// TODO: Implement Search

func (c *IndustryUseCase) Get(ctx context.Context, request *model.GetIndustryRequest) (*model.IndustryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	industry := new(entity.IndustryProfile)
	if err := c.IndustryRepository.FindByUserID(tx, industry, request.ID); err != nil {
		c.Log.Warnf("Failed find profile by user id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.IndustryToResponse(industry), nil

}

func (c *IndustryUseCase) Update(ctx context.Context, request *model.UpdateIndustryRequest, userID string, userRole string) (*model.IndustryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	industry := new(entity.IndustryProfile)
	if err := c.IndustryRepository.FindById(tx, industry, request.ID); err != nil {
		c.Log.Warnf("Failed find profile by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Authorization check: Skip for admin, otherwise check ownership
	if userRole != "admin" && industry.UserID != uuid.MustParse(userID) {
		c.Log.Warnf("User %s is not authorized to update profile %s", userID, request.ID)
		return nil, fiber.ErrForbidden
	}

	if request.TotalRecycledWeight != nil {
		industry.TotalRecycledWeight = *request.TotalRecycledWeight
	}

	if request.TotalWasteWeight != nil {
		industry.TotalWasteWeight = *request.TotalWasteWeight
	}

	if err := c.IndustryRepository.Update(tx, industry); err != nil {
		c.Log.Warnf("Failed to update profile: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.IndustryToResponse(industry), nil
}

func (c *IndustryUseCase) Delete(ctx context.Context, request *model.DeleteIndustryRequest) (*model.IndustryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find profile by id
	industry := new(entity.IndustryProfile)
	if err := c.IndustryRepository.FindById(tx, industry, request.ID); err != nil {
		c.Log.Warnf("Failed find profile by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete profile
	if err := c.IndustryRepository.Delete(tx, industry); err != nil {
		c.Log.Warnf("Failed delete profile : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.IndustryToResponse(industry), nil
}
