package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type CollectorManagementRepository struct {
	Repository[entity.CollectorManagement]
	Log *logrus.Logger
}

func NewCollectorManagementRepository(log *logrus.Logger) *CollectorManagementRepository {
	return &CollectorManagementRepository{
		Log: log,
	}
}

func (r *CollectorManagementRepository) Search(db *gorm.DB, request *model.SearchCollectorManagementRequest) ([]entity.CollectorManagement, int64, error) {
	var collectorManagements []entity.CollectorManagement
	if err := db.Scopes(r.FilterCollectorManagement(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&collectorManagements).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.CollectorManagement{}).Scopes(r.FilterCollectorManagement(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return collectorManagements, total, nil
}

func (r *CollectorManagementRepository) FilterCollectorManagement(request *model.SearchCollectorManagementRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if wasteBankID := request.WasteBankID; wasteBankID != "" {
			tx = tx.Where("waste_bank_id = ?", wasteBankID)
		}
		if collectorID := request.CollectorID; collectorID != "" {
			tx = tx.Where("collector_id = ?", collectorID)
		}
		if status := request.Status; status != "" {
			tx = tx.Where("status = ?", status)
		}
		return tx
	}
}

func (r *CollectorManagementRepository) FindByIdWithRelations(db *gorm.DB, collectorManagement *entity.CollectorManagement, id string) error {
	return db.Preload("WasteBank").Preload("Collector").Where("id = ?", id).First(collectorManagement).Error
}

func (r *CollectorManagementRepository) FindByWasteBankAndCollector(db *gorm.DB, collectorManagement *entity.CollectorManagement, wasteBankID, collectorID string) error {
	return db.Where("waste_bank_id = ? AND collector_id = ?", wasteBankID, collectorID).First(collectorManagement).Error
}

func (r *CollectorManagementRepository) FindByWasteBankAndCollectorExcludeID(db *gorm.DB, collectorManagement *entity.CollectorManagement, wasteBankID, collectorID, excludeID string) error {
	return db.Where("waste_bank_id = ? AND collector_id = ? AND id != ?", wasteBankID, collectorID, excludeID).First(collectorManagement).Error
}
