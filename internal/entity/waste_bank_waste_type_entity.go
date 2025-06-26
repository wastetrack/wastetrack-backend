package entity

import (
	"time"

	"github.com/google/uuid"
)

type WasteBankWasteType struct {
	ID                uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WasteBankID       uuid.UUID        `gorm:"column:waste_bank_id;not null"`
	WasteBank         WasteBankProfile `gorm:"foreignKey:WasteBankID"`
	WasteTypeID       uuid.UUID        `gorm:"column:waste_type_id;not null"`
	WasteType         WasteType        `gorm:"foreignKey:WasteTypeID"`
	CustomPricePerKgs int64            `gorm:"column:custom_price_per_kgs"`
	CreatedAt         time.Time        `gorm:"column:created_at;default:now()"`
	UpdatedAt         time.Time        `gorm:"column:updated_at;default:now()"`
}
