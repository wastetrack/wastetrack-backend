package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

type WasteBankRepository struct {
	Repository[entity.WasteBankProfile]
	Log *logrus.Logger
}

func NewWasteBankRepository(log *logrus.Logger) *WasteBankRepository {
	return &WasteBankRepository{
		Log: log,
	}
}

func (r *WasteBankRepository) FindByUserID(db *gorm.DB, profile *entity.WasteBankProfile, userID string) error {
	return db.Where("user_id = ?", userID).Preload("User").First(profile).Error
}
