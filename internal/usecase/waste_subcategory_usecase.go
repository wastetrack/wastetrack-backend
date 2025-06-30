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

type WasteSubCategoryUsecase struct {
	DB                         *gorm.DB
	Log                        *logrus.Logger
	Validate                   *validator.Validate
	WasteCategoryRepository    *repository.WasteCategoryRepository
	WasteSubCategoryRepository *repository.WasteSubCategoryRepository
}

func NewWasteSubCategoryUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteCategoryRepository *repository.WasteCategoryRepository, repo *repository.WasteSubCategoryRepository) *WasteSubCategoryUsecase {
	return &WasteSubCategoryUsecase{
		DB:                         db,
		Log:                        log,
		Validate:                   validate,
		WasteCategoryRepository:    wasteCategoryRepository,
		WasteSubCategoryRepository: repo,
	}
}

func (c *WasteSubCategoryUsecase) Create(ctx context.Context, request *model.WasteSubCategoryRequest) (*model.WasteSubCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}
	// check if category exists
	category := new(entity.WasteCategory)
	if err := c.WasteCategoryRepository.FindById(tx, category, request.CategoryID); err != nil {
		c.Log.Warnf("Failed to find waste category by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	wasteSubCategory := &entity.WasteSubcategory{
		Name:        request.Name,
		Description: request.Description,
		CategoryID:  category.ID,
	}

	if err := c.WasteSubCategoryRepository.Create(tx, wasteSubCategory); err != nil {
		c.Log.Warnf("Failed to create waste subcategory: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteSubCategoryToResponse(wasteSubCategory), nil
}

func (c *WasteSubCategoryUsecase) Get(ctx context.Context, request *model.GetWasteCategoryRequest) (*model.WasteSubCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteSubCategory := new(entity.WasteSubcategory)
	if err := c.WasteSubCategoryRepository.FindById(tx, wasteSubCategory, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste subcategory by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteSubCategoryToResponse(wasteSubCategory), nil
}

func (c *WasteSubCategoryUsecase) Update(ctx context.Context, request *model.UpdateWasteSubCategoryRequest) (*model.WasteSubCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteSubCategory := new(entity.WasteSubcategory)
	if err := c.WasteSubCategoryRepository.FindById(tx, wasteSubCategory, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste subcategory by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.Name != "" {
		wasteSubCategory.Name = request.Name
	}
	if request.Description != "" {
		wasteSubCategory.Description = request.Description
	}

	if err := c.WasteSubCategoryRepository.Update(tx, wasteSubCategory); err != nil {
		c.Log.Warnf("Failed to update waste subcategory: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteSubCategoryToResponse(wasteSubCategory), nil
}

func (c *WasteSubCategoryUsecase) Search(ctx context.Context, request *model.SearchWasteSubCategoryRequest) ([]model.WasteSubCategoryResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	subCategories, total, err := c.WasteSubCategoryRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to search waste subcategories")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteSubCategoryResponse, len(subCategories))
	for i, sub := range subCategories {
		responses[i] = *converter.WasteSubCategoryToResponse(&sub)
	}

	return responses, total, nil
}

func (c *WasteSubCategoryUsecase) Delete(ctx context.Context, request *model.DeleteWasteSubCategoryRequest) (*model.WasteSubCategoryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteSubCategory := new(entity.WasteSubcategory)
	if err := c.WasteSubCategoryRepository.FindById(tx, wasteSubCategory, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste subcategory by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := c.WasteSubCategoryRepository.Delete(tx, wasteSubCategory); err != nil {
		c.Log.Warnf("Failed to delete waste subcategory: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteSubCategoryToResponse(wasteSubCategory), nil
}
