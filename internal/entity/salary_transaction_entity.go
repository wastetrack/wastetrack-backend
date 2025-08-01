package entity

import (
	"time"

	"github.com/google/uuid"
)

type SalaryTransaction struct {
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SenderID        uuid.UUID `gorm:"column:sender_id;not null"`
	Sender          User      `gorm:"foreignKey:SenderID"`
	ReceiverID      uuid.UUID `gorm:"column:receiver_id;not null"`
	Receiver        User      `gorm:"foreignKey:ReceiverID"`
	Amount          int64     `gorm:"column:amount;default:0"`
	TransactionType string    `gorm:"column:transaction_type;not null"`
	IsDeleted       bool      `gorm:"column:is_deleted;default:false"`
	CreatedAt       time.Time `gorm:"column:created_at;default:now()"`
	Status          string    `gorm:"column:status;default:'pending'"`
	Notes           string    `gorm:"column:notes"`
}
