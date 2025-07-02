package usecase

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"github.com/wastetrack/wastetrack-backend/internal/model/converter"
	"github.com/wastetrack/wastetrack-backend/internal/repository"
	"gorm.io/gorm"
)

type WasteDropRequestItemUsecase struct {
	DB                             *gorm.DB
	Log                            *logrus.Logger
	Validate                       *validator.Validate
	WasteDropRequestItemRepository *repository.WasteDropRequestItemRepository
	WasteDropRequestRepository     *repository.WasteDropRequestRepository
	WasteTypeRepository            *repository.WasteTypeRepository
}

func NewWasteDropRequestItemUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteDropRequestItemRepo *repository.WasteDropRequestItemRepository, wasteDropRequestRepo *repository.WasteDropRequestRepository, wasteTypeRepo *repository.WasteTypeRepository) *WasteDropRequestItemUsecase {
	return &WasteDropRequestItemUsecase{
		DB:                             db,
		Log:                            log,
		Validate:                       validate,
		WasteDropRequestItemRepository: wasteDropRequestItemRepo,
		WasteDropRequestRepository:     wasteDropRequestRepo,
		WasteTypeRepository:            wasteTypeRepo,
	}
}

func (c *WasteDropRequestItemUsecase) Create(ctx context.Context, request *model.WasteDropRequestItemRequest) (*model.WasteDropRequestItemSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Check if waste drop request exists
	wasteDropRequest := new(entity.WasteDropRequest)
	if err := c.WasteDropRequestRepository.FindById(tx, wasteDropRequest, request.RequestID); err != nil {
		c.Log.Warnf("Failed to find waste drop request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Check if waste type exists
	wasteType := new(entity.WasteType)
	if err := c.WasteTypeRepository.FindById(tx, wasteType, request.WasteTypeID); err != nil {
		c.Log.Warnf("Failed to find waste type by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	wasteDropRequestItem := &entity.WasteDropRequestItem{
		RequestID:        wasteDropRequest.ID,
		WasteTypeID:      wasteType.ID,
		Quantity:         request.Quantity,
		VerifiedWeight:   request.VerifiedWeight,
		VerifiedSubtotal: request.VerifiedSubtotal,
	}

	if err := c.WasteDropRequestItemRepository.Create(tx, wasteDropRequestItem); err != nil {
		c.Log.Warnf("Failed to create waste drop request item: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestItemToSimpleResponse(wasteDropRequestItem), nil
}

func (c *WasteDropRequestItemUsecase) Get(ctx context.Context, request *model.GetWasteDropRequestItemRequest) (*model.WasteDropRequestItemResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteDropRequestItem := new(entity.WasteDropRequestItem)
	if err := c.WasteDropRequestItemRepository.FindByID(tx, wasteDropRequestItem, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste drop request item by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestItemToResponse(wasteDropRequestItem), nil
}

func (c *WasteDropRequestItemUsecase) Update(ctx context.Context, request *model.UpdateWasteDropRequestItemRequest) (*model.WasteDropRequestItemSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteDropRequestItem := new(entity.WasteDropRequestItem)
	if err := c.WasteDropRequestItemRepository.FindById(tx, wasteDropRequestItem, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste drop request item by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Validate related entities if they are being updated
	if request.RequestID != "" {
		wasteDropRequest := new(entity.WasteDropRequest)
		if err := c.WasteDropRequestRepository.FindById(tx, wasteDropRequest, request.RequestID); err != nil {
			c.Log.Warnf("Failed to find waste drop request by ID: %+v", err)
			return nil, fiber.ErrNotFound
		}
		wasteDropRequestItem.RequestID = wasteDropRequest.ID
	}

	if request.WasteTypeID != "" {
		wasteType := new(entity.WasteType)
		if err := c.WasteTypeRepository.FindById(tx, wasteType, request.WasteTypeID); err != nil {
			c.Log.Warnf("Failed to find waste type by ID: %+v", err)
			return nil, fiber.ErrNotFound
		}
		wasteDropRequestItem.WasteTypeID = wasteType.ID
	}

	if request.Quantity > 0 {
		wasteDropRequestItem.Quantity = request.Quantity
	}
	if request.VerifiedWeight > 0 {
		wasteDropRequestItem.VerifiedWeight = request.VerifiedWeight
	}
	if request.VerifiedSubtotal > 0 {
		wasteDropRequestItem.VerifiedSubtotal = request.VerifiedSubtotal
	}

	if err := c.WasteDropRequestItemRepository.Update(tx, wasteDropRequestItem); err != nil {
		c.Log.Warnf("Failed to update waste drop request item: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestItemToSimpleResponse(wasteDropRequestItem), nil
}

func (c *WasteDropRequestItemUsecase) Search(ctx context.Context, request *model.SearchWasteDropRequestItemRequest) ([]model.WasteDropRequestItemSimpleResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	requestItems, total, err := c.WasteDropRequestItemRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to search waste drop request items")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteDropRequestItemSimpleResponse, len(requestItems))
	for i, item := range requestItems {
		responses[i] = *converter.WasteDropRequestItemToSimpleResponse(&item)
	}

	return responses, total, nil
}

func (c *WasteDropRequestItemUsecase) Delete(ctx context.Context, request *model.DeleteWasteDropRequestItemRequest) (*model.WasteDropRequestItemSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteDropRequestItem := new(entity.WasteDropRequestItem)
	if err := c.WasteDropRequestItemRepository.FindById(tx, wasteDropRequestItem, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste drop request item by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := c.WasteDropRequestItemRepository.Delete(tx, wasteDropRequestItem); err != nil {
		c.Log.Warnf("Failed to delete waste drop request item: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestItemToSimpleResponse(wasteDropRequestItem), nil
}
