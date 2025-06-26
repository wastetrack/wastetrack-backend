package entity

import "github.com/google/uuid"

type IndustryProfile struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID             uuid.UUID `gorm:"column:user_id;unique;not null"`
	User               User      `gorm:"foreignKey:UserID"`
	TotalWasteWeight   int64     `gorm:"column:total_waste_weight;default:0"`
	TotalRecyledWeight int64     `gorm:"column:total_recyled_weight;default:0"`
}
