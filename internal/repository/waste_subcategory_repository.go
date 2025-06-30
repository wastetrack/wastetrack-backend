package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type WasteSubCategoryRepository struct {
	Repository[entity.WasteSubcategory]
	Log *logrus.Logger
}

func NewWasteSubCategoryRepository(log *logrus.Logger) *WasteSubCategoryRepository {
	return &WasteSubCategoryRepository{
		Log: log,
	}
}
func (r *WasteSubCategoryRepository) FindById(db *gorm.DB, entity *entity.WasteSubcategory, id string) error {
	return db.Where("id = ?", id).Preload("WasteCategory").Take(entity).Error
}

func (r *WasteSubCategoryRepository) Search(db *gorm.DB, request *model.SearchWasteSubCategoryRequest) ([]entity.WasteSubcategory, int64, error) {
	var wasteSubCategories []entity.WasteSubcategory
	if err := db.Scopes(r.FilterWasteSubCategory(request)).Preload("WasteCategory").Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&wasteSubCategories).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.WasteSubcategory{}).Scopes(r.FilterWasteSubCategory(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return wasteSubCategories, total, nil
}

func (r *WasteSubCategoryRepository) FilterWasteSubCategory(request *model.SearchWasteSubCategoryRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name := request.Name; name != "" {
			tx = tx.Where("name LIKE ?", "%"+name+"%")
		}
		if categoryID := request.CategoryID; categoryID != "" {
			tx = tx.Where("category_id = ?", categoryID)
		}
		return tx
	}
}
