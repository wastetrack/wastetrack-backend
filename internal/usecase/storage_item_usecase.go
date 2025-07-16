package usecase

import (
	"context"
	"time"

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
	WasteTypeRepository   *repository.WasteTypeRepository
}

func NewStorageItemUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	storageRepo *repository.StorageRepository,
	storageItemRepo *repository.StorageItemRepository,
	wasteTypeRepo *repository.WasteTypeRepository,
) *StorageItemUsecase {
	return &StorageItemUsecase{
		DB:                    db,
		Log:                   log,
		Validate:              validate,
		StorageRepository:     storageRepo,
		StorageItemRepository: storageItemRepo,
		WasteTypeRepository:   wasteTypeRepo,
	}
}

func (c *StorageItemUsecase) Create(ctx context.Context, request *model.StorageItemRequest) (*model.StorageItemSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Parse and validate UUIDs
	storageID, err := uuid.Parse(request.StorageID)
	if err != nil {
		c.Log.Warnf("Invalid storage ID: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteTypeID, err := uuid.Parse(request.WasteTypeID)
	if err != nil {
		c.Log.Warnf("Invalid waste type ID: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Validate weight is positive
	if request.WeightKgs <= 0 {
		c.Log.Warnf("Weight must be positive: %f", request.WeightKgs)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Weight must be greater than 0")
	}

	// Check if storage exists
	storage := new(entity.Storage)
	if err := c.StorageRepository.FindById(tx, storage, request.StorageID); err != nil {
		c.Log.Warnf("Failed to find storage by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Check if waste type exists
	wasteType := new(entity.WasteType)
	if err := c.WasteTypeRepository.FindById(tx, wasteType, request.WasteTypeID); err != nil {
		c.Log.Warnf("Failed to find waste type by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// NEW: Check if storage item with this storage_id and waste_type_id combination already exists
	var existingStorageItem entity.StorageItem
	err = tx.Where("storage_id = ? AND waste_type_id = ?", storageID, wasteTypeID).
		First(&existingStorageItem).Error

	switch err {
	case nil:
		// Storage item combination exists, add to existing weight
		c.Log.Infof("Found existing storage item for storage %s and waste type %s: adding %f kg to existing %f kg",
			storageID.String(), wasteTypeID.String(), request.WeightKgs, existingStorageItem.WeightKgs)

		existingStorageItem.WeightKgs += request.WeightKgs
		existingStorageItem.UpdatedAt = time.Now()

		if err := c.StorageItemRepository.Update(tx, &existingStorageItem); err != nil {
			c.Log.Warnf("Failed to update existing storage item: %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		if err := tx.Commit().Error; err != nil {
			c.Log.Warnf("Failed to commit transaction: %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		c.Log.Infof("Successfully updated existing storage item ID: %s with new total weight: %f kg",
			existingStorageItem.ID.String(), existingStorageItem.WeightKgs)

		return converter.StorageItemToSimpleResponse(&existingStorageItem), nil

	case gorm.ErrRecordNotFound:
		// Storage item combination doesn't exist, create new one
		c.Log.Infof("Creating new storage item for storage %s and waste type %s with weight %f kg",
			storageID.String(), wasteTypeID.String(), request.WeightKgs)

		storageItem := &entity.StorageItem{
			StorageID:   storageID,
			WasteTypeID: wasteTypeID,
			WeightKgs:   request.WeightKgs,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := c.StorageItemRepository.Create(tx, storageItem); err != nil {
			c.Log.Warnf("Failed to create storage item: %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		if err := tx.Commit().Error; err != nil {
			c.Log.Warnf("Failed to commit transaction: %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		c.Log.Infof("Successfully created new storage item ID: %s", storageItem.ID.String())
		return converter.StorageItemToSimpleResponse(storageItem), nil

	default:
		// Database error
		c.Log.Warnf("Database error while checking for existing storage item: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
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
