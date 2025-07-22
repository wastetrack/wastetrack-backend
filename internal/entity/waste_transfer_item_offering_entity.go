package entity

import (
	"github.com/google/uuid"
)

type WasteTransferItemOffering struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	TransferFormID uuid.UUID            `gorm:"column:transfer_request_id;not null"`
	TransferForm   WasteTransferRequest `gorm:"foreignKey:TransferFormID"`

	WasteTypeID uuid.UUID `gorm:"column:waste_type_id;not null"`
	WasteType   WasteType `gorm:"foreignKey:WasteTypeID"`

	// Offering details (initial proposal)
	OfferingWeight      float64 `gorm:"column:offering_weight"`        // DECIMAL - offered weight
	OfferingPricePerKgs int64   `gorm:"column:offering_price_per_kgs"` // DECIMAL - offered price per kg

	// Accepted details (filled when collector is assigned)
	AcceptedWeight      float64 `gorm:"column:accepted_weight;default:0"`        // DECIMAL - accepted weight
	AcceptedPricePerKgs int64   `gorm:"column:accepted_price_per_kgs;default:0"` // DECIMAL - accepted price per kg

	VerifiedWeight float64 `gorm:"column:verified_weight;default:0"` // DECIMAL - verified weight

	// Recycling process
	RecycledWeight float64 `gorm:"column:recycled_weight;default:0"` // DECIMAL - weight of actual recycled material
}

func (WasteTransferItemOffering) TableName() string {
	return "waste_transfer_items"
}
