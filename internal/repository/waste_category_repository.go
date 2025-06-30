package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type WasteCategoryRepository struct {
	Repository[entity.WasteCategory]
	Log *logrus.Logger
}

func NewWasteCategoryRepository(log *logrus.Logger) *WasteCategoryRepository {
	return &WasteCategoryRepository{
		Log: log,
	}
}

func (r *WasteCategoryRepository) Search(db *gorm.DB, request *model.SearchWasteCategoryRequest) ([]entity.WasteCategory, int64, error) {
	var wasteCategories []entity.WasteCategory
	if err := db.Scopes(r.FilterWasteCategory(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&wasteCategories).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.WasteCategory{}).Scopes(r.FilterWasteCategory(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return wasteCategories, total, nil
}

func (r *WasteCategoryRepository) FilterWasteCategory(request *model.SearchWasteCategoryRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name := request.Name; name != "" {
			tx = tx.Where("name LIKE ?", "%"+name+"%")
		}
		return tx
	}
}
