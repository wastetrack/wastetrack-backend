package repository

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
	"gorm.io/gorm"
)

type WasteDropRequestRepository struct {
	Repository[entity.WasteDropRequest]
	Log *logrus.Logger
}

func NewWasteDropRequestRepository(log *logrus.Logger) *WasteDropRequestRepository {
	return &WasteDropRequestRepository{
		Log: log,
	}
}

func (r *WasteDropRequestRepository) FindByID(db *gorm.DB, wasteDropRequest *entity.WasteDropRequest, id string) error {
	return db.Where("id = ?", id).Preload("AssignedCollector").Preload("Customer").Preload("WasteBank").First(wasteDropRequest).Error
}

func (r *WasteDropRequestRepository) Search(db *gorm.DB, request *model.SearchWasteDropRequest) ([]entity.WasteDropRequest, int64, error) {
	var wasteDropRequests []entity.WasteDropRequest

	// Build the query with distance calculation if coordinates provided
	query := db.Scopes(r.FilterWasteDropRequest(request))

	// Set order direction for created_at (default: DESC)
	orderDir := "DESC"
	if request.OrderDir == "asc" {
		orderDir = "ASC"
	}

	// If latitude and longitude are provided, calculate distance and order by it FIRST
	if request.Latitude != nil && request.Longitude != nil {
		distanceSelect := fmt.Sprintf(`*, 
    CASE 
        WHEN appointment_location IS NOT NULL THEN 
            ST_Distance(
                appointment_location, 
                ST_SetSRID(ST_MakePoint(%f, %f), 4326)
            )
        ELSE NULL 
    END as distance`,
			*request.Longitude, *request.Latitude)

		// Primary: Distance (nearest first), Secondary: created_at (based on order_dir)
		query = query.Select(distanceSelect).
			Order(fmt.Sprintf("distance ASC NULLS LAST, created_at %s", orderDir))
	} else {
		// No coordinates provided, only order by created_at
		query = query.Order(fmt.Sprintf("created_at %s", orderDir))
	}

	// Apply pagination and execute query
	if err := query.Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&wasteDropRequests).Error; err != nil {
		return nil, 0, err
	}

	// Count total records with same filters (without distance calculation for performance)
	var total int64 = 0
	if err := db.Model(&entity.WasteDropRequest{}).
		Scopes(r.FilterWasteDropRequest(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return wasteDropRequests, total, nil
}

func (r *WasteDropRequestRepository) FilterWasteDropRequest(request *model.SearchWasteDropRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if deliveryType := request.DeliveryType; deliveryType != "" {
			tx = tx.Where("delivery_type = ?", deliveryType)
		}
		if customerID := request.CustomerID; customerID != "" {
			tx = tx.Where("customer_id = ?", customerID)
		}
		if wasteBankID := request.WasteBankID; wasteBankID != "" {
			tx = tx.Where("waste_bank_id = ?", wasteBankID)
		}
		if assignedCollectorID := request.AssignedCollectorID; assignedCollectorID != "" {
			tx = tx.Where("assigned_collector_id = ?", assignedCollectorID)
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
		if isDeleted := request.IsDeleted; isDeleted != nil {
			tx = tx.Where("is_deleted = ?", *isDeleted)
		}
		return tx
	}
}

// FindByIDWithDistance finds a waste drop request by ID and calculates distance if coordinates are provided
func (r *WasteDropRequestRepository) FindByIDWithDistance(db *gorm.DB, wasteDropRequest *entity.WasteDropRequest, id string, lat, lng *float64) error {
	query := db.Where("id = ?", id).Preload("AssignedCollector").Preload("Customer").Preload("WasteBank")

	if lat != nil && lng != nil {
		distanceSelect := fmt.Sprintf(`*, 
			CASE 
				WHEN appointment_location IS NOT NULL THEN 
					ST_Distance(
						appointment_location, 
						ST_SetSRID(ST_MakePoint(%f, %f), 4326)
					)
				ELSE NULL 
			END as distance`,
			*lng, *lat)
		query = query.Select(distanceSelect)
	}

	return query.First(wasteDropRequest).Error
}

func (r *WasteDropRequestRepository) UpdateStatus(db *gorm.DB, id string, status string) error {
	return db.Model(&entity.WasteDropRequest{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *WasteDropRequestRepository) AssignCollector(db *gorm.DB, id string, collectorID string) error {
	return db.Model(&entity.WasteDropRequest{}).
		Where("id = ?", id).
		Update("assigned_collector_id", collectorID).Error
}
