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

type WasteTypeUsecase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validate                *validator.Validate
	WasteCategoryRepository *repository.WasteCategoryRepository
	WasteTypeRepository     *repository.WasteTypeRepository
}

func NewWasteTypeUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteCategoryRepository *repository.WasteCategoryRepository, wasteTypeRepository *repository.WasteTypeRepository) *WasteTypeUsecase {
	return &WasteTypeUsecase{
		DB:                      db,
		Log:                     log,
		Validate:                validate,
		WasteCategoryRepository: wasteCategoryRepository,
		WasteTypeRepository:     wasteTypeRepository,
	}
}

func (u *WasteTypeUsecase) Create(ctx context.Context, request *model.WasteTypeRequest) (*model.WasteTypeResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// check if category exists
	category := new(entity.WasteCategory)
	if err := u.WasteCategoryRepository.FindById(tx, category, request.CategoryID); err != nil {
		u.Log.Warnf("Category not found: %v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Waste Category not found")
	}

	wasteType := &entity.WasteType{
		Name:        request.Name,
		Description: request.Description,
		CategoryID:  category.ID,
	}

	if err := u.WasteTypeRepository.Create(tx, wasteType); err != nil {
		u.Log.Warnf("Failed to create waste type: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTypeToResponse(wasteType), nil
}

func (u *WasteTypeUsecase) Get(ctx context.Context, id string) (*model.WasteTypeResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	wasteType := new(entity.WasteType)
	if err := u.WasteTypeRepository.FindById(tx, wasteType, id); err != nil {
		u.Log.Warnf("Waste type not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTypeToResponse(wasteType), nil
}

func (u *WasteTypeUsecase) Update(ctx context.Context, request *model.UpdateWasteTypeRequest) (*model.WasteTypeResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteType := new(entity.WasteType)
	if err := u.WasteTypeRepository.FindById(tx, wasteType, request.ID); err != nil {
		u.Log.Warnf("Waste type not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	// Optional: check if new category ID exists (if you allow updating it)
	// if request.CategoryID != "" && request.CategoryID != wasteType.CategoryID.String() {
	//     category := new(entity.WasteCategory)
	//     if err := u.WasteCategoryRepository.FindById(tx, category, request.CategoryID); err != nil {
	//         u.Log.Warnf("New category not found: %v", err)
	//         return nil, fiber.ErrNotFound
	//     }
	//     wasteType.CategoryID = category.ID
	// }

	if request.Name != "" {
		wasteType.Name = request.Name
	}
	if request.Description != "" {
		wasteType.Description = request.Description
	}

	if err := u.WasteTypeRepository.Update(tx, wasteType); err != nil {
		u.Log.Warnf("Update failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTypeToResponse(wasteType), nil
}

func (u *WasteTypeUsecase) Search(ctx context.Context, request *model.SearchWasteTypeRequest) ([]model.WasteTypeResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	types, total, err := u.WasteTypeRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Warn("Search failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Commit failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteTypeResponse, len(types))
	for i, t := range types {
		responses[i] = *converter.WasteTypeToResponse(&t)
	}

	return responses, total, nil
}

func (u *WasteTypeUsecase) Delete(ctx context.Context, id string) (*model.WasteTypeResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	wasteType := new(entity.WasteType)
	if err := u.WasteTypeRepository.FindById(tx, wasteType, id); err != nil {
		u.Log.Warnf("Waste type not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := u.WasteTypeRepository.Delete(tx, wasteType); err != nil {
		u.Log.Warnf("Delete failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTypeToResponse(wasteType), nil
}
