package entity

import (
	"time"

	"github.com/google/uuid"
)

type WasteTransferForm struct {
	ID                     int64     `gorm:"primaryKey;autoIncrement"`
	SourceUserID           uuid.UUID `gorm:"column:source_user_id;not null"`
	SourceUser             User      `gorm:"foreignKey:SourceUserID"`
	DestinationUserID      uuid.UUID `gorm:"column:destination_user_id;not null"`
	DestinationUser        User      `gorm:"foreignKey:DestinationUserID"`
	FormType               int64     `gorm:"column:form_type"`
	TotalWeight            int64     `gorm:"column:total_weight;default:0"`
	TotalPrice             int64     `gorm:"column:total_price;default:0"`
	Status                 string    `gorm:"column:status"`
	SourcePhoneNumber      string    `gorm:"column:source_phone_number"`
	DestinationPhoneNumber string    `gorm:"column:destination_phone_number"`
	AppointmentDate        time.Time `gorm:"column:appointment_date;type:date"`
	AppointmentTime        time.Time `gorm:"column:appointment_time;type:time"`
	CreatedAt              time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt              time.Time `gorm:"column:updated_at;default:now()"`
}
