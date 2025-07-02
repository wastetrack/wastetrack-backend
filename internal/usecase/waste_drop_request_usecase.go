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
	DB                             *gorm.DB
	Log                            *logrus.Logger
	Validate                       *validator.Validate
	WasteDropRequestRepository     *repository.WasteDropRequestRepository
	UserRepository                 *repository.UserRepository
	WasteTypeRepository            *repository.WasteTypeRepository
	WasteDropRequestItemRepository *repository.WasteDropRequestItemRepository
	WasteBankPricedTypeRepository  *repository.WasteBankPricedTypeRepository
	CustomerRepository             *repository.CustomerRepository
	WasteBankRepository            *repository.WasteBankRepository
	WasteCollectorRepository       *repository.WasteCollectorRepository
}

func NewWasteDropRequestUsecase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, wasteDropRequestRepository *repository.WasteDropRequestRepository, userRepository *repository.UserRepository, wasteTypeRepository *repository.WasteTypeRepository, wasteDropRequestItemRepository *repository.WasteDropRequestItemRepository, wasteBankPricedTypeRepository *repository.WasteBankPricedTypeRepository, customerRepository *repository.CustomerRepository, wasteBankRepository *repository.WasteBankRepository, wasteCollectorRepository *repository.WasteCollectorRepository) *WasteDropRequestUsecase {
	return &WasteDropRequestUsecase{
		DB:                             db,
		Log:                            log,
		Validate:                       validate,
		WasteDropRequestRepository:     wasteDropRequestRepository,
		UserRepository:                 userRepository,
		WasteTypeRepository:            wasteTypeRepository,
		WasteDropRequestItemRepository: wasteDropRequestItemRepository,
		WasteBankPricedTypeRepository:  wasteBankPricedTypeRepository,
		CustomerRepository:             customerRepository,
		WasteBankRepository:            wasteBankRepository,
		WasteCollectorRepository:       wasteCollectorRepository,
	}
}

func (c *WasteDropRequestUsecase) Create(ctx context.Context, request *model.WasteDropRequestRequest) (*model.WasteDropRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Validate items arrays have same length
	if len(request.Items.WasteTypeIDs) != len(request.Items.Quantities) {
		c.Log.Warnf("WasteTypeIDs and Quantities arrays must have same length")
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

	// Validate waste type IDs
	wasteTypeIDs := make([]uuid.UUID, len(request.Items.WasteTypeIDs))
	for i, wasteTypeIDStr := range request.Items.WasteTypeIDs {
		wasteTypeID, err := uuid.Parse(wasteTypeIDStr)
		if err != nil {
			c.Log.Warnf("Invalid waste type ID: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		wasteTypeIDs[i] = wasteTypeID
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

	// Validate all waste types exist
	for _, wasteTypeID := range wasteTypeIDs {
		wasteType := new(entity.WasteType)
		if err := c.WasteTypeRepository.FindById(tx, wasteType, wasteTypeID.String()); err != nil {
			c.Log.Warnf("Failed to find waste type by ID %s: %+v", wasteTypeID.String(), err)
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

	// Create waste drop request items in batch
	wasteDropRequestItems := make([]*entity.WasteDropRequestItem, len(wasteTypeIDs))
	for i, wasteTypeID := range wasteTypeIDs {
		wasteDropRequestItems[i] = &entity.WasteDropRequestItem{
			RequestID:        wasteDropRequest.ID,
			WasteTypeID:      wasteTypeID,
			Quantity:         request.Items.Quantities[i],
			VerifiedWeight:   0.0, // Initial values
			VerifiedSubtotal: 0,   // Initial values
		}
	}

	if err := c.WasteDropRequestItemRepository.CreateBatch(tx, wasteDropRequestItems); err != nil {
		c.Log.Warnf("Failed to create waste drop request items: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteDropRequestToSimpleResponse(wasteDropRequest), nil
}

// Helper method to update customer profile with environmental impact
func (c *WasteDropRequestUsecase) updateCustomerProfile(tx *gorm.DB, customerID uuid.UUID, totalWeight float64, itemCount int64) error {
	// Find or create customer profile
	customerProfile := &entity.CustomerProfile{}
	err := c.CustomerRepository.FindByUserID(tx, customerProfile, customerID.String())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new profile if doesn't exist
			customerProfile = &entity.CustomerProfile{
				UserID:        customerID,
				CarbonDeficit: 0,
				WaterSaved:    0,
				BagsStored:    0,
				Trees:         0,
			}
			if err := c.CustomerRepository.Create(tx, customerProfile); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Calculate environmental impact based on waste weight
	// These are example calculations - adjust based on your environmental impact formulas
	carbonReduction := int64(totalWeight * 2.5) // 2.5kg CO2 reduction per kg waste
	waterSaved := int64(totalWeight * 1000)     // 1 liter per gram of waste
	bagsStored := itemCount                     // Each item represents a bag
	treesSaved := int64(totalWeight / 10)       // 1 tree saved per 10kg of waste

	// Update profile with accumulated values
	customerProfile.CarbonDeficit += carbonReduction
	customerProfile.WaterSaved += waterSaved
	customerProfile.BagsStored += bagsStored
	customerProfile.Trees += treesSaved

	return c.CustomerRepository.Update(tx, customerProfile)
}

// Helper method to update waste collector profile
func (c *WasteDropRequestUsecase) updateWasteCollectorProfile(tx *gorm.DB, collectorID uuid.UUID, totalWeight float64) error {
	c.Log.Infof("Updating waste collector profile for ID: %s with weight: %f", collectorID.String(), totalWeight)

	// Find or create waste collector profile
	collectorProfile := &entity.WasteCollectorProfile{}
	err := c.WasteCollectorRepository.FindByUserID(tx, collectorProfile, collectorID.String())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Infof("Creating new waste collector profile for user ID: %s", collectorID.String())
			// Create new profile if doesn't exist
			collectorProfile = &entity.WasteCollectorProfile{
				UserID:           collectorID,
				TotalWasteWeight: 0,
			}
			if err := c.WasteCollectorRepository.Create(tx, collectorProfile); err != nil {
				c.Log.Warnf("Failed to create waste collector profile: %+v", err)
				return err
			}
		} else {
			c.Log.Warnf("Failed to find waste collector profile: %+v", err)
			return err
		}
	}

	// Update total waste weight
	oldWeight := collectorProfile.TotalWasteWeight
	collectorProfile.TotalWasteWeight += totalWeight
	c.Log.Infof("Updating collector weight from %f to %f", oldWeight, collectorProfile.TotalWasteWeight)

	if err := c.WasteCollectorRepository.Update(tx, collectorProfile); err != nil {
		c.Log.Warnf("Failed to update waste collector profile: %+v", err)
		return err
	}

	c.Log.Infof("Successfully updated waste collector profile")
	return nil
}

// Helper method to update waste bank profile
func (c *WasteDropRequestUsecase) updateWasteBankProfile(tx *gorm.DB, wasteBankID uuid.UUID, totalWeight float64) error {
	c.Log.Infof("Updating waste bank profile for ID: %s with weight: %f", wasteBankID.String(), totalWeight)

	// Find or create waste bank profile
	wasteBankProfile := &entity.WasteBankProfile{}
	err := c.WasteBankRepository.FindByUserID(tx, wasteBankProfile, wasteBankID.String())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Infof("Creating new waste bank profile for user ID: %s", wasteBankID.String())
			// Create new profile if doesn't exist
			wasteBankProfile = &entity.WasteBankProfile{
				UserID:           wasteBankID,
				TotalWasteWeight: 0,
				TotalWorkers:     1, // Default to 1 worker
			}
			if err := c.WasteBankRepository.Create(tx, wasteBankProfile); err != nil {
				c.Log.Warnf("Failed to create waste bank profile: %+v", err)
				return err
			}
		} else {
			c.Log.Warnf("Failed to find waste bank profile: %+v", err)
			return err
		}
	}

	// Update total waste weight
	oldWeight := wasteBankProfile.TotalWasteWeight
	wasteBankProfile.TotalWasteWeight += totalWeight
	c.Log.Infof("Updating waste bank weight from %f to %f", oldWeight, wasteBankProfile.TotalWasteWeight)

	if err := c.WasteBankRepository.Update(tx, wasteBankProfile); err != nil {
		c.Log.Warnf("Failed to update waste bank profile: %+v", err)
		return err
	}

	c.Log.Infof("Successfully updated waste bank profile")
	return nil
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

func (c *WasteDropRequestUsecase) Complete(ctx context.Context, request *model.CompleteWasteDropRequest) (*model.WasteDropRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	if len(request.Items.WasteTypeIDs) != len(request.Items.Weights) {
		c.Log.Warnf("WasteTypeIDs and Weights arrays must have same length")
		return nil, fiber.ErrBadRequest
	}

	// Find waste drop request
	wasteDropRequest := new(entity.WasteDropRequest)
	if err := c.WasteDropRequestRepository.FindByID(tx, wasteDropRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste drop request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	if wasteDropRequest.WasteBankID == nil {
		c.Log.Warn("Cannot complete request without assigned waste bank")
		return nil, fiber.ErrBadRequest
	}

	// Parse waste type IDs and build map
	weightMap := make(map[uuid.UUID]float64)
	for i, idStr := range request.Items.WasteTypeIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.Log.Warnf("Invalid waste type ID: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		weightMap[id] = request.Items.Weights[i]
	}

	// Get existing items
	var existingItems []entity.WasteDropRequestItem
	if err := tx.Where("request_id = ?", wasteDropRequest.ID).Find(&existingItems).Error; err != nil {
		c.Log.Warnf("Failed to find waste drop request items: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if len(existingItems) == 0 {
		c.Log.Warnf("No waste drop request items found for request ID: %s", wasteDropRequest.ID)
		return nil, fiber.ErrNotFound
	}

	var totalVerifiedPrice int64
	var totalVerifiedWeight float64

	for _, item := range existingItems {
		weight, exists := weightMap[item.WasteTypeID]
		if !exists {
			c.Log.Warnf("Weight not provided for waste type ID: %s", item.WasteTypeID)
			return nil, fiber.ErrBadRequest
		}

		// Get price from WasteBankPricedType
		searchReq := &model.SearchWasteBankPricedTypeRequest{
			WasteBankID: wasteDropRequest.WasteBankID.String(),
			WasteTypeID: item.WasteTypeID.String(),
			Page:        1,
			Size:        1,
		}

		pricedTypes, _, err := c.WasteBankPricedTypeRepository.Search(tx, searchReq)
		if err != nil || len(pricedTypes) == 0 {
			c.Log.Warnf("Price not found for waste bank %s and type %s: %+v",
				wasteDropRequest.WasteBankID, item.WasteTypeID, err)
			return nil, fiber.ErrNotFound
		}

		price := pricedTypes[0].CustomPricePerKgs
		subtotal := int64(weight * float64(price))

		item.VerifiedWeight = weight
		item.VerifiedSubtotal = subtotal

		if err := tx.Save(&item).Error; err != nil {
			c.Log.Warnf("Failed to update item: %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		totalVerifiedPrice += subtotal
		totalVerifiedWeight += weight
	}

	if len(request.Items.WasteTypeIDs) != len(existingItems) {
		c.Log.Warnf("Mismatch between request items (%d) and DB items (%d)",
			len(request.Items.WasteTypeIDs), len(existingItems))
		return nil, fiber.ErrBadRequest
	}

	// Update main request
	wasteDropRequest.Status = "completed"
	wasteDropRequest.TotalPrice = totalVerifiedPrice

	if err := c.WasteDropRequestRepository.Update(tx, wasteDropRequest); err != nil {
		c.Log.Warnf("Failed to update waste drop request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update related profiles
	if err := c.updateCustomerProfile(tx, wasteDropRequest.CustomerID, totalVerifiedWeight, int64(len(existingItems))); err != nil {
		c.Log.Warnf("Failed to update customer profile: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if wasteDropRequest.AssignedCollectorID != nil {
		if err := c.updateWasteCollectorProfile(tx, *wasteDropRequest.AssignedCollectorID, totalVerifiedWeight); err != nil {
			c.Log.Warnf("Failed to update collector profile: %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if wasteDropRequest.WasteBankID != nil {
		if err := c.updateWasteBankProfile(tx, *wasteDropRequest.WasteBankID, totalVerifiedWeight); err != nil {
			c.Log.Warnf("Failed to update waste bank profile: %+v", err)
			return nil, fiber.ErrInternalServerError
		}
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
