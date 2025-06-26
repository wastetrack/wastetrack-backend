package entity

import "github.com/google/uuid"

type WasteTransferItemOffering struct {
	ID                  uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TransferFormID      int64             `gorm:"column:transfer_form_id;not null"`
	TransferForm        WasteTransferForm `gorm:"foreignKey:TransferFormID"`
	WasteTypeID         uuid.UUID         `gorm:"column:waste_type_id;not null"`
	WasteType           WasteType         `gorm:"foreignKey:WasteTypeID"`
	OfferingWeight      float64           `gorm:"column:offering_weight"`
	OfferingPricePerKgs float64           `gorm:"column:offering_price_per_kgs"`
	AcceptedWeight      float64           `gorm:"column:accepted_weight"`
	AcceptedPricePerKgs float64           `gorm:"column:accepted_price_per_kgs"`
}
