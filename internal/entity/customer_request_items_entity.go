package entity

import "github.com/google/uuid"

type CustomerRequestItem struct {
	ID               uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RequestID        uuid.UUID       `gorm:"column:request_id;not null"`
	Request          CustomerRequest `gorm:"foreignKey:RequestID"`
	WasteTypeID      uuid.UUID       `gorm:"column:waste_type_id;not null"`
	WasteType        WasteType       `gorm:"foreignKey:WasteTypeID"`
	Quantity         float64         `gorm:"column:quantity"`
	VerifiedWeight   float64         `gorm:"column:verified_weight"`
	VerifiedSubtotal int64           `gorm:"column:verified_subtotal"`
}
