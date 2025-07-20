package repository

import (
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type SalaryTransactionRepository struct {
	Repository[entity.SalaryTransaction]
	Log *logrus.Logger
}

func NewSalaryTransactionRepository(log *logrus.Logger) *SalaryTransactionRepository {
	return &SalaryTransactionRepository{
		Log: log,
	}
}

func (r *SalaryTransactionRepository) Search(db *gorm.DB, request *model.SearchSalaryTransactionRequest) ([]entity.SalaryTransaction, int64, error) {
	var salaryTransactions []entity.SalaryTransaction
	if err := db.Scopes(r.FilterSalaryTransaction(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&salaryTransactions).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.SalaryTransaction{}).Scopes(r.FilterSalaryTransaction(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return salaryTransactions, total, nil
}

func (r *SalaryTransactionRepository) FilterSalaryTransaction(request *model.SearchSalaryTransactionRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if senderID := request.SenderID; senderID != "" {
			tx = tx.Where("sender_id = ?", senderID)
		}
		if receiverID := request.ReceiverID; receiverID != "" {
			tx = tx.Where("receiver_id = ?", receiverID)
		}
		if transactionType := request.TransactionType; transactionType != "" {
			tx = tx.Where("transaction_type = ?", transactionType)
		}
		if status := request.Status; status != "" {
			tx = tx.Where("status = ?", status)
		}
		if notes := request.Notes; notes != "" {
			tx = tx.Where("notes ILIKE ?", "%"+notes+"%")
		}
		if isDeleted := request.IsDeleted; isDeleted != nil {
			tx = tx.Where("is_deleted = ?", *isDeleted)
		}
		return tx
	}
}

func (r *SalaryTransactionRepository) FindByIdWithRelations(db *gorm.DB, salaryTransaction *entity.SalaryTransaction, id string) error {
	return db.Preload("Sender").Preload("Receiver").Where("id = ?", id).First(salaryTransaction).Error
}
