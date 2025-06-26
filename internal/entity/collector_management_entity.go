package entity

import (
	"github.com/google/uuid"
)

type CollectorManagement struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WasteBankID uuid.UUID `gorm:"column:waste_bank_id;not null"`
	WasteBank   User      `gorm:"foreignKey:WasteBankID"`
	CollectorID uuid.UUID `gorm:"column:collector_id;not null"`
	Collector   User      `gorm:"foreignKey:CollectorID"`
	Status      string    `gorm:"column:status;default:'pending'"`
}
