package usecase

import (
	"context"
	"fmt"
	"math"
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

type WasteTransferRequestUsecase struct {
	DB                                  *gorm.DB
	Log                                 *logrus.Logger
	Validate                            *validator.Validate
	WasteTransferRequestRepository      *repository.WasteTransferRequestRepository
	WasteTransferItemOfferingRepository *repository.WasteTransferItemOfferingRepository
	UserRepository                      *repository.UserRepository
	WasteTypeRepository                 *repository.WasteTypeRepository
}

func NewWasteTransferRequestUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	wasteTransferRequestRepository *repository.WasteTransferRequestRepository,
	wasteTransferItemOfferingRepository *repository.WasteTransferItemOfferingRepository,
	userRepository *repository.UserRepository,
	wasteTypeRepository *repository.WasteTypeRepository,
) *WasteTransferRequestUsecase {
	return &WasteTransferRequestUsecase{
		DB:                                  db,
		Log:                                 log,
		Validate:                            validate,
		WasteTransferRequestRepository:      wasteTransferRequestRepository,
		WasteTransferItemOfferingRepository: wasteTransferItemOfferingRepository,
		UserRepository:                      userRepository,
		WasteTypeRepository:                 wasteTypeRepository,
	}
}

func (c *WasteTransferRequestUsecase) Create(ctx context.Context, request *model.WasteTransferRequestRequest) (*model.WasteTransferRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Validate items arrays have same length
	if len(request.Items.WasteTypeIDs) != len(request.Items.OfferingWeights) ||
		len(request.Items.WasteTypeIDs) != len(request.Items.OfferingPricesPerKgs) {
		c.Log.Warnf("WasteTypeIDs, OfferingWeights, and OfferingPricesPerKgs arrays must have same length")
		return nil, fiber.ErrBadRequest
	}

	// Parse UUIDs
	sourceUserID, err := uuid.Parse(request.SourceUserID)
	if err != nil {
		c.Log.Warnf("Invalid source user ID: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	destinationUserID, err := uuid.Parse(request.DestinationUserID)
	if err != nil {
		c.Log.Warnf("Invalid destination user ID: %+v", err)
		return nil, fiber.ErrBadRequest
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

	// Check if source user exists
	sourceUser := new(entity.User)
	if err := c.UserRepository.FindById(tx, sourceUser, request.SourceUserID); err != nil {
		c.Log.Warnf("Failed to find source user by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Check if destination user exists
	destinationUser := new(entity.User)
	if err := c.UserRepository.FindById(tx, destinationUser, request.DestinationUserID); err != nil {
		c.Log.Warnf("Failed to find destination user by ID: %+v", err)
		return nil, fiber.ErrNotFound
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

	var appointmentStartTime, appointmentEndTime types.TimeOnly
	if request.AppointmentStartTime != "" {
		startTime, err := time.Parse("15:04:05Z07:00", request.AppointmentStartTime)
		if err != nil {
			c.Log.Warnf("Invalid appointment start time format: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		appointmentStartTime = types.NewTimeOnly(startTime)
	}

	if request.AppointmentEndTime != "" {
		endTime, err := time.Parse("15:04:05Z07:00", request.AppointmentEndTime)
		if err != nil {
			c.Log.Warnf("Invalid appointment end time format: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		appointmentEndTime = types.NewTimeOnly(endTime)
	}

	wasteTransferRequest := &entity.WasteTransferRequest{
		SourceUserID:           sourceUserID,
		DestinationUserID:      destinationUserID,
		FormType:               request.FormType,
		TotalWeight:            0, // Will be calculated from items
		TotalPrice:             0, // Will be calculated from items
		Status:                 "pending",
		ImageURL:               request.ImageURL,
		Notes:                  request.Notes,
		SourcePhoneNumber:      request.SourcePhoneNumber,
		DestinationPhoneNumber: request.DestinationPhoneNumber,
		AppointmentDate:        appointmentDate,
		AppointmentStartTime:   appointmentStartTime,
		AppointmentEndTime:     appointmentEndTime,
	}
	// Handle appointment location if provided
	if request.AppointmentLocation != nil {
		wasteTransferRequest.AppointmentLocation = &types.Point{
			Lat: request.AppointmentLocation.Latitude,
			Lng: request.AppointmentLocation.Longitude,
		}
	}

	if err := c.WasteTransferRequestRepository.Create(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to create waste transfer request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create waste transfer item offerings in batch
	var totalOfferingWeight float64
	var totalOfferingPrice float64
	wasteTransferItems := make([]*entity.WasteTransferItemOffering, len(wasteTypeIDs))
	for i, wasteTypeID := range wasteTypeIDs {
		weight := request.Items.OfferingWeights[i]
		pricePerKg := request.Items.OfferingPricesPerKgs[i]

		wasteTransferItems[i] = &entity.WasteTransferItemOffering{
			TransferFormID:      wasteTransferRequest.ID,
			WasteTypeID:         wasteTypeID,
			OfferingWeight:      weight,
			OfferingPricePerKgs: pricePerKg,
			AcceptedWeight:      0, // Initial values
			AcceptedPricePerKgs: 0, // Initial values
		}

		totalOfferingWeight += weight
		totalOfferingPrice += weight * pricePerKg
	}

	if err := c.WasteTransferItemOfferingRepository.CreateBatch(tx, wasteTransferItems); err != nil {
		c.Log.Warnf("Failed to create waste transfer item offerings: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update total weight and price
	wasteTransferRequest.TotalWeight = totalOfferingWeight
	wasteTransferRequest.TotalPrice = int64(totalOfferingPrice)

	if err := c.WasteTransferRequestRepository.Update(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to update waste transfer request totals: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferRequestToSimpleResponse(wasteTransferRequest), nil
}
func (c *WasteTransferRequestUsecase) Get(ctx context.Context, request *model.GetWasteTransferRequest) (*model.WasteTransferRequestResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteTransferRequest := new(entity.WasteTransferRequest)

	// Use FindByIDWithDistance if coordinates are provided, otherwise use FindByID
	if request.Latitude != nil && request.Longitude != nil {
		if err := c.WasteTransferRequestRepository.FindByIDWithDistance(tx, wasteTransferRequest, request.ID, request.Latitude, request.Longitude); err != nil {
			c.Log.Warnf("Failed to find waste transfer request by ID with distance: %+v", err)
			return nil, fiber.ErrNotFound
		}
	} else {
		if err := c.WasteTransferRequestRepository.FindByID(tx, wasteTransferRequest, request.ID); err != nil {
			c.Log.Warnf("Failed to find waste transfer request by ID: %+v", err)
			return nil, fiber.ErrNotFound
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferRequestToResponse(wasteTransferRequest), nil
}

func (c *WasteTransferRequestUsecase) AssignCollectorByWasteType(ctx context.Context, request *model.AssignCollectorByWasteTypeRequest) (*model.WasteTransferRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Parse and validate the transfer request ID
	transferFormUUID, err := uuid.Parse(request.ID)
	if err != nil {
		c.Log.Warnf("Invalid transfer request ID: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	var collectorID *uuid.UUID
	if request.AssignedCollectorID != "" {
		parsedUUID, err := uuid.Parse(request.AssignedCollectorID)
		if err != nil {
			c.Log.Warnf("Invalid collector ID: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		collectorID = &parsedUUID

		// Validate that the collector exists
		collector := new(entity.User)
		if err := c.UserRepository.FindById(tx, collector, request.AssignedCollectorID); err != nil {
			c.Log.Warnf("Collector not found: %v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, "Waste Collector not found")
		}
	}
	// Collector ID is optional - some industries can collect waste without a specific collector

	// Get transfer request
	wasteTransferRequest := new(entity.WasteTransferRequest)
	if err := c.WasteTransferRequestRepository.FindByID(tx, wasteTransferRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste transfer request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Validate status
	if wasteTransferRequest.Status != "pending" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Can only assign collector to pending requests")
	}

	// Get current items
	currentItems, err := c.WasteTransferItemOfferingRepository.FindByTransferFormID(tx, transferFormUUID)
	if err != nil {
		c.Log.Warnf("Failed to find waste transfer items: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if len(currentItems) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No items found for this transfer request")
	}

	// Create waste type pricing map
	wasteTypePricing := make(map[uuid.UUID]model.AssignCollectorWasteTypeRequest)
	for _, wt := range request.WasteTypes {
		wasteTypeID, err := uuid.Parse(wt.WasteTypeID)
		if err != nil {
			c.Log.Warnf("Invalid waste type ID: %s", wt.WasteTypeID)
			return nil, fiber.ErrBadRequest
		}

		// Validate pricing values
		if wt.AcceptedWeight < 0 || wt.AcceptedPricePerKgs < 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "Weight and price must be non-negative")
		}

		wasteTypePricing[wasteTypeID] = wt
	}

	// Update items based on waste type pricing
	var totalAcceptedWeight float64
	var totalAcceptedPrice float64

	for _, item := range currentItems {
		if pricing, exists := wasteTypePricing[item.WasteTypeID]; exists {
			// Validate accepted weight doesn't exceed offered weight
			if pricing.AcceptedWeight > item.OfferingWeight {
				return nil, fiber.NewError(fiber.StatusBadRequest,
					fmt.Sprintf("Accepted weight (%.2f) cannot exceed offered weight (%.2f) for waste type: %s",
						pricing.AcceptedWeight, item.OfferingWeight, item.WasteTypeID))
			}

			// Apply the waste type pricing to this item
			item.AcceptedWeight = pricing.AcceptedWeight
			item.AcceptedPricePerKgs = pricing.AcceptedPricePerKgs

			if err := c.WasteTransferItemOfferingRepository.Update(tx, &item); err != nil {
				c.Log.Warnf("Failed to update waste transfer item: %+v", err)
				return nil, fiber.ErrInternalServerError
			}

			totalAcceptedWeight += pricing.AcceptedWeight
			totalAcceptedPrice += pricing.AcceptedWeight * pricing.AcceptedPricePerKgs
		} else {
			return nil, fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("Missing pricing for waste type: %s", item.WasteTypeID))
		}
	}

	// Assign collector if one is provided
	if collectorID != nil {
		if err := c.WasteTransferRequestRepository.AssignCollector(tx, request.ID, *collectorID); err != nil {
			c.Log.Warnf("Failed to assign collector: %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	// Update the waste transfer request
	wasteTransferRequest.AssignedCollectorID = collectorID
	wasteTransferRequest.Status = "assigned"
	// Use proper rounding for float to int conversion
	wasteTransferRequest.TotalWeight = totalAcceptedWeight
	wasteTransferRequest.TotalPrice = int64(math.Round(totalAcceptedPrice))

	if err := c.WasteTransferRequestRepository.Update(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to update waste transfer request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferRequestToSimpleResponse(wasteTransferRequest), nil
}

func (c *WasteTransferRequestUsecase) Update(ctx context.Context, request *model.UpdateWasteTransferRequest) (*model.WasteTransferRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteTransferRequest := new(entity.WasteTransferRequest)
	if err := c.WasteTransferRequestRepository.FindByID(tx, wasteTransferRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste transfer request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Update fields if provided
	if request.FormType != "" {
		wasteTransferRequest.FormType = request.FormType
	}
	if request.Status != "" {
		wasteTransferRequest.Status = request.Status
	}
	if request.AppointmentDate != "" {
		appointmentDate, err := time.Parse("2006-01-02", request.AppointmentDate)
		if err != nil {
			c.Log.Warnf("Invalid appointment date format: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		wasteTransferRequest.AppointmentDate = appointmentDate
	}
	if request.AppointmentStartTime != "" {
		startTime, err := time.Parse("15:04:05Z07:00", request.AppointmentStartTime)
		if err != nil {
			c.Log.Warnf("Invalid appointment start time format: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		wasteTransferRequest.AppointmentStartTime = types.NewTimeOnly(startTime)
	}
	if request.AppointmentEndTime != "" {
		endTime, err := time.Parse("15:04:05Z07:00", request.AppointmentEndTime)
		if err != nil {
			c.Log.Warnf("Invalid appointment end time format: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		wasteTransferRequest.AppointmentEndTime = types.NewTimeOnly(endTime)
	}

	if err := c.WasteTransferRequestRepository.Update(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to update waste transfer request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferRequestToSimpleResponse(wasteTransferRequest), nil
}
func (c *WasteTransferRequestUsecase) CompleteRequest(ctx context.Context, request *model.CompleteWasteTransferRequest) (*model.WasteTransferRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Validate items arrays have same length
	if len(request.Items.WasteTypeIDs) != len(request.Items.Weights) {
		c.Log.Warnf("WasteTypeIDs and Weights arrays must have same length")
		return nil, fiber.ErrBadRequest
	}

	// Parse and validate the transfer request ID
	transferFormUUID, err := uuid.Parse(request.ID)
	if err != nil {
		c.Log.Warnf("Invalid transfer request ID: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Get transfer request
	wasteTransferRequest := new(entity.WasteTransferRequest)
	if err := c.WasteTransferRequestRepository.FindByID(tx, wasteTransferRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste transfer request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Validate status - can only complete assigned or in_progress requests
	if wasteTransferRequest.Status != "assigned" && wasteTransferRequest.Status != "in_progress" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Can only complete assigned or in_progress requests")
	}

	// Get current items
	currentItems, err := c.WasteTransferItemOfferingRepository.FindByTransferFormID(tx, transferFormUUID)
	if err != nil {
		c.Log.Warnf("Failed to find waste transfer items: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if len(currentItems) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No items found for this transfer request")
	}

	// Parse waste type IDs and validate
	wasteTypeUUIDs := make([]uuid.UUID, len(request.Items.WasteTypeIDs))
	for i, wasteTypeIDStr := range request.Items.WasteTypeIDs {
		wasteTypeID, err := uuid.Parse(wasteTypeIDStr)
		if err != nil {
			c.Log.Warnf("Invalid waste type ID: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		wasteTypeUUIDs[i] = wasteTypeID
	}

	// Validate all weights are non-negative
	for i, weight := range request.Items.Weights {
		if weight < 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("Weight at index %d must be non-negative", i))
		}
	}

	// Create waste type weight mapping
	wasteTypeWeights := make(map[uuid.UUID]float64)
	for i, wasteTypeID := range wasteTypeUUIDs {
		wasteTypeWeights[wasteTypeID] = request.Items.Weights[i]
	}

	// Update items with verified weights
	var totalVerifiedWeight float64
	var totalVerifiedPrice float64
	updatedItemsCount := 0

	for _, item := range currentItems {
		if verifiedWeight, exists := wasteTypeWeights[item.WasteTypeID]; exists {
			// Validate verified weight doesn't exceed accepted weight (if accepted weight is set)
			if item.AcceptedWeight > 0 && verifiedWeight > item.AcceptedWeight {
				return nil, fiber.NewError(fiber.StatusBadRequest,
					fmt.Sprintf("Verified weight (%.2f) cannot exceed accepted weight (%.2f) for waste type: %s",
						verifiedWeight, item.AcceptedWeight, item.WasteTypeID))
			}

			// Update the verified weight
			item.VerifiedWeight = verifiedWeight

			if err := c.WasteTransferItemOfferingRepository.Update(tx, &item); err != nil {
				c.Log.Warnf("Failed to update waste transfer item: %+v", err)
				return nil, fiber.ErrInternalServerError
			}

			totalVerifiedWeight += verifiedWeight
			// Use accepted price per kg for calculation, fallback to offering price if not set
			pricePerKg := item.AcceptedPricePerKgs
			if pricePerKg == 0 {
				pricePerKg = item.OfferingPricePerKgs
			}
			totalVerifiedPrice += verifiedWeight * pricePerKg
			updatedItemsCount++
		}
	}

	// Ensure all provided waste types were found and updated
	if updatedItemsCount != len(request.Items.WasteTypeIDs) {
		return nil, fiber.NewError(fiber.StatusBadRequest,
			"Some waste types not found in this transfer request")
	}

	// Update the waste transfer request
	wasteTransferRequest.Status = "completed"
	wasteTransferRequest.TotalWeight = totalVerifiedWeight
	wasteTransferRequest.TotalPrice = int64(math.Round(totalVerifiedPrice))

	if err := c.WasteTransferRequestRepository.Update(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to update waste transfer request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferRequestToSimpleResponse(wasteTransferRequest), nil
}
func (c *WasteTransferRequestUsecase) RecycleRequest(ctx context.Context, request *model.RecycleWasteTransferRequest) (*model.WasteTransferRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Validate items arrays have same length
	if len(request.Items.WasteTypeIDs) != len(request.Items.Weights) {
		c.Log.Warnf("WasteTypeIDs and Weights arrays must have same length")
		return nil, fiber.ErrBadRequest
	}

	// Parse and validate the transfer request ID
	transferFormUUID, err := uuid.Parse(request.ID)
	if err != nil {
		c.Log.Warnf("Invalid transfer request ID: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Get transfer request
	wasteTransferRequest := new(entity.WasteTransferRequest)
	if err := c.WasteTransferRequestRepository.FindByID(tx, wasteTransferRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste transfer request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Validate status - can only recycle completed requests
	if wasteTransferRequest.Status != "completed" && wasteTransferRequest.Status != "recycling_in_process" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Can only recycle completed or recycling in process requests")
	}

	// Get current items
	currentItems, err := c.WasteTransferItemOfferingRepository.FindByTransferFormID(tx, transferFormUUID)
	if err != nil {
		c.Log.Warnf("Failed to find waste transfer items: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if len(currentItems) == 0 {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No items found for this transfer request")
	}

	// Parse waste type IDs and validate
	wasteTypeUUIDs := make([]uuid.UUID, len(request.Items.WasteTypeIDs))
	for i, wasteTypeIDStr := range request.Items.WasteTypeIDs {
		wasteTypeID, err := uuid.Parse(wasteTypeIDStr)
		if err != nil {
			c.Log.Warnf("Invalid waste type ID: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		wasteTypeUUIDs[i] = wasteTypeID
	}

	// Validate all weights are non-negative
	for i, weight := range request.Items.Weights {
		if weight < 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("Weight at index %d must be non-negative", i))
		}
	}

	// Create waste type weight mapping
	wasteTypeWeights := make(map[uuid.UUID]float64)
	for i, wasteTypeID := range wasteTypeUUIDs {
		wasteTypeWeights[wasteTypeID] = request.Items.Weights[i]
	}

	// Update items with recycled weights
	var totalRecycledWeight float64
	var totalRecycledPrice float64
	updatedItemsCount := 0

	for _, item := range currentItems {
		if recycledWeight, exists := wasteTypeWeights[item.WasteTypeID]; exists {
			// Validate recycled weight doesn't exceed verified weight
			if item.VerifiedWeight > 0 && recycledWeight > item.VerifiedWeight {
				return nil, fiber.NewError(fiber.StatusBadRequest,
					fmt.Sprintf("Recycled weight (%.2f) cannot exceed verified weight (%.2f) for waste type: %s",
						recycledWeight, item.VerifiedWeight, item.WasteTypeID))
			}

			// Update the recycled weight
			item.RecycledWeight = recycledWeight

			if err := c.WasteTransferItemOfferingRepository.Update(tx, &item); err != nil {
				c.Log.Warnf("Failed to update waste transfer item: %+v", err)
				return nil, fiber.ErrInternalServerError
			}

			totalRecycledWeight += recycledWeight
			// Use accepted price per kg for calculation, fallback to offering price if not set
			pricePerKg := item.AcceptedPricePerKgs
			if pricePerKg == 0 {
				pricePerKg = item.OfferingPricePerKgs
			}
			totalRecycledPrice += recycledWeight * pricePerKg
			updatedItemsCount++
		}
	}

	// Ensure all provided waste types were found and updated
	if updatedItemsCount != len(request.Items.WasteTypeIDs) {
		return nil, fiber.NewError(fiber.StatusBadRequest,
			"Some waste types not found in this transfer request")
	}

	// Check that ALL items in the request now have recycled_weight > 0
	// This ensures the entire request is properly recycled before changing status
	// for _, item := range currentItems {
	// 	if item.RecycledWeight <= 0 {
	// 		return nil, fiber.NewError(fiber.StatusBadRequest,
	// 			fmt.Sprintf("All items must have recycled weight > 0. Item with waste type %s has recycled weight: %.2f",
	// 				item.WasteTypeID, item.RecycledWeight))
	// 	}
	// }

	// Update the waste transfer request
	wasteTransferRequest.Status = "recycled"

	if err := c.WasteTransferRequestRepository.Update(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to update waste transfer request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferRequestToSimpleResponse(wasteTransferRequest), nil
}

func (c *WasteTransferRequestUsecase) Search(ctx context.Context, request *model.SearchWasteTransferRequest) ([]model.WasteTransferRequestSimpleResponse, int64, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.WithError(err).Warn("Invalid request body")
		return nil, 0, fiber.ErrBadRequest
	}

	// Set default pagination values if not provided
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Size <= 0 {
		request.Size = 10
	}

	// The repository Search method already handles distance calculation and sorting
	// if Latitude and Longitude are provided in the request
	wasteTransferRequests, total, err := c.WasteTransferRequestRepository.Search(tx, request)
	if err != nil {
		c.Log.WithError(err).Warn("Failed to search waste transfer requests")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.WithError(err).Error("Failed to commit transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.WasteTransferRequestSimpleResponse, len(wasteTransferRequests))
	for i, transferRequest := range wasteTransferRequests {
		responses[i] = *converter.WasteTransferRequestToSimpleResponse(&transferRequest)
	}

	return responses, total, nil
}

func (c *WasteTransferRequestUsecase) Delete(ctx context.Context, request *model.DeleteWasteTransferRequest) (*model.WasteTransferRequestSimpleResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	wasteTransferRequest := new(entity.WasteTransferRequest)
	if err := c.WasteTransferRequestRepository.FindByID(tx, wasteTransferRequest, request.ID); err != nil {
		c.Log.Warnf("Failed to find waste transfer request by ID: %+v", err)
		return nil, fiber.ErrNotFound
	}

	// Parse UUID for finding items
	transferFormUUID, err := uuid.Parse(request.ID)
	if err != nil {
		c.Log.Warnf("Invalid transfer form ID: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Delete related items first
	items, err := c.WasteTransferItemOfferingRepository.FindByTransferFormID(tx, transferFormUUID)
	if err != nil {
		c.Log.Warnf("Failed to find waste transfer items: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	for _, item := range items {
		if err := c.WasteTransferItemOfferingRepository.Delete(tx, &item); err != nil {
			c.Log.Warnf("Failed to delete waste transfer item: %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := c.WasteTransferRequestRepository.Delete(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to delete waste transfer request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferRequestToSimpleResponse(wasteTransferRequest), nil
}
