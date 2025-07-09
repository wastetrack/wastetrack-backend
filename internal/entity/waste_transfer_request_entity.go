package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/types"
)

type WasteTransferRequest struct {
	ID                     uuid.UUID                   `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SourceUserID           uuid.UUID                   `gorm:"column:source_user_id;not null"`
	SourceUser             User                        `gorm:"foreignKey:SourceUserID"`
	DestinationUserID      uuid.UUID                   `gorm:"column:destination_user_id;not null"`
	DestinationUser        User                        `gorm:"foreignKey:DestinationUserID"`
	FormType               string                      `gorm:"column:form_type"`
	TotalWeight            int64                       `gorm:"column:total_weight;default:0"`
	TotalPrice             int64                       `gorm:"column:total_price;default:0"`
	Status                 string                      `gorm:"column:status;default:'pending'"`
	ImageURL               string                      `gorm:"column:image_url"`
	Notes                  string                      `gorm:"column:notes"`
	SourcePhoneNumber      string                      `gorm:"column:source_phone_number"`
	DestinationPhoneNumber string                      `gorm:"column:destination_phone_number"`
	AppointmentDate        time.Time                   `gorm:"column:appointment_date;type:date"`
	AppointmentStartTime   types.TimeOnly              `gorm:"column:appointment_start_time;type:timetz"`
	AppointmentEndTime     types.TimeOnly              `gorm:"column:appointment_end_time;type:timetz"`
	AppointmentLocation    *types.Point                `gorm:"type:geography(POINT,4326);"`
	CreatedAt              time.Time                   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt              time.Time                   `gorm:"column:updated_at;autoUpdateTime"`
	Items                  []WasteTransferItemOffering `gorm:"foreignKey:TransferFormID"`
	Distance               *float64                    `gorm:"->" json:"-"`
}
