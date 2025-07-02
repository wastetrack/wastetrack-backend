package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type WasteTypeRepository struct {
	Repository[entity.WasteType]
	Log *logrus.Logger
}

func NewWasteTypeRepository(log *logrus.Logger) *WasteTypeRepository {
	return &WasteTypeRepository{
		Log: log,
	}
}

func (r *WasteTypeRepository) FindById(db *gorm.DB, entity *entity.WasteType, id string) error {
	return db.
		Where("id = ?", id).
		Preload("WasteCategory").
		Take(entity).
		Error
}

func (r *WasteTypeRepository) Search(db *gorm.DB, request *model.SearchWasteTypeRequest) ([]entity.WasteType, int64, error) {
	var wasteTypes []entity.WasteType
	if err := db.
		Scopes(r.FilterWasteType(request)).
		Preload("WasteCategory").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&wasteTypes).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.WasteType{}).
		Scopes(r.FilterWasteType(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return wasteTypes, total, nil
}

func (r *WasteTypeRepository) FilterWasteType(request *model.SearchWasteTypeRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name := request.Name; name != "" {
			tx = tx.Where("name ILIKE ?", "%"+name+"%")
		}
		if categoryID := request.CategoryID; categoryID != "" {
			tx = tx.Where("category_id = ?", categoryID)
		}
		return tx
	}
}
