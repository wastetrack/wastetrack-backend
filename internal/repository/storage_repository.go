package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type StorageRepository struct {
	Repository[entity.Storage]
	Log *logrus.Logger
}

func NewStorageRepository(log *logrus.Logger) *StorageRepository {
	return &StorageRepository{
		Log: log,
	}
}

func (r *StorageRepository) FindById(db *gorm.DB, storage *entity.Storage, id string) error {
	return db.Where("id = ?", id).Preload("User").First(storage).Error
}

func (r *StorageRepository) Search(db *gorm.DB, request *model.SearchStorageRequest) ([]entity.Storage, int64, error) {
	var storages []entity.Storage

	query := db.Scopes(r.FilterStorage(request))

	if err := query.Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&storages).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.Storage{}).Scopes(r.FilterStorage(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return storages, total, nil
}

func (r *StorageRepository) FilterStorage(request *model.SearchStorageRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if request.UserID != "" {
			tx = tx.Where("user_id = ?", request.UserID)
		}
		if request.IsForRecycledMaterial != nil {
			tx = tx.Where("is_for_recycled_material = ?", *request.IsForRecycledMaterial)
		}
		if request.MinLength != nil {
			tx = tx.Where("length >= ?", *request.MinLength)
		}
		if request.MaxLength != nil {
			tx = tx.Where("length <= ?", *request.MaxLength)
		}
		if request.MinWidth != nil {
			tx = tx.Where("width >= ?", *request.MinWidth)
		}
		if request.MaxWidth != nil {
			tx = tx.Where("width <= ?", *request.MaxWidth)
		}
		if request.MinHeight != nil {
			tx = tx.Where("height >= ?", *request.MinHeight)
		}
		if request.MaxHeight != nil {
			tx = tx.Where("height <= ?", *request.MaxHeight)
		}
		return tx
	}
}
