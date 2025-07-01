package entity

import "github.com/google/uuid"

type WasteType struct {
	ID               uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID       uuid.UUID        `gorm:"column:category_id;not null"`
	SubcategoryID    uuid.UUID        `gorm:"column:subcategory_id"`
	Name             string           `gorm:"column:name;not null"`
	Description      string           `gorm:"column:description"`
	WasteCategory    WasteCategory    `gorm:"foreignKey:CategoryID"`
	WasteSubcategory WasteSubcategory `gorm:"foreignKey:SubcategoryID"`
}
