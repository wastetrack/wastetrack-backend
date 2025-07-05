package usecase

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/model/converter"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"gorm.io/gorm"
)

type SalaryTransactionUsecase struct {
	DB                          *gorm.DB
	Log                         *logrus.Logger
	Validate                    *validator.Validate
	SalaryTransactionRepository *repository.SalaryTransactionRepository
	UserRepository              *repository.UserRepository
}

func NewSalaryTransactionUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, salaryTransactionRepository *repository.SalaryTransactionRepository, userRepository *repository.UserRepository) *SalaryTransactionUsecase {
	return &SalaryTransactionUsecase{
		DB:                          db,
		Log:                         log,
		Validate:                    validate,
		SalaryTransactionRepository: salaryTransactionRepository,
		UserRepository:              userRepository,
	}
}

func (u *SalaryTransactionUsecase) Create(ctx context.Context, request *model.SalaryTransactionRequest) (*model.SalaryTransactionSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Parse UUIDs
	senderID, err := uuid.Parse(request.SenderID)
	if err != nil {
		u.Log.Warnf("Invalid sender ID: %v", err)
		return nil, fiber.ErrBadRequest
	}

	receiverID, err := uuid.Parse(request.ReceiverID)
	if err != nil {
		u.Log.Warnf("Invalid receiver ID: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if sender exists
	sender := new(entity.User)
	if err := u.UserRepository.FindById(tx, sender, request.SenderID); err != nil {
		u.Log.Warnf("Sender not found: %v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Sender not found")
	}

	// Check if receiver exists
	receiver := new(entity.User)
	if err := u.UserRepository.FindById(tx, receiver, request.ReceiverID); err != nil {
		u.Log.Warnf("Receiver not found: %v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Receiver not found")
	}

	// Check if sender has sufficient balance
	if sender.Balance < request.Amount {
		u.Log.Warnf("Insufficient balance: sender_id=%s, balance=%d, required=%d", request.SenderID, sender.Balance, request.Amount)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Insufficient balance")
	}

	// Perform balance transfer
	sender.Balance -= request.Amount
	receiver.Balance += request.Amount

	// Update sender balance
	if err := u.UserRepository.Update(tx, sender); err != nil {
		u.Log.Warnf("Failed to update sender balance: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update receiver balance
	if err := u.UserRepository.Update(tx, receiver); err != nil {
		u.Log.Warnf("Failed to update receiver balance: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	salaryTransaction := &entity.SalaryTransaction{
		SenderID:        senderID,
		ReceiverID:      receiverID,
		TransactionType: request.TransactionType,
		Amount:          request.Amount,
		Status:          request.Status,
		Notes:           request.Notes,
	}

	if err := u.SalaryTransactionRepository.Create(tx, salaryTransaction); err != nil {
		u.Log.Warnf("Failed to create salary transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SalaryTransactionToSimpleResponse(salaryTransaction), nil
}

func (u *SalaryTransactionUsecase) Get(ctx context.Context, id string) (*model.SalaryTransactionResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	salaryTransaction := new(entity.SalaryTransaction)
	if err := u.SalaryTransactionRepository.FindByIdWithRelations(tx, salaryTransaction, id); err != nil {
		u.Log.Warnf("Salary transaction not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SalaryTransactionToResponse(salaryTransaction), nil
}

func (u *SalaryTransactionUsecase) Update(ctx context.Context, request *model.UpdateSalaryTransactionRequest) (*model.SalaryTransactionSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	salaryTransaction := new(entity.SalaryTransaction)
	if err := u.SalaryTransactionRepository.FindById(tx, salaryTransaction, request.ID); err != nil {
		u.Log.Warnf("Salary transaction not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	// Update transaction type if provided
	if request.TransactionType != "" {
		salaryTransaction.TransactionType = request.TransactionType
	}

	// Update status if provided
	if request.Status != "" {
		salaryTransaction.Status = request.Status
	}

	// Update notes if provided
	if request.Notes != "" {
		salaryTransaction.Notes = request.Notes
	}

	if err := u.SalaryTransactionRepository.Update(tx, salaryTransaction); err != nil {
		u.Log.Warnf("Update failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SalaryTransactionToSimpleResponse(salaryTransaction), nil
}

func (u *SalaryTransactionUsecase) Search(ctx context.Context, request *model.SearchSalaryTransactionRequest) ([]model.SalaryTransactionSimpleResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	// Use search without preloads for simple response
	salaryTransactions, total, err := u.SalaryTransactionRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Warn("Search failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Commit failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.SalaryTransactionSimpleResponse, len(salaryTransactions))
	for i, st := range salaryTransactions {
		responses[i] = *converter.SalaryTransactionToSimpleResponse(&st)
	}

	return responses, total, nil
}

func (u *SalaryTransactionUsecase) Delete(ctx context.Context, id string) (*model.SalaryTransactionSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	salaryTransaction := new(entity.SalaryTransaction)
	if err := u.SalaryTransactionRepository.FindById(tx, salaryTransaction, id); err != nil {
		u.Log.Warnf("Salary transaction not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := u.SalaryTransactionRepository.Delete(tx, salaryTransaction); err != nil {
		u.Log.Warnf("Delete failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.SalaryTransactionToSimpleResponse(salaryTransaction), nil
}
