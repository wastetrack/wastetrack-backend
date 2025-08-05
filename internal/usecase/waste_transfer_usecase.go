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
	"github.com/wastetrack/wastetrack-backend/pkg/timezone"

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
	// Storage repositories
	StorageRepository     *repository.StorageRepository
	StorageItemRepository *repository.StorageItemRepository
	// NEW: Profile repositories
	IndustryRepository          *repository.IndustryRepository
	WasteBankRepository         *repository.WasteBankRepository
	SalaryTransactionRepository *repository.SalaryTransactionRepository
}

func NewWasteTransferRequestUsecase(
	db *gorm.DB,
	log *logrus.Logger,
	validate *validator.Validate,
	wasteTransferRequestRepository *repository.WasteTransferRequestRepository,
	wasteTransferItemOfferingRepository *repository.WasteTransferItemOfferingRepository,
	userRepository *repository.UserRepository,
	wasteTypeRepository *repository.WasteTypeRepository,
	storageRepository *repository.StorageRepository,
	storageItemRepository *repository.StorageItemRepository,
	industryRepository *repository.IndustryRepository,
	wasteBankRepository *repository.WasteBankRepository,
	salaryTransactionRepository *repository.SalaryTransactionRepository,
) *WasteTransferRequestUsecase {
	return &WasteTransferRequestUsecase{
		DB:                                  db,
		Log:                                 log,
		Validate:                            validate,
		WasteTransferRequestRepository:      wasteTransferRequestRepository,
		WasteTransferItemOfferingRepository: wasteTransferItemOfferingRepository,
		UserRepository:                      userRepository,
		WasteTypeRepository:                 wasteTypeRepository,
		StorageRepository:                   storageRepository,
		StorageItemRepository:               storageItemRepository,
		IndustryRepository:                  industryRepository,
		WasteBankRepository:                 wasteBankRepository,
		SalaryTransactionRepository:         salaryTransactionRepository,
	}
}

// NEW: Helper method to find or create storage for raw materials (not recycled)
func (c *WasteTransferRequestUsecase) findOrCreateRawMaterialStorage(tx *gorm.DB, userID uuid.UUID) (*entity.Storage, error) {
	c.Log.Infof("Finding or creating raw material storage for user ID: %s", userID.String())

	// Try to find existing storage for raw materials (not recycled)
	searchReq := &model.SearchStorageRequest{
		UserID:                userID.String(),
		IsForRecycledMaterial: &[]bool{false}[0], // Pointer to false
		Page:                  1,
		Size:                  1,
	}

	storages, _, err := c.StorageRepository.Search(tx, searchReq)
	if err != nil {
		c.Log.Warnf("Failed to search storage: %+v", err)
		return nil, err
	}

	// If storage exists, return the first one
	if len(storages) > 0 {
		c.Log.Infof("Found existing raw material storage ID: %s", storages[0].ID.String())
		return &storages[0], nil
	}

	// Create new storage if none exists
	c.Log.Infof("Creating new raw material storage for user")
	storage := &entity.Storage{
		UserID:                userID,
		Length:                10.0, // Default dimensions - you might want to make these configurable
		Width:                 10.0,
		Height:                3.0,
		IsForRecycledMaterial: false, // Raw materials storage
	}

	if err := c.StorageRepository.Create(tx, storage); err != nil {
		c.Log.Warnf("Failed to create storage: %+v", err)
		return nil, err
	}

	c.Log.Infof("Successfully created new raw material storage ID: %s", storage.ID.String())
	return storage, nil
}

// NEW: Helper method to subtract items from source storage
func (c *WasteTransferRequestUsecase) subtractFromSourceStorage(tx *gorm.DB, storageID uuid.UUID, items []entity.WasteTransferItemOffering) error {
	c.Log.Infof("Subtracting %d items from source storage ID: %s", len(items), storageID.String())

	for _, item := range items {
		if item.AcceptedWeight <= 0 {
			c.Log.Warnf("Skipping item with zero or negative accepted weight: %f", item.AcceptedWeight)
			continue
		}

		// Check if storage item exists for this waste type
		var existingStorageItem entity.StorageItem
		err := tx.Where("storage_id = ? AND waste_type_id = ?", storageID, item.WasteTypeID).
			First(&existingStorageItem).Error

		if err == nil {
			// Storage item exists, subtract from existing weight
			c.Log.Infof("Subtracting from existing storage item for waste type %s: removing %f kg from existing %f kg",
				item.WasteTypeID.String(), item.AcceptedWeight, existingStorageItem.WeightKgs)

			if existingStorageItem.WeightKgs < item.AcceptedWeight {
				return fmt.Errorf("insufficient stock in storage for waste type %s: available %f kg, requested %f kg",
					item.WasteTypeID.String(), existingStorageItem.WeightKgs, item.AcceptedWeight)
			}

			existingStorageItem.WeightKgs -= item.AcceptedWeight
			existingStorageItem.UpdatedAt = time.Now()

			// If weight becomes zero or negative, delete the storage item
			if existingStorageItem.WeightKgs <= 0 {
				c.Log.Infof("Deleting storage item for waste type %s as weight is now %f kg",
					item.WasteTypeID.String(), existingStorageItem.WeightKgs)

				if err := c.StorageItemRepository.Delete(tx, &existingStorageItem); err != nil {
					c.Log.Warnf("Failed to delete storage item: %+v", err)
					return err
				}
			} else {
				if err := c.StorageItemRepository.Update(tx, &existingStorageItem); err != nil {
					c.Log.Warnf("Failed to update existing storage item: %+v", err)
					return err
				}
			}
		} else if err == gorm.ErrRecordNotFound {
			// Storage item doesn't exist - this is an error for subtraction
			return fmt.Errorf("cannot subtract waste type %s: not found in source storage", item.WasteTypeID.String())
		} else {
			// Database error
			c.Log.Warnf("Database error while checking storage item: %+v", err)
			return err
		}
	}

	c.Log.Infof("Successfully subtracted all items from source storage")
	return nil
}

// NEW: Helper method to add items to destination storage
func (c *WasteTransferRequestUsecase) addToDestinationStorage(tx *gorm.DB, storageID uuid.UUID, items []entity.WasteTransferItemOffering) error {
	c.Log.Infof("Adding %d items to destination storage ID: %s", len(items), storageID.String())

	for _, item := range items {
		if item.VerifiedWeight <= 0 {
			c.Log.Warnf("Skipping item with zero or negative verified weight: %f", item.VerifiedWeight)
			continue
		}

		// Check if storage item already exists for this waste type
		var existingStorageItem entity.StorageItem
		err := tx.Where("storage_id = ? AND waste_type_id = ?", storageID, item.WasteTypeID).
			First(&existingStorageItem).Error

		if err == nil {
			// Storage item exists, add to existing weight
			c.Log.Infof("Updating existing storage item for waste type %s: adding %f kg to existing %f kg",
				item.WasteTypeID.String(), item.VerifiedWeight, existingStorageItem.WeightKgs)

			existingStorageItem.WeightKgs += item.VerifiedWeight
			existingStorageItem.UpdatedAt = time.Now()

			if err := c.StorageItemRepository.Update(tx, &existingStorageItem); err != nil {
				c.Log.Warnf("Failed to update existing storage item: %+v", err)
				return err
			}
		} else if err == gorm.ErrRecordNotFound {
			// Storage item doesn't exist, create new one
			c.Log.Infof("Creating new storage item for waste type %s with weight %f kg",
				item.WasteTypeID.String(), item.VerifiedWeight)

			newStorageItem := &entity.StorageItem{
				StorageID:   storageID,
				WasteTypeID: item.WasteTypeID,
				WeightKgs:   item.VerifiedWeight,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}

			if err := c.StorageItemRepository.Create(tx, newStorageItem); err != nil {
				c.Log.Warnf("Failed to create new storage item: %+v", err)
				return err
			}
		} else {
			// Database error
			c.Log.Warnf("Database error while checking storage item: %+v", err)
			return err
		}
	}

	c.Log.Infof("Successfully processed all items for destination storage")
	return nil
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

	// Parse appointment date
	appointmentDate, err := time.Parse("2006-01-02", request.AppointmentDate)
	if err != nil {
		c.Log.Warnf("Invalid appointment date format: %+v", err)
		return nil, fiber.ErrBadRequest
	}

	var appointmentStartTime, appointmentEndTime types.TimeOnly
	var startTimeLocation, endTimeLocation *time.Location

	// Parse and validate appointment start time
	if request.AppointmentStartTime != "" {
		startTime, location, err := timezone.ParseTimeWithTimezone(request.AppointmentStartTime)
		if err != nil {
			c.Log.Warnf("Invalid appointment start time format: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		startTimeLocation = location
		appointmentStartTime = types.NewTimeOnly(startTime)

		// Check if appointment start time is in the past
		if timezone.IsDateTimeInPast(appointmentDate, startTime, location) {
			c.Log.Warnf("Appointment start time cannot be in the past")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Appointment start time cannot be in the past")
		}
	}

	// Parse and validate appointment end time
	if request.AppointmentEndTime != "" {
		endTime, location, err := timezone.ParseTimeWithTimezone(request.AppointmentEndTime)
		if err != nil {
			c.Log.Warnf("Invalid appointment end time format: %+v", err)
			return nil, fiber.ErrBadRequest
		}
		endTimeLocation = location
		appointmentEndTime = types.NewTimeOnly(endTime)

		// Check if appointment end time is in the past
		if timezone.IsDateTimeInPast(appointmentDate, endTime, location) {
			c.Log.Warnf("Appointment end time cannot be in the past")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Appointment end time cannot be in the past")
		}
	}

	// Additional validation: if both start and end times are provided, ensure end time is after start time
	if request.AppointmentStartTime != "" && request.AppointmentEndTime != "" {
		startTimeWithDate := time.Date(
			appointmentDate.Year(),
			appointmentDate.Month(),
			appointmentDate.Day(),
			appointmentStartTime.Time.Hour(),
			appointmentStartTime.Time.Minute(),
			appointmentStartTime.Time.Second(),
			0,
			startTimeLocation,
		)

		endTimeWithDate := time.Date(
			appointmentDate.Year(),
			appointmentDate.Month(),
			appointmentDate.Day(),
			appointmentEndTime.Time.Hour(),
			appointmentEndTime.Time.Minute(),
			appointmentEndTime.Time.Second(),
			0,
			endTimeLocation,
		)

		if endTimeWithDate.Before(startTimeWithDate) || endTimeWithDate.Equal(startTimeWithDate) {
			c.Log.Warnf("Appointment end time must be after start time")
			return nil, fiber.NewError(fiber.StatusBadRequest, "Appointment end time must be after start time")
		}
	}

	// Validate appointment date is not in the past (even without time)
	today := time.Now().In(timezone.WIB).Truncate(24 * time.Hour)
	appointmentDateOnly := appointmentDate.Truncate(24 * time.Hour)
	if appointmentDateOnly.Before(today) {
		c.Log.Warnf("Appointment date cannot be in the past")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Appointment date cannot be in the past")
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
	var totalOfferingPrice int64
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
		totalOfferingPrice += int64(weight) * pricePerKg
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
	var totalAcceptedPrice int64

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
			totalAcceptedPrice += int64(math.Round(pricing.AcceptedWeight)) * pricing.AcceptedPricePerKgs
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
	wasteTransferRequest.TotalPrice = totalAcceptedPrice

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

	// Validate status - can only complete assigned or collecting requests
	if wasteTransferRequest.Status != "assigned" && wasteTransferRequest.Status != "collecting" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Can only complete assigned or collecting requests")
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
	var totalVerifiedPrice int64
	updatedItemsCount := 0

	for i := range currentItems {
		if verifiedWeight, exists := wasteTypeWeights[currentItems[i].WasteTypeID]; exists {
			// Validate verified weight doesn't exceed accepted weight (if accepted weight is set)
			if currentItems[i].AcceptedWeight > 0 && verifiedWeight > currentItems[i].AcceptedWeight {
				return nil, fiber.NewError(fiber.StatusBadRequest,
					fmt.Sprintf("Verified weight (%.2f) cannot exceed accepted weight (%.2f) for waste type: %s",
						verifiedWeight, currentItems[i].AcceptedWeight, currentItems[i].WasteTypeID))
			}

			// Update the verified weight
			currentItems[i].VerifiedWeight = verifiedWeight

			if err := c.WasteTransferItemOfferingRepository.UpdateVerifiedWeight(tx, &currentItems[i]); err != nil {
				c.Log.Warnf("Failed to update waste transfer item verified weight: %+v", err)
				return nil, fiber.ErrInternalServerError
			}

			totalVerifiedWeight += verifiedWeight
			// Use accepted price per kg for calculation, fallback to offering price if not set
			pricePerKg := currentItems[i].AcceptedPricePerKgs
			if pricePerKg == 0 {
				pricePerKg = currentItems[i].OfferingPricePerKgs
			}
			totalVerifiedPrice += int64(float64(pricePerKg) * verifiedWeight)
			updatedItemsCount++
		}
	}

	// Ensure all provided waste types were found and updated
	if updatedItemsCount != len(request.Items.WasteTypeIDs) {
		return nil, fiber.NewError(fiber.StatusBadRequest,
			"Some waste types not found in this transfer request")
	}

	// Create salary transaction record for the payment
	// salaryTransaction := &entity.SalaryTransaction{
	// 	SenderID:        buyer.ID,  // Buyer is the sender (payer)
	// 	ReceiverID:      seller.ID, // Seller is the receiver (payee)
	// 	TransactionType: "waste_payment",
	// 	Amount:          totalPaymentAmount,
	// 	Status:          "completed",
	// 	Notes:           fmt.Sprintf("Payment for waste transfer request: %s", wasteTransferRequest.ID.String()),
	// }

	// // Create the salary transaction record
	// if err := c.SalaryTransactionRepository.Create(tx, salaryTransaction); err != nil {
	// 	c.Log.Warnf("Failed to create salary transaction: %+v", err)
	// 	return nil, fiber.ErrInternalServerError
	// }

	// c.Log.Infof("Successfully created waste payment transaction: %s", salaryTransaction.ID.String())

	// Handle storage operations
	c.Log.Infof("Starting storage operations for waste transfer completion")

	// Find or create source storage (raw materials)
	sourceStorage, err := c.findOrCreateRawMaterialStorage(tx, wasteTransferRequest.SourceUserID)
	if err != nil {
		c.Log.Warnf("Failed to find or create source storage: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Find or create destination storage (raw materials)
	destinationStorage, err := c.findOrCreateRawMaterialStorage(tx, wasteTransferRequest.DestinationUserID)
	if err != nil {
		c.Log.Warnf("Failed to find or create destination storage: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Subtract from source storage using accepted weights
	if err := c.subtractFromSourceStorage(tx, sourceStorage.ID, currentItems); err != nil {
		c.Log.Warnf("Failed to subtract from source storage: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Storage operation failed: %v", err))
	}

	// Add to destination storage using verified weights
	if err := c.addToDestinationStorage(tx, destinationStorage.ID, currentItems); err != nil {
		c.Log.Warnf("Failed to add to destination storage: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("Successfully completed storage operations: subtracted from source storage %s, added to destination storage %s",
		sourceStorage.ID.String(), destinationStorage.ID.String())

	// Update destination user profile
	if err := c.updateDestinationUserProfile(tx, wasteTransferRequest.DestinationUserID, totalVerifiedWeight); err != nil {
		c.Log.Warnf("Failed to update destination user profile: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Update the waste transfer request
	wasteTransferRequest.Status = "completed"
	wasteTransferRequest.TotalWeight = totalVerifiedWeight
	wasteTransferRequest.TotalPrice = int64(totalVerifiedPrice)

	if err := c.WasteTransferRequestRepository.Update(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to update waste transfer request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	c.Log.Infof("Successfully completed waste transfer request with payment and storage integration")
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
		if err := c.WasteTransferItemOfferingRepository.SoftDelete(tx, &item); err != nil {
			c.Log.Warnf("Failed to delete waste transfer item: %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	if err := c.WasteTransferRequestRepository.SoftDelete(tx, wasteTransferRequest); err != nil {
		c.Log.Warnf("Failed to delete waste transfer request: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.WasteTransferRequestToSimpleResponse(wasteTransferRequest), nil
}

func (c *WasteTransferRequestUsecase) updateDestinationUserProfile(tx *gorm.DB, destinationUserID uuid.UUID, totalWeight float64) error {
	c.Log.Infof("Updating destination user profile for ID: %s with weight: %f", destinationUserID.String(), totalWeight)

	// Get the destination user to check their role
	user := &entity.User{}
	if err := c.UserRepository.FindById(tx, user, destinationUserID.String()); err != nil {
		c.Log.Warnf("Failed to find destination user: %+v", err)
		return err
	}

	// Update profile based on user role
	switch user.Role {
	case "industry":
		return c.updateIndustryProfile(tx, destinationUserID, totalWeight, 0) // 0 for recycled weight in completion
	case "waste_bank":
		return c.updateWasteBankProfileTransfer(tx, destinationUserID, totalWeight)
	default:
		c.Log.Infof("User role %s does not require profile weight update", user.Role)
		return nil // No error, just no update needed for other roles
	}
}

// NEW: Helper method to update industry profile
func (c *WasteTransferRequestUsecase) updateIndustryProfile(tx *gorm.DB, userID uuid.UUID, wasteWeight float64, recycledWeight float64) error {
	c.Log.Infof("Updating industry profile for user ID: %s with waste weight: %f, recycled weight: %f",
		userID.String(), wasteWeight, recycledWeight)

	// Find or create industry profile
	industryProfile := &entity.IndustryProfile{}
	err := c.IndustryRepository.FindByUserID(tx, industryProfile, userID.String())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Infof("Creating new industry profile for user ID: %s", userID.String())
			// Create new profile if doesn't exist
			industryProfile = &entity.IndustryProfile{
				UserID:              userID,
				TotalWasteWeight:    0,
				TotalRecycledWeight: 0,
			}
			if err := c.IndustryRepository.Create(tx, industryProfile); err != nil {
				c.Log.Warnf("Failed to create industry profile: %+v", err)
				return err
			}
		} else {
			c.Log.Warnf("Failed to find industry profile: %+v", err)
			return err
		}
	}

	// Update weights
	oldWasteWeight := industryProfile.TotalWasteWeight
	oldRecycledWeight := industryProfile.TotalRecycledWeight
	industryProfile.TotalWasteWeight += wasteWeight
	industryProfile.TotalRecycledWeight += recycledWeight

	c.Log.Infof("Updating industry profile: waste weight from %f to %f, recycled weight from %f to %f",
		oldWasteWeight, industryProfile.TotalWasteWeight, oldRecycledWeight, industryProfile.TotalRecycledWeight)

	if err := c.IndustryRepository.Update(tx, industryProfile); err != nil {
		c.Log.Warnf("Failed to update industry profile: %+v", err)
		return err
	}

	c.Log.Infof("Successfully updated industry profile")
	return nil
}

// NEW: Helper method to update waste bank profile (reusing existing logic but improved)
func (c *WasteTransferRequestUsecase) updateWasteBankProfileTransfer(tx *gorm.DB, userID uuid.UUID, totalWeight float64) error {
	c.Log.Infof("Updating waste bank profile for user ID: %s with weight: %f", userID.String(), totalWeight)

	// Find or create waste bank profile
	wasteBankProfile := &entity.WasteBankProfile{}
	err := c.WasteBankRepository.FindByUserID(tx, wasteBankProfile, userID.String())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Log.Infof("Creating new waste bank profile for user ID: %s", userID.String())
			// Create new profile if doesn't exist
			wasteBankProfile = &entity.WasteBankProfile{
				UserID:           userID,
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
