package repository

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type WasteTransferRequestRepository struct {
	Repository[entity.WasteTransferRequest]
	Log *logrus.Logger
}

func NewWasteTransferRequestRepository(log *logrus.Logger) *WasteTransferRequestRepository {
	return &WasteTransferRequestRepository{
		Log: log,
	}
}

// FindByID retrieves a transfer request by its ID and preloads related entities.
func (r *WasteTransferRequestRepository) FindByID(db *gorm.DB, wasteTransferRequest *entity.WasteTransferRequest, id string) error {
	return db.Where("id = ?", id).
		Preload("SourceUser").
		Preload("DestinationUser").
		Preload("Items").
		Preload("Items.WasteType").
		First(wasteTransferRequest).Error
}

// FindByIDWithDistance finds a waste transfer request by ID and calculates distance if latitude and longitude are provided
func (r *WasteTransferRequestRepository) FindByIDWithDistance(db *gorm.DB, wasteTransferRequest *entity.WasteTransferRequest, id string, lat, lng *float64) error {
	query := db.Where("id = ?", id).
		Preload("SourceUser").
		Preload("DestinationUser").
		Preload("Items").
		Preload("Items.WasteType")
	if lat != nil && lng != nil {
		// Calculate distance in kilometers
		distanceSelect := fmt.Sprintf(`*, 
			CASE 
				WHEN appointment_location IS NOT NULL THEN 
					ST_Distance(
						appointment_location, 
						ST_SetSRID(ST_MakePoint(%f, %f), 4326)
					)
				ELSE NULL 
			END as distance`, *lng, *lat)
		query = query.Select(distanceSelect)
	}
	return query.First(wasteTransferRequest).Error
}

// Search retrieves waste transfer requests based on search filters.
// If Latitude and Longitude are provided in the request, it calculates and returns the distance.
func (r *WasteTransferRequestRepository) Search(db *gorm.DB, request *model.SearchWasteTransferRequest) ([]entity.WasteTransferRequest, int64, error) {
	var wasteTransferRequests []entity.WasteTransferRequest

	// Start with applying filters
	query := db.Scopes(r.FilterWasteTransferRequest(request))

	// If coordinates are provided, select the distance calculation and order by it
	if request.Latitude != nil && request.Longitude != nil {
		distanceSelect := fmt.Sprintf(`*, 
			CASE 
				WHEN appointment_location IS NOT NULL THEN 
					ST_Distance(
						appointment_location, 
						ST_SetSRID(ST_MakePoint(%f, %f), 4326)
					)
				ELSE NULL 
			END as distance`, *request.Longitude, *request.Latitude)
		query = query.Select(distanceSelect).Order("distance ASC NULLS LAST")
	}

	// Apply pagination and execute query
	if err := query.Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&wasteTransferRequests).Error; err != nil {
		return nil, 0, err
	}

	// Count total records with the same filters (without the distance calculation for performance)
	var total int64 = 0
	if err := db.Model(&entity.WasteTransferRequest{}).
		Scopes(r.FilterWasteTransferRequest(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return wasteTransferRequests, total, nil
}

// FilterWasteTransferRequest applies search filters to the GORM query.
func (r *WasteTransferRequestRepository) FilterWasteTransferRequest(request *model.SearchWasteTransferRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if sourceUserID := request.SourceUserID; sourceUserID != "" {
			tx = tx.Where("source_user_id = ?", sourceUserID)
		}
		if destinationUserID := request.DestinationUserID; destinationUserID != "" {
			tx = tx.Where("destination_user_id = ?", destinationUserID)
		}
		if formType := request.FormType; formType != "" {
			tx = tx.Where("form_type = ?", formType)
		}
		if status := request.Status; status != "" {
			tx = tx.Where("status = ?", status)
		}
		if appointmentDate := request.AppointmentDate; appointmentDate != "" {
			tx = tx.Where("appointment_date = ?", appointmentDate)
		}
		if appointmentStartTime := request.AppointmentStartTime; appointmentStartTime != "" {
			tx = tx.Where("appointment_start_time >= ?", appointmentStartTime)
		}
		if appointmentEndTime := request.AppointmentEndTime; appointmentEndTime != "" {
			tx = tx.Where("appointment_end_time <= ?", appointmentEndTime)
		}
		return tx
	}
}

// UpdateStatus updates the status of a waste transfer request.
func (r *WasteTransferRequestRepository) UpdateStatus(db *gorm.DB, id string, status string) error {
	return db.Model(&entity.WasteTransferRequest{}).
		Where("id = ?", id).
		Update("status", status).Error
}
