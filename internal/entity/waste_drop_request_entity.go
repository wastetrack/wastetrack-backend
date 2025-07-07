package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/types"
)

type WasteDropRequest struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DeliveryType string    `gorm:"type:delivery_type;not null"` // ENUM: 'pickup', 'dropoff'

	CustomerID uuid.UUID `gorm:"column:customer_id;not null"`
	Customer   User      `gorm:"foreignKey:CustomerID"`

	UserPhoneNumber string     `gorm:"column:user_phone_number"`
	WasteBankID     *uuid.UUID `gorm:"column:waste_bank_id;not null"` // Nullable
	WasteBank       User       `gorm:"foreignKey:WasteBankID"`

	AssignedCollectorID *uuid.UUID `gorm:"column:assigned_collector_id"` // Nullable
	AssignedCollector   *User      `gorm:"foreignKey:AssignedCollectorID"`

	TotalPrice int64  `gorm:"column:total_price;default:0"`
	ImageURL   string `gorm:"column:image_url"`
	Status     string `gorm:"type:request_status;default:'pending'"` // ENUM

	AppointmentLocation  *types.Point   `gorm:"type:geography(POINT,4326);"` // Requires custom handling for GEOGRAPHY(Point,4326)
	AppointmentDate      time.Time      `gorm:"type:date"`
	AppointmentStartTime types.TimeOnly `gorm:"type:timetz"`
	AppointmentEndTime   types.TimeOnly `gorm:"type:timetz"`

	Notes     string    `gorm:"column:notes"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Distance field for calculated distance (not stored in DB)
	Distance *float64 `gorm:"-" json:"-"`
}
