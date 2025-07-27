package entity

import (
	"time"

	"github.com/google/uuid"
)

type PointConversion struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"column:user_id;not null"`
	User      User      `gorm:"foreignKey:UserID"`
	Amount    int64     `gorm:"column:amount;default:0"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	Status    string    `gorm:"column:status;default:'pending'"`
	IsDeleted bool      `gorm:"column:is_deleted;default:false"`
}
