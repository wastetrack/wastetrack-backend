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

type WasteSubCategoryUsecase struct {
	DB                         *gorm.DB
	Log                        *logrus.Logger
	Validate                   *validator.Validate
	WasteSubCategoryRepository *repository.WasteSubCategoryRepository
}

func NewWasteSubCategoryUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteSubCategoryRepository *repository.WasteSubCategoryRepository) *WasteSubCategoryUsecase {
	return &WasteSubCategoryUsecase{
		DB:                         db,
		Log:                        log,
		Validate:                   validate,
		WasteSubCategoryRepository: wasteSubCategoryRepository,
	}
}

func (c *WasteSubCategoryUsecase) Create(ctx context.Context, request *model.WasteSubCategoryRequest) (*model.WasteSubCategoryResponse, error) {
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
	wasteSubCategory := &entity.WasteSubcategory{
		Name:        request.Name,
		Description: request.Description,
		CategoryID:  uuid.MustParse(request.CategoryID),
	}

	if err := c.WasteSubCategoryRepository.Create(tx, wasteSubCategory); err != nil {
		c.Log.Warnf("Failed to create waste category: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCategoryToResponse(wasteSubCategory), nil
}
func (c *WasteSubCategoryUsecase) Get(ctx context.Context, request *model.GetWasteCategoryRequest) (*model.WasteSubCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	wasteCategory := new(entity.WasteCategory)
	if err := c.WasteSubCategoryRepository.FindById(tx, wasteCategory, request.ID); err != nil {
		c.Log.Warnf("Failed find waste category by id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCategoryToResponse(wasteCategory), nil

}

func (c *WasteSubCategoryUsecase) Update(ctx context.Context, request *model.UpdateWasteCategoryRequest) (*model.WasteSubCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	wasteCategory := new(entity.WasteCategory)
	if err := c.WasteSubCategoryRepository.FindById(tx, wasteCategory, request.ID); err != nil {
		c.Log.Warnf("Failed find waste category by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		wasteCategory.Name = request.Name
	}
	if request.Description != "" {
		wasteCategory.Description = request.Description
	}

	if err := c.WasteSubCategoryRepository.Update(tx, wasteCategory); err != nil {
		c.Log.Warnf("Failed to update waste category: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCategoryToResponse(wasteCategory), nil
}

func (c *WasteSubCategoryUsecase) Search(ctx context.Context, request *model.SearchWasteCategoryRequest) ([]model.WasteSubCategoryResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warnf("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}
	wasteCategories, total, err := c.WasteSubCategoryRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste categories")
		return nil, 0, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteSubCategoryResponse, len(wasteCategories))
	for i, wasteCategory := range wasteCategories {
		responses[i] = *converter.WasteCategoryToResponse(&wasteCategory)
	}
	return responses, total, nil
}
func (c *WasteSubCategoryUsecase) Delete(ctx context.Context, request *model.DeleteWasteCategoryRequest) (*model.WasteSubCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find waste category by id
	wasteCategory := new(entity.WasteCategory)
	if err := c.WasteSubCategoryRepository.FindById(tx, wasteCategory, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete waste category
	if err := c.WasteSubCategoryRepository.Delete(tx, wasteCategory); err != nil {
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
