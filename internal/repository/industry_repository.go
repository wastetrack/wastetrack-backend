package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

type IndustryRepository struct {
	Repository[entity.IndustryProfile]
	Log *logrus.Logger
}

func NewIndustryRepository(log *logrus.Logger) *IndustryRepository {
	return &IndustryRepository{
		Log: log,
	}
}

func (r *IndustryRepository) FindByUserID(db *gorm.DB, profile *entity.IndustryProfile, userID string) error {
	return db.Where("user_id = ?", userID).Preload("User").First(profile).Error
}
