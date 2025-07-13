package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
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
