package usecase

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/model/converter"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"github.com/wastetrack/wastetrack-backend/internal/types"
	"gorm.io/gorm"
)

type WasteDropRequestUsecase struct {
	DB                         *gorm.DB
	Log                        *logrus.Logger
	Validate                   *validator.Validate
	WasteDropRequestRepository *repository.WasteDropRequestRepository
	UserRepository             *repository.UserRepository
}

func NewWasteDropRequestUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteDropRequestRepository *repository.WasteDropRequestRepository, userRepository *repository.UserRepository) *WasteDropRequestUsecase {
	return &WasteDropRequestUsecase{
		DB:                         db,
		Log:                        log,
		Validate:                   validate,
		WasteDropRequestRepository: wasteDropRequestRepository,
		UserRepository:             userRepository,
	}
}

func (c *WasteDropRequestUsecase) Create(ctx context.Context, request *model.WasteDropRequestRequest) (*model.WasteDropRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Parse UUIDs
	customerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		c.Log.Warnf("Invalid customer ID: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	var wasteBankID *uuid.UUID
	if request.WasteBankID != "" {
		id, err := uuid.Parse(request.WasteBankID)
		if err != nil {
			c.Log.Warnf("Invalid waste bank ID: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		wasteBankID = &id
	}

	// Check if customer exists
	customer := new(entity.User)
	if err := c.UserRepository.FindById(tx, customer, request.CustomerID); err != nil {
		c.Log.Warnf("Failed to find customer by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Check if waste bank exists (if provided)
	if wasteBankID != nil {
		wasteBank := new(entity.User)
		if err := c.UserRepository.FindById(tx, wasteBank, request.WasteBankID); err != nil {
			c.Log.Warnf("Failed to find waste bank by ID: %+v", err)
			return nil, fiber.ErrNotFound
		}
	}

	// Parse appointment date and times
	appointmentDate, err := time.Parse("2006-01-02", request.AppointmentDate)
	if err != nil {
		c.Log.Warnf("Invalid appointment date format: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	appointmentStartTime, err := time.Parse("15:04:05Z07:00", request.AppointmentStartTime)
	if err != nil {
		c.Log.Warnf("Invalid appointment start time format: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	appointmentEndTime, err := time.Parse("15:04:05Z07:00", request.AppointmentEndTime)
	if err != nil {
		c.Log.Warnf("Invalid appointment end time format: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteDropRequest := &entity.WasteDropRequest{
		DeliveryType:         request.DeliveryType,
		CustomerID:           customerID,
		UserPhoneNumber:      request.UserPhoneNumber,
		WasteBankID:          wasteBankID,
		TotalPrice:           request.TotalPrice,
		ImageURL:             request.ImageURL,
		Status:               "pending",
		AppointmentDate:      appointmentDate,
		AppointmentStartTime: types.NewTimeOnly(appointmentStartTime),
		AppointmentEndTime:   types.NewTimeOnly(appointmentEndTime),
		Notes:                request.Notes,
	}

	// Handle appointment location if provided
	if request.AppointmentLocation != nil {
		wasteDropRequest.AppointmentLocation = &types.Point{
			Lat: request.AppointmentLocation.Latitude,
			Lng: request.AppointmentLocation.Longitude,
		}
	}

	if err := c.WasteDropRequestRepository.Create(tx, wasteDropRequest); err != nil {
		c.Log.Warnf("Failed to create waste drop request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestToSimpleResponse(wasteDropRequest), nil
}

func (c *WasteDropRequestUsecase) Get(ctx context.Context, request *model.GetWasteDropRequest) (*model.WasteDropRequestResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteDropRequest := new(entity.WasteDropRequest)
	if err := c.WasteDropRequestRepository.FindByID(tx, wasteDropRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste drop request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestToResponse(wasteDropRequest), nil
}

func (c *WasteDropRequestUsecase) Update(ctx context.Context, request *model.UpdateWasteDropRequest) (*model.WasteDropRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteDropRequest := new(entity.WasteDropRequest)
	if err := c.WasteDropRequestRepository.FindByID(tx, wasteDropRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste drop request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Update fields if provided
	if request.DeliveryType != "" {
		wasteDropRequest.DeliveryType = request.DeliveryType
	}
	if request.Status != "" {
		wasteDropRequest.Status = request.Status
	}
	if request.AssignedCollectorID != "" {
		collectorID, err := uuid.Parse(request.AssignedCollectorID)
		if err != nil {
			c.Log.Warnf("Invalid collector ID: %+v", err)
			return nil, fiber.ErrBadRequest
		}

		// Check if collector exists
		collector := new(entity.User)
		if err := c.UserRepository.FindById(tx, collector, request.AssignedCollectorID); err != nil {
			c.Log.Warnf("Failed to find collector by ID: %+v", err)
			return nil, fiber.ErrNotFound
		}

		wasteDropRequest.AssignedCollectorID = &collectorID
	}

	if err := c.WasteDropRequestRepository.Update(tx, wasteDropRequest); err != nil {
		c.Log.Warnf("Failed to update waste drop request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestToSimpleResponse(wasteDropRequest), nil
}

func (c *WasteDropRequestUsecase) Search(ctx context.Context, request *model.SearchWasteDropRequest) ([]model.WasteDropRequestSimpleResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	wasteDropRequests, total, err := c.WasteDropRequestRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to search waste drop requests")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteDropRequestSimpleResponse, len(wasteDropRequests))
	for i, request := range wasteDropRequests {
		responses[i] = *converter.WasteDropRequestToSimpleResponse(&request)
	}

	return responses, total, nil
}

func (c *WasteDropRequestUsecase) Delete(ctx context.Context, request *model.DeleteWasteDropRequest) (*model.WasteDropRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteDropRequest := new(entity.WasteDropRequest)
	if err := c.WasteDropRequestRepository.FindByID(tx, wasteDropRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste drop request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := c.WasteDropRequestRepository.Delete(tx, wasteDropRequest); err != nil {
		c.Log.Warnf("Failed to delete waste drop request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestToSimpleResponse(wasteDropRequest), nil
}
