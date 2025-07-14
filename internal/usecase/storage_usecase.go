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

type StorageUsecase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	StorageRepository *repository.StorageRepository
	UserRepository    *repository.UserRepository
}

func NewStorageUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, storageRepository *repository.StorageRepository, userRepository *repository.UserRepository) *StorageUsecase {
	return &StorageUsecase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		StorageRepository: storageRepository,
		UserRepository:    userRepository,
	}
}

func (u *StorageUsecase) Create(ctx context.Context, request *model.StorageRequest) (*model.StorageSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.UserID); err != nil {
		u.Log.Warnf("User not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	storage := &entity.Storage{
		UserID:                user.ID,
		Length:                request.Length,
		Width:                 request.Width,
		Height:                request.Height,
		IsForRecycledMaterial: request.IsForRecycledMaterial,
	}

	if err := u.StorageRepository.Create(tx, storage); err != nil {
		u.Log.Warnf("Failed to create storage: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StorageToSimpleResponse(storage), nil
}

func (u *StorageUsecase) Get(ctx context.Context, id string) (*model.StorageResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	storage := new(entity.Storage)
	if err := u.StorageRepository.FindById(tx, storage, id); err != nil {
		u.Log.Warnf("Storage not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StorageToResponse(storage), nil
}

func (u *StorageUsecase) Update(ctx context.Context, request *model.UpdateStorageRequest) (*model.StorageSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	storage := new(entity.Storage)
	if err := u.StorageRepository.FindById(tx, storage, request.ID); err != nil {
		u.Log.Warnf("Storage not found: %v", err)
		return nil, fiber.ErrNotFound
	}
	if storage.UserID != uuid.MustParse(request.UserID) {
		return nil, fiber.NewError(fiber.StatusForbidden, "You are not the owner of this storage")
	}

	if request.Height != 0 {
		storage.Height = request.Height
	}
	if request.Width != 0 {
		storage.Width = request.Width
	}
	if request.Length != 0 {
		storage.Length = request.Length
	}
	if request.IsForRecycledMaterial != nil {
		storage.IsForRecycledMaterial = *request.IsForRecycledMaterial
	}

	if err := u.StorageRepository.Update(tx, storage); err != nil {
		u.Log.Warnf("Update failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StorageToSimpleResponse(storage), nil
}

func (u *StorageUsecase) Search(ctx context.Context, request *model.SearchStorageRequest) ([]model.StorageSimpleResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	storages, total, err := u.StorageRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Warn("Search failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Commit failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.StorageSimpleResponse, len(storages))
	for i, s := range storages {
		responses[i] = *converter.StorageToSimpleResponse(&s)
	}

	return responses, total, nil
}

func (u *StorageUsecase) Delete(ctx context.Context, id string) (*model.StorageSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	storage := new(entity.Storage)
	if err := u.StorageRepository.FindById(tx, storage, id); err != nil {
		u.Log.Warnf("Storage not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := u.StorageRepository.Delete(tx, storage); err != nil {
		u.Log.Warnf("Delete failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.StorageToSimpleResponse(storage), nil
}
