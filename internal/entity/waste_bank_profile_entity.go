package entity

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/types"
)

type WasteBankProfile struct {
	ID               uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           uuid.UUID      `gorm:"column:user_id;unique;not null"`
	TotalWasteWeight float64        `gorm:"column:total_waste_weight;default:0"`
	TotalWorkers     int64          `gorm:"column:total_workers;default:0"`
	OpenTime         types.TimeOnly `gorm:"column:open_time;type:time"`
	CloseTime        types.TimeOnly `gorm:"column:close_time;type:time"`
	User             User           `gorm:"foreignKey:UserID"`
}
