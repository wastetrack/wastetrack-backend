package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

type WasteCollectorRepository struct {
	Repository[entity.WasteCollectorProfile]
	Log *logrus.Logger
}

func NewWasteCollectorRepository(log *logrus.Logger) *WasteCollectorRepository {
	return &WasteCollectorRepository{
		Log: log,
	}
}

func (r *WasteCollectorRepository) FindByUserID(db *gorm.DB, profile *entity.WasteCollectorProfile, userID string) error {
	return db.Where("user_id = ?", userID).Preload("User").First(profile).Error
}
