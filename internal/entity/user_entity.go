package entity

import (
	"time"

	"github.com/google/uuid"
)

// PostGIS-compatible point using Postgres' GEOGRAPHY(Point, 4326)
type User struct {
	ID                     uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username               string     `gorm:"column:username;unique;not null"`
	Email                  string     `gorm:"column:email;unique;not null"`
	Password               string     `gorm:"column:password;not null"`
	Role                   string     `gorm:"column:role;type:user_role;default:'customer';not null"`
	PhoneNumber            string     `gorm:"column:phone_number"`
	Institution            string     `gorm:"column:institution"`
	Address                string     `gorm:"column:address"`
	City                   string     `gorm:"column:city"`
	Province               string     `gorm:"column:province"`
	Points                 int64      `gorm:"column:points;default:0"`
	Balance                int64      `gorm:"column:balance;default:0"`
	Location               *string    `gorm:"type:geometry(POINT,4326);"`
	IsEmailVerified        bool       `gorm:"column:is_email_verified;default:false"`
	EmailVerificationToken string     `gorm:"column:email_verification_token"`
	ResetPasswordToken     string     `gorm:"column:reset_password_token"`
	ResetPasswordExpiry    *time.Time `gorm:"column:reset_password_expiry"`
	CreatedAt              time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt              time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}
