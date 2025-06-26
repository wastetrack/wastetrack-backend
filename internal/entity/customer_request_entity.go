package entity

import (
	"time"

	"github.com/google/uuid"
)

type CustomerRequest struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	DeliveryType         string    `gorm:"column:delivery_type;not null"`
	CustomerID           uuid.UUID `gorm:"column:customer_id;not null"`
	Customer             User      `gorm:"foreignKey:CustomerID"`
	UserPhoneNumber      string    `gorm:"column:user_phone_number"`
	WasteBankID          uuid.UUID `gorm:"column:waste_bank_id"`
	WasteBank            User      `gorm:"foreignKey:WasteBankID"`
	TotalPrice           int64     `gorm:"column:total_price;default:0"`
	Status               string    `gorm:"column:status;default:'pending'"`
	AppointmentDate      time.Time `gorm:"column:appointment_date;type:date"`
	AppointmentStartTime time.Time `gorm:"column:appointment_start_time;type:time"`
	AppointmentEndTime   time.Time `gorm:"column:appointment_end_time;type:time"`
	Notes                string    `gorm:"column:notes"`
	CreatedAt            time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt            time.Time `gorm:"column:updated_at;default:now()"`
}
