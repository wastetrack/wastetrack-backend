package entity

import "github.com/google/uuid"

type GovernmentProfile struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"column:user_id;unique;not null"`
	User   User      `gorm:"foreignKey:UserID"`
}
