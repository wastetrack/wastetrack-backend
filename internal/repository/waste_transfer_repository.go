package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
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
		Preload("AssignedCollector"). // NEW: Preload assigned collector
		Preload("Items").
		Preload("Items.WasteType").
		First(wasteTransferRequest).Error
}

// FindByIDWithDistance finds a waste transfer request by ID and calculates distance if latitude and longitude are provided
func (r *WasteTransferRequestRepository) FindByIDWithDistance(db *gorm.DB, wasteTransferRequest *entity.WasteTransferRequest, id string, lat, lng *float64) error {
	query := db.Where("id = ?", id).
		Preload("SourceUser").
		Preload("DestinationUser").
		Preload("AssignedCollector"). // NEW: Preload assigned collector
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
	// Set order direction for created_at (default: DESC)
	orderDir := "DESC"
	if request.OrderDir == "asc" {
		orderDir = "ASC"
	}

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
		query = query.Select(distanceSelect).Order(fmt.Sprintf("distance ASC NULLS LAST, created_at %s", orderDir))
	} else {
		// No coordinates provided, only order by created_at
		query = query.Order(fmt.Sprintf("created_at %s", orderDir))
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
		// NEW: Filter by assigned collector
		if assignedCollectorID := request.AssignedCollectorID; assignedCollectorID != "" {
			tx = tx.Where("assigned_collector_id = ?", assignedCollectorID)
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
		if isDeleted := request.IsDeleted; isDeleted != nil {
			tx = tx.Where("is_deleted = ?", *isDeleted)
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
func (r *WasteTransferItemOfferingRepository) UpdateVerifiedWeight(tx *gorm.DB, item *entity.WasteTransferItemOffering) error {
	return tx.Model(item).Select("verified_weight").Updates(item).Error
}

// NEW: AssignCollector assigns a collector to a waste transfer request and updates status to "assigned"
func (r *WasteTransferRequestRepository) AssignCollector(db *gorm.DB, id string, collectorID uuid.UUID) error {
	return db.Model(&entity.WasteTransferRequest{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"assigned_collector_id": collectorID,
			"status":                "assigned",
		}).Error
}

// NEW: FindByCollectorID finds all waste transfer requests assigned to a specific collector
func (r *WasteTransferRequestRepository) FindByCollectorID(db *gorm.DB, collectorID string) ([]entity.WasteTransferRequest, error) {
	var requests []entity.WasteTransferRequest
	err := db.Where("assigned_collector_id = ?", collectorID).
		Preload("SourceUser").
		Preload("DestinationUser").
		Preload("AssignedCollector").
		Preload("Items").
		Preload("Items.WasteType").
		Find(&requests).Error
	return requests, err
}

func (r *WasteTransferRequestRepository) GetTopOfftakers(db *gorm.DB, request *model.GovernmentDashboardRequest) ([]model.TopOfftaker, error) {
	var topOfftakers []model.TopOfftaker

	// Build the query with JOIN to get user details in a single query
	// Using COALESCE to handle potential NULL values and using correct column names
	query := db.Model(&entity.WasteTransferRequest{}).
		Select(`
			destination_users.id,
			COALESCE(destination_users.username, '') as name,
			COALESCE(destination_users.institution, '') as institution,
			COALESCE(destination_users.city, '') as city,
			COALESCE(destination_users.province, '') as province,
			SUM(waste_transfer_requests.total_weight) as total_weight,
			SUM(waste_transfer_requests.total_price) as total_price
		`).
		Joins("JOIN users AS destination_users ON waste_transfer_requests.destination_user_id = destination_users.id").
		Where("waste_transfer_requests.form_type = ? AND waste_transfer_requests.status = ? AND waste_transfer_requests.is_deleted = ?",
			"industry_request", "completed", false).
		Group(`
			destination_users.id,
			destination_users.username,
			destination_users.institution,
			destination_users.city,
			destination_users.province
		`).
		Order("total_weight DESC, total_price DESC").
		Limit(3)

	// Apply date filters if provided
	if request.StartMonth != "" {
		startDate, err := time.Parse("2006-01", request.StartMonth)
		if err != nil {
			return nil, fmt.Errorf("invalid start_month format: %v", err)
		}
		query = query.Where("waste_transfer_requests.created_at >= ?", startDate)
	}

	if request.EndMonth != "" {
		endDate, err := time.Parse("2006-01", request.EndMonth)
		if err != nil {
			return nil, fmt.Errorf("invalid end_month format: %v", err)
		}
		endDate = endDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("waste_transfer_requests.created_at <= ?", endDate)
	}

	// Apply location filters if provided
	if request.Province != "" {
		query = query.Where("destination_users.province ILIKE ?", "%"+request.Province+"%")
	}
	if request.City != "" {
		query = query.Where("destination_users.city ILIKE ?", "%"+request.City+"%")
	}

	// Execute the query directly into the result struct
	rows, err := query.Rows()
	if err != nil {
		return nil, fmt.Errorf("database error occurred while getting top offtakers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var topOfftaker model.TopOfftaker
		var totalWeight float64
		var id uuid.UUID

		if err := rows.Scan(
			&id,
			&topOfftaker.Name,
			&topOfftaker.Institution,
			&topOfftaker.City,
			&topOfftaker.Province,
			&totalWeight,
			&topOfftaker.TotalPrice,
		); err != nil {
			return nil, fmt.Errorf("error scanning top offtaker data: %v", err)
		}

		topOfftaker.ID = id.String()
		topOfftaker.TotalWeight = totalWeight
		topOfftakers = append(topOfftakers, topOfftaker)
	}

	return topOfftakers, nil
}
func (r *WasteTransferRequestRepository) GetWasteBankTransferTotals(db *gorm.DB, request *model.GovernmentDashboardRequest) (map[string]float64, error) {
	var results []struct {
		UserID      string  `json:"user_id"`
		TotalWeight float64 `json:"total_weight"`
	}

	// Use helper functions for date preparation
	startMonth, startDate := r.prepareDateParams(request.StartMonth)
	endMonth, endDate := r.prepareEndDateParams(request.EndMonth)

	query := db.Raw(`
		SELECT 
			wtr.destination_user_id::text as user_id,
			COALESCE(SUM(wtr.total_weight), 0) as total_weight
		FROM waste_transfer_requests wtr
		JOIN users u ON wtr.destination_user_id = u.id
		WHERE wtr.form_type = 'waste_bank_request'
		  AND wtr.status = 'completed'
		  AND wtr.is_deleted = false
		  AND u.role IN ('waste_bank_unit', 'waste_bank_central')
		  AND (? = '' OR wtr.created_at >= ?::timestamp)
		  AND (? = '' OR wtr.created_at <= ?::timestamp)
		  AND (? = '' OR u.province ILIKE ?)
		  AND (? = '' OR u.city ILIKE ?)
		GROUP BY wtr.destination_user_id
	`,
		startMonth, startDate,
		endMonth, endDate,
		request.Province, "%"+request.Province+"%",
		request.City, "%"+request.City+"%")

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get waste transfer totals: %v", err)
	}

	// Convert to map
	totals := make(map[string]float64)
	for _, result := range results {
		totals[result.UserID] = result.TotalWeight
	}

	return totals, nil
}

func (r *WasteTransferRequestRepository) prepareDateParams(monthStr string) (string, string) {
	if monthStr == "" {
		return "", ""
	}

	if parsed, err := time.Parse("2006-01", monthStr); err == nil {
		return monthStr, parsed.Format("2006-01-02 15:04:05")
	}

	return "", ""
}
func (r *WasteTransferRequestRepository) prepareEndDateParams(monthStr string) (string, string) {
	if monthStr == "" {
		return "", ""
	}

	if parsed, err := time.Parse("2006-01", monthStr); err == nil {
		// Set to last moment of the month: last day + 23:59:59
		endOfMonth := parsed.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		return monthStr, endOfMonth.Format("2006-01-02 15:04:05")
	}

	return "", ""
}
