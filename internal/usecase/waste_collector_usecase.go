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

type WasteCollectorUseCase struct {
	DB                       *gorm.DB
	Log                      *logrus.Logger
	Validate                 *validator.Validate
	WasteCollectorRepository *repository.WasteCollectorRepository
}

func NewWasteCollectorUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteCollectorRepository *repository.WasteCollectorRepository) *WasteCollectorUseCase {
	return &WasteCollectorUseCase{
		DB:                       db,
		Log:                      log,
		Validate:                 validate,
		WasteCollectorRepository: wasteCollectorRepository,
	}
}

// TODO: Implement Search

func (c *WasteCollectorUseCase) Get(ctx context.Context, request *model.GetWasteCollectorRequest) (*model.WasteCollectorResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	wasteCollector := new(entity.WasteCollectorProfile)
	if err := c.WasteCollectorRepository.FindByUserID(tx, wasteCollector, request.ID); err != nil {
		c.Log.Warnf("Failed find profile by user id : %+v", err)
		return nil, fiber.ErrNotFound
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCollectorToResponse(wasteCollector), nil

}

func (c *WasteCollectorUseCase) Update(ctx context.Context, request *model.UpdateWasteCollectorRequest, userID string, userRole string) (*model.WasteCollectorResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteCollector := new(entity.WasteCollectorProfile)
	if err := c.WasteCollectorRepository.FindById(tx, wasteCollector, request.ID); err != nil {
		c.Log.Warnf("Failed find subject by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Authorization check: Skip for admin, otherwise check ownership
	if userRole != "admin" && wasteCollector.UserID != uuid.MustParse(userID) {
		c.Log.Warnf("User %s is not authorized to update waste collector %s", userID, request.ID)
		return nil, fiber.ErrForbidden
	}

	if request.TotalWasteWeight != nil {
		wasteCollector.TotalWasteWeight = *request.TotalWasteWeight
	}

	if err := c.WasteCollectorRepository.Update(tx, wasteCollector); err != nil {
		c.Log.Warnf("Failed to update waste collector: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCollectorToResponse(wasteCollector), nil
}

func (c *WasteCollectorUseCase) Delete(ctx context.Context, request *model.DeleteWasteCollectorRequest) (*model.WasteCollectorResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Find waste collector by id
	wasteCollector := new(entity.WasteCollectorProfile)
	if err := c.WasteCollectorRepository.FindById(tx, wasteCollector, request.ID); err != nil {
		c.Log.Warnf("Failed find waste collector by id : %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Delete waste collector
	if err := c.WasteCollectorRepository.Delete(tx, wasteCollector); err != nil {
		c.Log.Warnf("Failed delete waste collector : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteCollectorToResponse(wasteCollector), nil
}
