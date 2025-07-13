package entity

import "github.com/google/uuid"

type Storage struct {
	ID     uuid.UUID `gorm:"primaryKey;autoIncrement"`
	UserID uuid.UUID `gorm:"column:user_id;not null"`
	User   User      `gorm:"foreignKey:UserID"`
	Length float64   `gorm:"column:length"`
	Width  float64   `gorm:"column:width"`
	Height float64   `gorm:"column:height"`
}

func (Storage) TableName() string {
	return "storage"
}
