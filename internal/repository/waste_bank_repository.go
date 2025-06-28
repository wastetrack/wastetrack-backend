package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
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
