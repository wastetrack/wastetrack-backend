package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Repository[entity.CustomerProfile]
	Log *logrus.Logger
}

func NewCustomerRepository(log *logrus.Logger) *CustomerRepository {
	return &CustomerRepository{
		Log: log,
	}
}

func (r *CustomerRepository) FindByUserID(db *gorm.DB, profile *entity.CustomerProfile, userID string) error {
	return db.Where("user_id = ?", userID).Preload("User").First(profile).Error
}

func (r *CustomerRepository) FindByUserIDNoPreload(db *gorm.DB, profile *entity.CustomerProfile, userID string) error {
	return db.Where("user_id = ?", userID).First(profile).Error
}
