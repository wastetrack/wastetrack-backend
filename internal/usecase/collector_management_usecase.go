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

type CollectorManagementUsecase struct {
	DB                            *gorm.DB
	Log                           *logrus.Logger
	Validate                      *validator.Validate
	CollectorManagementRepository *repository.CollectorManagementRepository
	UserRepository                *repository.UserRepository
}

func NewCollectorManagementUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, collectorManagementRepository *repository.CollectorManagementRepository, userRepository *repository.UserRepository) *CollectorManagementUsecase {
	return &CollectorManagementUsecase{
		DB:                            db,
		Log:                           log,
		Validate:                      validate,
		CollectorManagementRepository: collectorManagementRepository,
		UserRepository:                userRepository,
	}
}

func (u *CollectorManagementUsecase) Create(ctx context.Context, request *model.CollectorManagementRequest) (*model.CollectorManagementSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Parse UUIDs
	wasteBankID, err := uuid.Parse(request.WasteBankID)
	if err != nil {
		u.Log.Warnf("Invalid waste bank ID: %v", err)
		return nil, fiber.ErrBadRequest
	}

	collectorID, err := uuid.Parse(request.CollectorID)
	if err != nil {
		u.Log.Warnf("Invalid collector ID: %v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if waste bank exists
	wasteBank := new(entity.User)
	if err := u.UserRepository.FindById(tx, wasteBank, request.WasteBankID); err != nil {
		u.Log.Warnf("Waste bank not found: %v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Waste bank not found")
	}

	// Check if collector exists
	collector := new(entity.User)
	if err := u.UserRepository.FindById(tx, collector, request.CollectorID); err != nil {
		u.Log.Warnf("Collector not found: %v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "Collector not found")
	}

	// Check for duplicate combination
	existing := new(entity.CollectorManagement)
	if err := u.CollectorManagementRepository.FindByWasteBankAndCollector(tx, existing, request.WasteBankID, request.CollectorID); err == nil {
		u.Log.Warnf("Duplicate collector management found: waste_bank_id=%s, collector_id=%s", request.WasteBankID, request.CollectorID)
		return nil, fiber.NewError(fiber.StatusConflict, "This collector is already assigned to this waste bank")
	} else if err != gorm.ErrRecordNotFound {
		u.Log.Warnf("Error checking for duplicate: %v", err)
		return nil, fiber.ErrInternalServerError
	}

	collectorManagement := &entity.CollectorManagement{
		WasteBankID: wasteBankID,
		CollectorID: collectorID,
		Status:      request.Status,
	}

	if err := u.CollectorManagementRepository.Create(tx, collectorManagement); err != nil {
		u.Log.Warnf("Failed to create collector management: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CollectorManagementToSimpleResponse(collectorManagement), nil
}

func (u *CollectorManagementUsecase) Get(ctx context.Context, id string) (*model.CollectorManagementResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	collectorManagement := new(entity.CollectorManagement)
	if err := u.CollectorManagementRepository.FindByIdWithRelations(tx, collectorManagement, id); err != nil {
		u.Log.Warnf("Collector management not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CollectorManagementToResponse(collectorManagement), nil
}

func (u *CollectorManagementUsecase) Update(ctx context.Context, request *model.UpdateCollectorManagementRequest) (*model.CollectorManagementSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	collectorManagement := new(entity.CollectorManagement)
	if err := u.CollectorManagementRepository.FindById(tx, collectorManagement, request.ID); err != nil {
		u.Log.Warnf("Collector management not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	// Update waste bank ID if provided
	if request.WasteBankID != "" {
		wasteBankID, err := uuid.Parse(request.WasteBankID)
		if err != nil {
			u.Log.Warnf("Invalid waste bank ID: %v", err)
			return nil, fiber.ErrBadRequest
		}

		wasteBank := new(entity.User)
		if err := u.UserRepository.FindById(tx, wasteBank, request.WasteBankID); err != nil {
			u.Log.Warnf("Waste bank not found: %v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, "Waste bank not found")
		}
		collectorManagement.WasteBankID = wasteBankID
	}

	// Update collector ID if provided
	if request.CollectorID != "" {
		collectorID, err := uuid.Parse(request.CollectorID)
		if err != nil {
			u.Log.Warnf("Invalid collector ID: %v", err)
			return nil, fiber.ErrBadRequest
		}

		collector := new(entity.User)
		if err := u.UserRepository.FindById(tx, collector, request.CollectorID); err != nil {
			u.Log.Warnf("Collector not found: %v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, "Collector not found")
		}
		collectorManagement.CollectorID = collectorID
	}

	// Check for duplicate combination if either waste bank or collector is being updated
	if request.WasteBankID != "" || request.CollectorID != "" {
		existing := new(entity.CollectorManagement)
		finalWasteBankID := collectorManagement.WasteBankID.String()
		finalCollectorID := collectorManagement.CollectorID.String()

		if err := u.CollectorManagementRepository.FindByWasteBankAndCollectorExcludeID(tx, existing, finalWasteBankID, finalCollectorID, request.ID); err == nil {
			u.Log.Warnf("Duplicate collector management found during update: waste_bank_id=%s, collector_id=%s", finalWasteBankID, finalCollectorID)
			return nil, fiber.NewError(fiber.StatusConflict, "This collector is already assigned to this waste bank")
		} else if err != gorm.ErrRecordNotFound {
			u.Log.Warnf("Error checking for duplicate during update: %v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	// Update status if provided
	if request.Status != "" {
		collectorManagement.Status = request.Status
	}

	if err := u.CollectorManagementRepository.Update(tx, collectorManagement); err != nil {
		u.Log.Warnf("Update failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CollectorManagementToSimpleResponse(collectorManagement), nil
}

func (u *CollectorManagementUsecase) Search(ctx context.Context, request *model.SearchCollectorManagementRequest) ([]model.CollectorManagementSimpleResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	// Use search without preloads for simple response
	collectorManagements, total, err := u.CollectorManagementRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Warn("Search failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("Commit failed")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.CollectorManagementSimpleResponse, len(collectorManagements))
	for i, cm := range collectorManagements {
		responses[i] = *converter.CollectorManagementToSimpleResponse(&cm)
	}

	return responses, total, nil
}

func (u *CollectorManagementUsecase) Delete(ctx context.Context, id string) (*model.CollectorManagementSimpleResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	collectorManagement := new(entity.CollectorManagement)
	if err := u.CollectorManagementRepository.FindById(tx, collectorManagement, id); err != nil {
		u.Log.Warnf("Collector management not found: %v", err)
		return nil, fiber.ErrNotFound
	}

	if err := u.CollectorManagementRepository.Delete(tx, collectorManagement); err != nil {
		u.Log.Warnf("Delete failed: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Commit error: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CollectorManagementToSimpleResponse(collectorManagement), nil
}
