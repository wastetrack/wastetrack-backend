package usecase

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/model/converter"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"gorm.io/gorm"
)

type WasteCategoryUsecase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validate                *validator.Validate
	WasteCategoryRepository *repository.WasteCategoryRepository
}

func NewWasteCategoryUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteCategoryRepository *repository.WasteCategoryRepository) *WasteCategoryUsecase {
	return &WasteCategoryUsecase{
		DB:                      db,
		Log:                     log,
		Validate:                validate,
		WasteCategoryRepository: wasteCategoryRepository,
	}
}

func (c *WasteCategoryUsecase) Create(ctx context.Context, request *model.WasteCategoryRequest) (*model.WasteCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		c.Log.Warnf("Failed to start transaction: %+v", tx.Error)
		return nil, fiber.ErrInternalServerError
	}
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	wasteCategory := &entity.WasteCategory{
		Name:        request.Name,
		Description: request.Description,
	}

	if err := c.WasteCategoryRepository.Create(tx, wasteCategory); err != nil {
		c.Log.Warnf("Failed to create waste category: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCategoryToResponse(wasteCategory), nil
}
func (c *WasteCategoryUsecase) Get(ctx context.Context, request *model.GetWasteCategoryRequest) (*model.WasteCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	wasteCategory := new(entity.WasteCategory)
	if err := c.WasteCategoryRepository.FindById(tx, wasteCategory, request.ID); err != nil {
		c.Log.Warnf("Failed find waste category by id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCategoryToResponse(wasteCategory), nil

}

func (c *WasteCategoryUsecase) Update(ctx context.Context, request *model.UpdateWasteCategoryRequest) (*model.WasteCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	wasteCategory := new(entity.WasteCategory)
	if err := c.WasteCategoryRepository.FindById(tx, wasteCategory, request.ID); err != nil {
		c.Log.Warnf("Failed find waste category by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		wasteCategory.Name = request.Name
	}
	if request.Description != "" {
		wasteCategory.Description = request.Description
	}

	if err := c.WasteCategoryRepository.Update(tx, wasteCategory); err != nil {
		c.Log.Warnf("Failed to update waste category: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCategoryToResponse(wasteCategory), nil
}

func (c *WasteCategoryUsecase) Search(ctx context.Context, request *model.SearchWasteCategoryRequest) ([]model.WasteCategoryResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warnf("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}
	wasteCategories, total, err := c.WasteCategoryRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste categories")
		return nil, 0, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteCategoryResponse, len(wasteCategories))
	for i, wasteCategory := range wasteCategories {
		responses[i] = *converter.WasteCategoryToResponse(&wasteCategory)
	}
	return responses, total, nil
}
func (c *WasteCategoryUsecase) Delete(ctx context.Context, request *model.DeleteWasteCategoryRequest) (*model.WasteCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find waste category by id
	wasteCategory := new(entity.WasteCategory)
	if err := c.WasteCategoryRepository.FindById(tx, wasteCategory, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete waste category
	if err := c.WasteCategoryRepository.Delete(tx, wasteCategory); err != nil {
		c.Log.Warnf("Failed delete waste category : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCategoryToResponse(wasteCategory), nil
}
