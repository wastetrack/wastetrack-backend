package entity

import (
	"github.com/google/uuid"
)

type WasteDropRequestItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	RequestID uuid.UUID        `gorm:"column:request_id;not null"`
	Request   WasteDropRequest `gorm:"foreignKey:RequestID"`

	WasteTypeID uuid.UUID `gorm:"column:waste_type_id;not null"`
	WasteType   WasteType `gorm:"foreignKey:WasteTypeID"`

	Quantity            int64   `gorm:"column:quantity"`        // BIGINT
	VerifiedWeight      float64 `gorm:"column:verified_weight"` // DECIMAL
	VerifiedPricePerKgs int64   `gorm:"column:verified_price_per_kgs"`
	VerifiedSubtotal    int64   `gorm:"column:verified_subtotal"` // BIGINT
	IsDeleted           bool    `gorm:"column:is_deleted;default:false"`
}
