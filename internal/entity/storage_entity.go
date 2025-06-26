package entity

import "github.com/google/uuid"

type Storage struct {
	ID     int64     `gorm:"primaryKey;autoIncrement"`
	UserID uuid.UUID `gorm:"column:user_id;not null"`
	User   User      `gorm:"foreignKey:UserID"`
	Length int64     `gorm:"column:length"`
	Width  int64     `gorm:"column:width"`
	Height int64     `gorm:"column:height"`
}
