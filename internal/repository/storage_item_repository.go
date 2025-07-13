package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
)

type StorageItemRepository struct {
	Repository[entity.StorageItem]
	Log *logrus.Logger
}

func NewStorageItemRepository(log *logrus.Logger) *StorageItemRepository {
	return &StorageItemRepository{
		Log: log,
	}
}
