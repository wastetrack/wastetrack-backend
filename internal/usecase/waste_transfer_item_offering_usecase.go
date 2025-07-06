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

type WasteTransferItemOfferingUsecase struct {
	DB                                  *gorm.DB
	Log                                 *logrus.Logger
	Validate                            *validator.Validate
	WasteTransferItemOfferingRepository *repository.WasteTransferItemOfferingRepository
	WasteTransferRequestRepository      *repository.WasteTransferRequestRepository
	WasteTypeRepository                 *repository.WasteTypeRepository
}

func NewWasteTransferItemOfferingUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	wasteTransferItemOfferingRepository *repository.WasteTransferItemOfferingRepository,
	wasteTransferRequestRepository *repository.WasteTransferRequestRepository,
	wasteTypeRepository *repository.WasteTypeRepository,
) *WasteTransferItemOfferingUsecase {
	return &WasteTransferItemOfferingUsecase{
		DB:                                  db,
		Log:                                 log,
		Validate:                            validate,
		WasteTransferItemOfferingRepository: wasteTransferItemOfferingRepository,
		WasteTransferRequestRepository:      wasteTransferRequestRepository,
		WasteTypeRepository:                 wasteTypeRepository,
	}
}

func (c *WasteTransferItemOfferingUsecase) Get(ctx context.Context, request *model.GetWasteTransferItemOfferingRequest) (*model.WasteTransferItemOfferingResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	item := new(entity.WasteTransferItemOffering)
	if err := c.WasteTransferItemOfferingRepository.FindByID(tx, item, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste transfer item offering by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferItemOfferingToResponse(item), nil
}

func (c *WasteTransferItemOfferingUsecase) Search(ctx context.Context, request *model.SearchWasteTransferItemOfferingRequest) ([]model.WasteTransferItemOfferingSimpleResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	items, total, err := c.WasteTransferItemOfferingRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to search waste transfer item offerings")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteTransferItemOfferingSimpleResponse, len(items))
	for i, item := range items {
		responses[i] = *converter.WasteTransferItemOfferingToSimpleResponse(&item)
	}

	return responses, total, nil
}
