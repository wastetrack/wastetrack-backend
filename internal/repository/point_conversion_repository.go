package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type PointConversionRepository struct {
	Repository[entity.PointConversion]
	Log *logrus.Logger
}

func NewPointConversionRepository(log *logrus.Logger) *PointConversionRepository {
	return &PointConversionRepository{
		Log: log,
	}
}

func (r *PointConversionRepository) Search(db *gorm.DB, request *model.SearchPointConversionRequest) ([]entity.PointConversion, int64, error) {
	var pointConversions []entity.PointConversion
	if err := db.Scopes(r.FilterPointConversion(request), r.OrderPointConversion(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&pointConversions).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.PointConversion{}).Scopes(r.FilterPointConversion(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return pointConversions, total, nil
}

func (r *PointConversionRepository) FilterPointConversion(request *model.SearchPointConversionRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if userID := request.UserID; userID != "" {
			tx = tx.Where("user_id = ?", userID)
		}
		if amount := request.Amount; amount > 0 {
			tx = tx.Where("amount = ?", amount)
		}
		if status := request.Status; status != "" {
			tx = tx.Where("status = ?", status)
		}
		if isDeleted := request.IsDeleted; isDeleted != nil {
			tx = tx.Where("is_deleted = ?", *isDeleted)
		}
		return tx
	}
}

func (r *PointConversionRepository) OrderPointConversion(request *model.SearchPointConversionRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		orderBy := request.OrderBy
		orderDir := request.OrderDir

		// Set defaults if not provided
		if orderBy == "" {
			orderBy = "created_at"
		}
		if orderDir == "" {
			orderDir = "desc"
		}

		// Validate order direction
		if orderDir != "asc" && orderDir != "desc" {
			orderDir = "desc"
		}

		// Validate order by column (whitelist approach for security)
		validColumns := map[string]bool{
			"id":         true,
			"user_id":    true,
			"amount":     true,
			"created_at": true,
			"status":     true,
			"is_deleted": true,
		}

		if !validColumns[orderBy] {
			orderBy = "created_at"
		}

		return tx.Order(orderBy + " " + orderDir)
	}
}

func (r *PointConversionRepository) FindByIdWithRelations(db *gorm.DB, pointConversion *entity.PointConversion, id string) error {
	return db.Preload("User").Where("id = ?", id).First(pointConversion).Error
}
