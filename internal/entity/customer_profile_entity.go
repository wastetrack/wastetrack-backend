package entity

import (
	"github.com/google/uuid"
)

type CustomerProfile struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID        uuid.UUID `gorm:"column:user_id;unique;not null"`
	User          User      `gorm:"foreignKey:UserID"`
	CarbonDeficit int64     `gorm:"column:carbon_deficit;default:0"`
	WaterSaved    int64     `gorm:"column:water_saved;default:0"`
	BagsStored    int64     `gorm:"column:bags_stored;default:0"`
	Trees         int64     `gorm:"column:trees;default:0"`
}
