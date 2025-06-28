package usecase

import (
	"context"
	"fmt"
	"time"

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

type WasteBankUsecase struct {
	DB                  *gorm.DB
	Log                 *logrus.Logger
	Validate            *validator.Validate
	WasteBankRepository *repository.WasteBankRepository
}

func NewWasteBankUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteBankRepository *repository.WasteBankRepository) *WasteBankUsecase {
	return &WasteBankUsecase{
		DB:                  db,
		Log:                 log,
		Validate:            validate,
		WasteBankRepository: wasteBankRepository,
	}
}

func (c *WasteBankUsecase) Create(ctx context.Context, request *model.WasteBankRequest) (*model.WasteBankResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		c.Log.Warnf("Failed to start transaction: %+v", tx.Error)
		return nil, fiber.ErrInternalServerError
	}
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	userID, err := uuid.Parse(request.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id format: %w", err)
	}

	wasteBank := &entity.WasteBankProfile{
		UserID: userID,
	}
	// Handle optional fields
	if request.TotalWasteWeight != nil {
		wasteBank.TotalWasteWeight = *request.TotalWasteWeight
	}

	if request.TotalWorkers != nil {
		wasteBank.TotalWorkers = *request.TotalWorkers
	}

	// Handle time parsing if provided
	if request.OpenTime != nil && *request.OpenTime != "" {
		openTime, err := time.Parse("15:04:05", *request.OpenTime)
		if err != nil {
			return nil, fmt.Errorf("invalid open_time format: %w", err)
		}
		wasteBank.OpenTime = openTime
	}

	if request.CloseTime != nil && *request.CloseTime != "" {
		closeTime, err := time.Parse("15:04:05", *request.CloseTime)
		if err != nil {
			return nil, fmt.Errorf("invalid close_time format: %w", err)
		}
		wasteBank.CloseTime = closeTime
	}

	if err := c.WasteBankRepository.Create(tx, wasteBank); err != nil {
		c.Log.Warnf("Failed to create waste bank: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteBankToResponse(wasteBank), nil
}

//TODO: Search,Get

func (c *WasteBankUsecase) Update(ctx context.Context, request *model.WasteBankUpdateRequest) (*model.WasteBankResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	wasteBank := new(entity.WasteBankProfile)
	if err := c.WasteBankRepository.FindById(tx, wasteBank, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if request.TotalWasteWeight != nil {
		wasteBank.TotalWasteWeight = *request.TotalWasteWeight
	}

	if request.TotalWorkers != nil {
		wasteBank.TotalWorkers = *request.TotalWorkers
	}

	if request.OpenTime != nil && *request.OpenTime != "" {
		openTime, err := time.Parse("15:04:05", *request.OpenTime)
		if err != nil {
			return nil, fmt.Errorf("invalid open_time format: %w", err)
		}
		wasteBank.OpenTime = openTime
	}

	if request.CloseTime != nil && *request.CloseTime != "" {
		closeTime, err := time.Parse("15:04:05", *request.CloseTime)
		if err != nil {
			return nil, fmt.Errorf("invalid close_time format: %w", err)
		}
		wasteBank.CloseTime = closeTime
	}

	if err := c.WasteBankRepository.Update(tx, wasteBank); err != nil {
		c.Log.Warnf("Failed to update waste bank: %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteBankToResponse(wasteBank), nil
}

func (c *WasteBankUsecase) Delete(ctx context.Context, request *model.DeleteWasteBankRequest) (*model.WasteBankResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find waste bank by id
	wasteBank := new(entity.WasteBankProfile)
	if err := c.WasteBankRepository.FindById(tx, wasteBank, request.ID); err != nil {
		c.Log.Warnf("Failed find waste bank by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete waste bank
	if err := c.WasteBankRepository.Delete(tx, wasteBank); err != nil {
		c.Log.Warnf("Failed delete waste bank : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteBankToResponse(wasteBank), nil
}
