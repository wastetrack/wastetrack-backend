package entity

import (
	"time"

	"github.com/google/uuid"
)

type StorageItem struct {
	ID          uuid.UUID `gorm:"primaryKey;autoIncrement"`
	StorageID   uuid.UUID `gorm:"column:storage_id;not null"`
	Storage     Storage   `gorm:"foreignKey:StorageID"`
	WasteTypeID uuid.UUID `gorm:"column:waste_type_id;not null"`
	WasteType   WasteType `gorm:"foreignKey:WasteTypeID"`
	WeightKgs   float64   `gorm:"column:weight_kgs"`
	CreatedAt   time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:now()"`
}
