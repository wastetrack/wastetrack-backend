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

type WasteBankPricedTypeUsecase struct {
	DB                            *gorm.DB
	Log                           *logrus.Logger
	Validate                      *validator.Validate
	WasteBankPricedTypeRepository *repository.WasteBankPricedTypeRepository
	WasteTypeRepository           *repository.WasteTypeRepository
}

func NewWasteBankPricedTypeUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	wasteBankPricedTypeRepo *repository.WasteBankPricedTypeRepository,
	wasteTypeRepo *repository.WasteTypeRepository,
) *WasteBankPricedTypeUsecase {
	return &WasteBankPricedTypeUsecase{
		DB: db, Log: log, Validate: validate,
		WasteBankPricedTypeRepository: wasteBankPricedTypeRepo,
		WasteTypeRepository:           wasteTypeRepo,
	}
}
func (uc *WasteBankPricedTypeUsecase) CreateBatch(ctx context.Context, requests []model.WasteBankPricedTypeRequest) ([]*model.WasteBankPricedTypeResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if len(requests) == 0 {
		return nil, fiber.ErrBadRequest
	}

	// Validate all items
	for _, req := range requests {
		if err := uc.Validate.Struct(req); err != nil {
			uc.Log.Warn("Validation error in batch request: ", err)
			return nil, fiber.ErrBadRequest
		}
	}

	// Optional: check all WasteType IDs exist (to avoid partial failure)
	for _, req := range requests {
		wasteType := new(entity.WasteType)
		if err := uc.WasteTypeRepository.FindById(tx, wasteType, req.WasteTypeID); err != nil {
			uc.Log.Warnf("Invalid WasteTypeID in batch: %s", req.WasteTypeID)
			return nil, fiber.ErrNotFound
		}
	}

	var entities []*entity.WasteBankPricedType
	for _, req := range requests {
		entities = append(entities, &entity.WasteBankPricedType{
			WasteBankID:       uuid.MustParse(req.WasteBankID),
			WasteTypeID:       uuid.MustParse(req.WasteTypeID),
			CustomPricePerKgs: req.CustomPricePerKgs,
		})
	}

	if err := uc.WasteBankPricedTypeRepository.CreateBatch(tx, entities); err != nil {
		uc.Log.Warn("Failed to batch insert WasteBankPricedType: ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("Failed to commit batch insert: ", err)
		return nil, fiber.ErrInternalServerError
	}

	// Convert to response
	var responses []*model.WasteBankPricedTypeResponse
	for _, e := range entities {
		responses = append(responses, converter.WasteBankPricedTypeToResponse(e))
	}

	return responses, nil
}
func (uc *WasteBankPricedTypeUsecase) Create(ctx context.Context, request *model.WasteBankPricedTypeRequest) (*model.WasteBankPricedTypeResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		uc.Log.Warn("Invalid request: ", err)
		return nil, fiber.ErrBadRequest
	}

	// check if waste type exists
	wasteType := new(entity.WasteType)
	if err := uc.WasteTypeRepository.FindById(tx, wasteType, request.WasteTypeID); err != nil {
		uc.Log.Warnf("Failed to find waste type by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	wpt := &entity.WasteBankPricedType{
		WasteBankID:       uuid.MustParse(request.WasteBankID),
		WasteTypeID:       uuid.MustParse(request.WasteTypeID),
		CustomPricePerKgs: request.CustomPricePerKgs,
	}

	if err := uc.WasteBankPricedTypeRepository.Create(tx, wpt); err != nil {
		uc.Log.Warn("Failed to create WasteBankPricedType: ", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("Failed to commit transaction: ", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteBankPricedTypeToResponse(wpt), nil
}

func (uc *WasteBankPricedTypeUsecase) Get(ctx context.Context, request *model.GetWasteBankPricedTypeRequest) (*model.WasteBankPricedTypeResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		return nil, fiber.ErrBadRequest
	}

	wpt := new(entity.WasteBankPricedType)
	if err := uc.WasteBankPricedTypeRepository.FindById(tx, wpt, request.ID); err != nil {
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteBankPricedTypeToResponse(wpt), nil
}

func (uc *WasteBankPricedTypeUsecase) Update(ctx context.Context, request *model.UpdateWasteBankPricedTypeRequest) (*model.WasteBankPricedTypeResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		return nil, fiber.ErrBadRequest
	}

	wpt := new(entity.WasteBankPricedType)
	if err := uc.WasteBankPricedTypeRepository.FindById(tx, wpt, request.ID); err != nil {
		return nil, fiber.ErrNotFound
	}

	wpt.CustomPricePerKgs = request.CustomPricePerKgs

	if err := uc.WasteBankPricedTypeRepository.Update(tx, wpt); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteBankPricedTypeToResponse(wpt), nil
}

func (uc *WasteBankPricedTypeUsecase) Delete(ctx context.Context, request *model.DeleteWasteBankPricedTypeRequest) (*model.WasteBankPricedTypeResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := uc.Validate.Struct(request); err != nil {
		return nil, fiber.ErrBadRequest
	}

	wpt := new(entity.WasteBankPricedType)
	if err := uc.WasteBankPricedTypeRepository.FindById(tx, wpt, request.ID); err != nil {
		return nil, fiber.ErrNotFound
	}

	if err := uc.WasteBankPricedTypeRepository.Delete(tx, wpt); err != nil {
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteBankPricedTypeToResponse(wpt), nil
}

func (c *WasteBankPricedTypeUsecase) Search(ctx context.Context, request *model.SearchWasteBankPricedTypeRequest) ([]model.WasteBankPricedTypeListResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warnf("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}
	wasteBankPricedTypes, total, err := c.WasteBankPricedTypeRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warnf("Failed to search waste bank priced types")
		return nil, 0, fiber.ErrInternalServerError
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteBankPricedTypeListResponse, len(wasteBankPricedTypes))
	for i, wasteBankPricedType := range wasteBankPricedTypes {
		responses[i] = *converter.WasteBankPricedTypeToListResponse(&wasteBankPricedType)
	}
	return responses, total, nil
}
