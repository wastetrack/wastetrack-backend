package entity

import "github.com/google/uuid"

type WasteType struct {
	ID            uuid.UUID        `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CategoryID    uuid.UUID        `gorm:"column:category_id;not null"`
	Category      WasteCategory    `gorm:"foreignKey:CategoryID"`
	SubcategoryID uuid.UUID        `gorm:"column:subcategory_id"`
	Subcategory   WasteSubcategory `gorm:"foreignKey:SubcategoryID"`
	Name          string           `gorm:"column:name;not null"`
}
