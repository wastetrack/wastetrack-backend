package repository

import (
	"fmt"

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

func (r *StorageRepository) GetWasteBankStorageVolumes(db *gorm.DB, request *model.GovernmentDashboardRequest) (map[string]float64, error) {
	var results []struct {
		UserID      string  `json:"user_id"`
		TotalVolume float64 `json:"total_volume"`
	}

	query := db.Raw(`
		SELECT 
			s.user_id::text as user_id,
			SUM(s.length * s.width * s.height) as total_volume
		FROM storage s
		JOIN users u ON s.user_id = u.id
		WHERE s.is_deleted = false
		  AND u.role IN ('waste_bank_unit', 'waste_bank_central')
		  AND ($1 = '' OR u.province ILIKE $2)
		  AND ($3 = '' OR u.city ILIKE $4)
		GROUP BY s.user_id
		HAVING SUM(s.length * s.width * s.height) > 0
	`,
		request.Province, "%"+request.Province+"%",
		request.City, "%"+request.City+"%")

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get storage volumes: %v", err)
	}

	// Convert to map for easy lookup
	volumes := make(map[string]float64)
	for _, result := range results {
		volumes[result.UserID] = result.TotalVolume
	}

	return volumes, nil
}
