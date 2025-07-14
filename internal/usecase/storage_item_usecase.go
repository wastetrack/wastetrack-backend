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

type StorageItemUsecase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	StorageRepository     *repository.StorageRepository
	StorageItemRepository *repository.StorageItemRepository
}

func NewStorageItemUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	storageRepo *repository.StorageRepository,
	storageItemRepo *repository.StorageItemRepository,
) *StorageItemUsecase {
	return &StorageItemUsecase{
		DB:                    db,
		Log:                   log,
		Validate:              validate,
		StorageRepository:     storageRepo,
		StorageItemRepository: storageItemRepo,
	}
}

func (u *StorageItemUsecase) Create(ctx context.Context, request *model.StorageItemRequest) (*model.StorageItemSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Verify storage exists
	storage := new(entity.Storage)
	if err := u.StorageRepository.FindById(tx, storage, request.StorageID); err != nil {
		u.Log.Warnf("Storage not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	storageItem := &entity.StorageItem{
		StorageID:   uuid.MustParse(request.StorageID),
		WasteTypeID: uuid.MustParse(request.WasteTypeID),
		WeightKgs:   request.WeightKgs,
	}

	if err := u.StorageItemRepository.Create(tx, storageItem); err != nil {
		u.Log.Warnf("Failed to create storage item: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StorageItemToSimpleResponse(storageItem), nil
}

func (u *StorageItemUsecase) Get(ctx context.Context, id string) (*model.StorageItemResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item := new(entity.StorageItem)
	if err := u.StorageItemRepository.FindById(tx, item, id); err != nil {
		u.Log.Warnf("Storage item not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StorageItemToResponse(item), nil
}

func (u *StorageItemUsecase) Update(ctx context.Context, request *model.UpdateStorageItemRequest) (*model.StorageItemSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}
	storage := new(entity.Storage)
	if err := u.StorageRepository.FindById(tx, storage, request.StorageID); err != nil {
		u.Log.Warnf("Storage not found: %v", err)
		return nil, fiber.ErrNotFound
	}
	item := new(entity.StorageItem)
	if err := u.StorageItemRepository.FindById(tx, item, request.ID); err != nil {
		u.Log.Warnf("Storage item not found: %v", err)
		return nil, fiber.ErrNotFound
	}
	if storage.UserID != uuid.MustParse(request.UserID) {
		return nil, fiber.NewError(fiber.StatusForbidden, "You are not the owner of this storage")
	}

	if request.Weight <= 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Weight must be greater than 0")
	}

	if request.Weight != 0 {
		item.WeightKgs = request.Weight
	}

	if err := u.StorageItemRepository.Update(tx, item); err != nil {
		u.Log.Warnf("Failed to update storage item: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StorageItemToSimpleResponse(item), nil
}

func (u *StorageItemUsecase) Search(ctx context.Context, request *model.SearchStorageItemRequest) ([]model.StorageItemSimpleResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := u.StorageItemRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Warn("Search failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Commit failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.StorageItemSimpleResponse, len(items))
	for i, item := range items {
		responses[i] = *converter.StorageItemToSimpleResponse(&item)
	}

	return responses, total, nil
}

func (u *StorageItemUsecase) Delete(ctx context.Context, id string) (*model.StorageItemSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	item := new(entity.StorageItem)
	if err := u.StorageItemRepository.FindById(tx, item, id); err != nil {
		u.Log.Warnf("Storage item not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := u.StorageItemRepository.Delete(tx, item); err != nil {
		u.Log.Warnf("Delete failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StorageItemToSimpleResponse(item), nil
}
