package repository

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type WasteDropRequestItemRepository struct {
	Repository[entity.WasteDropRequestItem]
	Log *logrus.Logger
}

func NewWasteDropRequestItemRepository(log *logrus.Logger) *WasteDropRequestItemRepository {
	return &WasteDropRequestItemRepository{
		Log: log,
	}
}

func (r *WasteDropRequestItemRepository) FindByID(db *gorm.DB, wasteDropRequestItem *entity.WasteDropRequestItem, id string) error {
	return db.Where("id = ?", id).Preload("Request").Preload("WasteType").First(wasteDropRequestItem).Error
}

func (r *WasteDropRequestItemRepository) Search(db *gorm.DB, request *model.SearchWasteDropRequestItemRequest) ([]entity.WasteDropRequestItem, int64, error) {
	var wasteDropRequestItems []entity.WasteDropRequestItem

	// Apply filters and pagination (no preloads for simple response)
	if err := db.Scopes(r.FilterWasteDropRequestItem(request)).
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&wasteDropRequestItems).Error; err != nil {
		return nil, 0, err
	}

	// Count total records with same filters
	var total int64 = 0
	if err := db.Model(&entity.WasteDropRequestItem{}).
		Scopes(r.FilterWasteDropRequestItem(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return wasteDropRequestItems, total, nil
}

func (r *WasteDropRequestItemRepository) FilterWasteDropRequestItem(request *model.SearchWasteDropRequestItemRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if requestID := request.RequestID; requestID != "" {
			tx = tx.Where("request_id = ?", requestID)
		}
		if wasteTypeID := request.WasteTypeID; wasteTypeID != "" {
			tx = tx.Where("waste_type_id = ?", wasteTypeID)
		}
		return tx
	}
}

func (r *WasteDropRequestItemRepository) FindByDropFormID(db *gorm.DB, dropFormID uuid.UUID) ([]entity.WasteDropRequestItem, error) {
	var items []entity.WasteDropRequestItem
	err := db.Where("request_id = ?", dropFormID).
		Preload("WasteType").
		Find(&items).Error
	return items, err
}

func (r *WasteDropRequestItemRepository) SoftDelete(db *gorm.DB, wasteDropRequestItem *entity.WasteDropRequestItem) error {
	return db.Delete(wasteDropRequestItem).Error
}
