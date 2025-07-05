package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func SalaryTransactionToSimpleResponse(salaryTransaction *entity.SalaryTransaction) *model.SalaryTransactionSimpleResponse {
	return &model.SalaryTransactionSimpleResponse{
		ID:              salaryTransaction.ID.String(),
		SenderID:        salaryTransaction.SenderID.String(),
		ReceiverID:      salaryTransaction.ReceiverID.String(),
		TransactionType: salaryTransaction.TransactionType,
		Amount:          salaryTransaction.Amount,
		CreatedAt:       salaryTransaction.CreatedAt.String(),
		Status:          salaryTransaction.Status,
		Notes:           salaryTransaction.Notes,
	}
}

func SalaryTransactionToResponse(salaryTransaction *entity.SalaryTransaction) *model.SalaryTransactionResponse {
	var sender *model.UserResponse
	if salaryTransaction.SenderID != uuid.Nil {
		sender = UserToResponse(&salaryTransaction.Sender)
	}
	var receiver *model.UserResponse
	if salaryTransaction.ReceiverID != uuid.Nil {
		receiver = UserToResponse(&salaryTransaction.Receiver)
	}
	return &model.SalaryTransactionResponse{
		ID:              salaryTransaction.ID.String(),
		SenderID:        salaryTransaction.SenderID.String(),
		ReceiverID:      salaryTransaction.ReceiverID.String(),
		TransactionType: salaryTransaction.TransactionType,
		Amount:          salaryTransaction.Amount,
		CreatedAt:       salaryTransaction.CreatedAt.String(),
		Status:          salaryTransaction.Status,
		Notes:           salaryTransaction.Notes,
		Sender:          sender,
		Receiver:        receiver,
	}
}
