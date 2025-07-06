package entity

import (
	"github.com/google/uuid"
)

type WasteTransferItemOffering struct {
	ID                  uuid.UUID            `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TransferFormID      uuid.UUID            `gorm:"column:transfer_request_id;not null"`
	TransferForm        WasteTransferRequest `gorm:"foreignKey:TransferFormID"`
	WasteTypeID         uuid.UUID            `gorm:"column:waste_type_id;not null"`
	WasteType           WasteType            `gorm:"foreignKey:WasteTypeID"`
	OfferingWeight      float64              `gorm:"column:offering_weight;default:0"`
	OfferingPricePerKgs float64              `gorm:"column:offering_price_per_kgs;default:0"`
	AcceptedWeight      float64              `gorm:"column:accepted_weight;default:0"`
	AcceptedPricePerKgs float64              `gorm:"column:accepted_price_per_kgs;default:0"`
}

func (WasteTransferItemOffering) TableName() string {
	return "waste_transfer_items"
}
