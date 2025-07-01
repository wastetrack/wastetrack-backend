package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type WasteBankPricedTypeRepository struct {
	Repository[entity.WasteBankPricedType]
	Log *logrus.Logger
}

func NewWasteBankPricedTypeRepository(log *logrus.Logger) *WasteBankPricedTypeRepository {
	return &WasteBankPricedTypeRepository{
		Log: log,
	}
}

func (r *WasteBankPricedTypeRepository) FindById(db *gorm.DB, wpt *entity.WasteBankPricedType, id string) error {
	return db.
		Preload("WasteBank").
		Preload("WasteType").
		Where("id = ?", id).
		Take(wpt).
		Error
}

func (r *WasteBankPricedTypeRepository) Search(db *gorm.DB, req *model.SearchWasteBankPricedTypeRequest) ([]entity.WasteBankPricedType, int64, error) {
	var result []entity.WasteBankPricedType
	var total int64

	query := db.Model(&entity.WasteBankPricedType{}).
		Scopes(r.FilterWasteBankPricedType(req))

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Offset((req.Page - 1) * req.Size).
		Limit(req.Size).
		Find(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func (r *WasteBankPricedTypeRepository) FilterWasteBankPricedType(req *model.SearchWasteBankPricedTypeRequest) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if req.WasteBankID != "" {
			tx = tx.Where("waste_bank_id = ?", req.WasteBankID)
		}
		if req.WasteTypeID != "" {
			tx = tx.Where("waste_type_id = ?", req.WasteTypeID)
		}
		return tx
	}
}
