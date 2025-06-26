package entity

import (
	"time"

	"github.com/google/uuid"
)

type WasteBankProfile struct {
	ID               uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           uuid.UUID `gorm:"column:user_id;unique;not null"`
	User             User      `gorm:"foreignKey:UserID"`
	Balance          int64     `gorm:"column:balance;default:0"`
	TotalWasteWeight int64     `gorm:"column:total_waste_weight;default:0"`
	TotalWorkers     int64     `gorm:"column:total_workers;default:0"`
	OpenTime         time.Time `gorm:"column:open_time;type:time"`
	CloseTime        time.Time `gorm:"column:close_time;type:time"`
}
