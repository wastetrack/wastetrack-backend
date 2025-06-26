package entity

import "github.com/google/uuid"

type WasteSubcategory struct {
	ID          uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID  uuid.UUID     `gorm:"column:category_id;not null"`
	Category    WasteCategory `gorm:"foreignKey:CategoryID"`
	Name        string        `gorm:"column:name;not null"`
	Description string        `gorm:"column:description"`
}
