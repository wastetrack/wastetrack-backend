package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type StorageItemRepository struct {
	Repository[entity.StorageItem]
	Log *logrus.Logger
}

func NewStorageItemRepository(log *logrus.Logger) *StorageItemRepository {
	return &StorageItemRepository{
		Log: log,
	}
}

func (r *StorageItemRepository) FindById(db *gorm.DB, item *entity.StorageItem, id string) error {
	return db.Where("id = ?", id).
		Preload("Storage").
		Preload("WasteType").
		Preload("WasteType.WasteCategory").
		First(item).Error
}

func (r *StorageItemRepository) Search(db *gorm.DB, request *model.SearchStorageItemRequest) ([]entity.StorageItem, int64, error) {
	var items []entity.StorageItem

	query := db.Scopes(r.FilterStorageItem(request))

	switch request.OrderByWeightKgs {
	case "asc":
		query = query.Order("weight_kgs ASC")
	case "desc":
		query = query.Order("weight_kgs DESC")
	}

	if err := query.
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.StorageItem{}).
		Scopes(r.FilterStorageItem(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *StorageItemRepository) FilterStorageItem(request *model.SearchStorageItemRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if request.StorageID != "" {
			tx = tx.Where("storage_id = ?", request.StorageID)
		}
		if request.WasteTypeID != "" {
			tx = tx.Where("waste_type_id = ?", request.WasteTypeID)
		}
		return tx
	}
}
