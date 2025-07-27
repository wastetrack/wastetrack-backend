package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/types"
)

type WasteTransferRequest struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	SourceUserID uuid.UUID `gorm:"column:source_user_id;not null"`
	SourceUser   User      `gorm:"foreignKey:SourceUserID"`

	DestinationUserID uuid.UUID `gorm:"column:destination_user_id;not null"`
	DestinationUser   User      `gorm:"foreignKey:DestinationUserID"`

	// NEW: Added assigned collector functionality
	AssignedCollectorID *uuid.UUID `gorm:"column:assigned_collector_id"` // Nullable
	AssignedCollector   *User      `gorm:"foreignKey:AssignedCollectorID"`

	FormType               string  `gorm:"column:form_type"`
	IsPaid                 bool    `gorm:"column:is_paid;default:false"`
	TotalWeight            float64 `gorm:"column:total_weight;default:0"`
	TotalPrice             int64   `gorm:"column:total_price;default:0"`
	Status                 string  `gorm:"type:request_status;default:'pending'"` // ENUM: pending, assigned, in_progress, completed, cancelled
	ImageURL               string  `gorm:"column:image_url"`
	Notes                  string  `gorm:"column:notes"`
	SourcePhoneNumber      string  `gorm:"column:source_phone_number"`
	DestinationPhoneNumber string  `gorm:"column:destination_phone_number"`

	AppointmentDate      time.Time      `gorm:"type:date"`
	AppointmentStartTime types.TimeOnly `gorm:"type:timetz"`
	AppointmentEndTime   types.TimeOnly `gorm:"type:timetz"`
	AppointmentLocation  *types.Point   `gorm:"type:geography(POINT,4326);"`
	IsDeleted            bool           `gorm:"column:is_deleted;default:false"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	Items []WasteTransferItemOffering `gorm:"foreignKey:TransferFormID"`

	// Distance field for calculated distance (not stored in DB)
	Distance *float64 `gorm:"->" json:"-"`
}
