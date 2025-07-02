package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
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
