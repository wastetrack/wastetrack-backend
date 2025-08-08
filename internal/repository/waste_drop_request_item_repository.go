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

type WasteDropRequestItemRepository struct {
	Repository[entity.WasteDropRequestItem]
	Log *logrus.Logger
}

func NewWasteDropRequestItemRepository(log *logrus.Logger) *WasteDropRequestItemRepository {
	return &WasteDropRequestItemRepository{
		Log: log,
	}
}

func (r *WasteDropRequestItemRepository) FindByID(db *gorm.DB, wasteDropRequestItem *entity.WasteDropRequestItem, id string) error {
	return db.Where("id = ?", id).Preload("Request").Preload("WasteType").First(wasteDropRequestItem).Error
}

func (r *WasteDropRequestItemRepository) Search(db *gorm.DB, request *model.SearchWasteDropRequestItemRequest) ([]entity.WasteDropRequestItem, int64, error) {
	var wasteDropRequestItems []entity.WasteDropRequestItem

	// Apply filters and pagination (no preloads for simple response)
	if err := db.Scopes(r.FilterWasteDropRequestItem(request)).
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&wasteDropRequestItems).Error; err != nil {
		return nil, 0, err
	}

	// Count total records with same filters
	var total int64 = 0
	if err := db.Model(&entity.WasteDropRequestItem{}).
		Scopes(r.FilterWasteDropRequestItem(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return wasteDropRequestItems, total, nil
}

func (r *WasteDropRequestItemRepository) FilterWasteDropRequestItem(request *model.SearchWasteDropRequestItemRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if requestID := request.RequestID; requestID != "" {
			tx = tx.Where("request_id = ?", requestID)
		}
		if wasteTypeID := request.WasteTypeID; wasteTypeID != "" {
			tx = tx.Where("waste_type_id = ?", wasteTypeID)
		}
		return tx
	}
}

func (r *WasteDropRequestItemRepository) FindByDropFormID(db *gorm.DB, dropFormID uuid.UUID) ([]entity.WasteDropRequestItem, error) {
	var items []entity.WasteDropRequestItem
	err := db.Where("request_id = ?", dropFormID).
		Preload("WasteType").
		Find(&items).Error
	return items, err
}

func (r *WasteDropRequestItemRepository) SoftDelete(db *gorm.DB, wasteDropRequestItem *entity.WasteDropRequestItem) error {
	return db.Delete(wasteDropRequestItem).Error
}

// GetTotalCollectedWaste calculates the total verified weight for government dashboard
func (r *WasteDropRequestItemRepository) GetTotalCollectedWaste(db *gorm.DB, request *model.GovernmentDashboardRequest) (float64, error) {
	var totalWeight float64

	query := db.Model(&entity.WasteDropRequestItem{}).
		Select("COALESCE(SUM(verified_weight), 0)").
		Where("waste_drop_request_items.is_deleted = ?", false)

	// Apply date filters if provided
	if request.StartMonth != "" || request.EndMonth != "" {
		// Join with waste_drop_requests to filter by creation date
		query = query.Joins("JOIN waste_drop_requests ON waste_drop_request_items.request_id = waste_drop_requests.id")

		if request.StartMonth != "" {
			startDate, err := time.Parse("2006-01", request.StartMonth)
			if err != nil {
				return 0, fmt.Errorf("invalid start_month format: %v", err)
			}
			query = query.Where("waste_drop_requests.created_at >= ?", startDate)
		}

		if request.EndMonth != "" {
			endDate, err := time.Parse("2006-01", request.EndMonth)
			if err != nil {
				return 0, fmt.Errorf("invalid end_month format: %v", err)
			}
			// Set to last moment of the month
			endDate = endDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			query = query.Where("waste_drop_requests.created_at <= ?", endDate)
		}
	}

	// Apply location filters if provided
	if request.Province != "" || request.City != "" {
		// Join with users through waste_drop_requests to filter by location
		query = query.Joins("JOIN users ON waste_drop_requests.customer_id = users.id")

		if request.Province != "" {
			query = query.Where("users.province ILIKE ?", "%"+request.Province+"%")
		}

		if request.City != "" {
			query = query.Where("users.city ILIKE ?", "%"+request.City+"%")
		}
	}

	// Execute the query
	if err := query.Scan(&totalWeight).Error; err != nil {
		return 0, fmt.Errorf("database error occurred while calculating total collected waste")
	}

	return totalWeight, nil
}

// GetCollectionTrendsByRole returns monthly collection trends grouped by user role
func (r *WasteDropRequestItemRepository) GetCollectionTrendsByRole(db *gorm.DB, request *model.GovernmentDashboardRequest) (map[string]map[string]float64, error) {
	// Structure: map[month]map[role]weight
	trends := make(map[string]map[string]float64)

	// Query for waste drop items - group by month and waste bank role
	wasteDropQuery := db.Model(&entity.WasteDropRequestItem{}).
		Select("DATE_TRUNC('month', waste_drop_requests.created_at) as month, waste_banks.role, COALESCE(SUM(waste_drop_request_items.verified_weight), 0) as total_weight").
		Joins("JOIN waste_drop_requests ON waste_drop_request_items.request_id = waste_drop_requests.id").
		Joins("JOIN users AS waste_banks ON waste_drop_requests.waste_bank_id = waste_banks.id").
		Where("waste_drop_request_items.is_deleted = ?", false).
		Group("DATE_TRUNC('month', waste_drop_requests.created_at), waste_banks.role").
		Order("month")

	// Apply date filters for waste drop
	if request.StartMonth != "" {
		startDate, err := time.Parse("2006-01", request.StartMonth)
		if err != nil {
			return nil, fmt.Errorf("invalid start_month format: %v", err)
		}
		wasteDropQuery = wasteDropQuery.Where("waste_drop_requests.created_at >= ?", startDate)
	}

	if request.EndMonth != "" {
		endDate, err := time.Parse("2006-01", request.EndMonth)
		if err != nil {
			return nil, fmt.Errorf("invalid end_month format: %v", err)
		}
		endDate = endDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		wasteDropQuery = wasteDropQuery.Where("waste_drop_requests.created_at <= ?", endDate)
	}

	// Apply location filters for waste drop if provided
	if request.Province != "" || request.City != "" {
		if request.Province != "" {
			wasteDropQuery = wasteDropQuery.Where("waste_banks.province ILIKE ?", "%"+request.Province+"%")
		}
		if request.City != "" {
			wasteDropQuery = wasteDropQuery.Where("waste_banks.city ILIKE ?", "%"+request.City+"%")
		}
	}

	// Execute waste drop query
	wasteDropRows, err := wasteDropQuery.Rows()
	if err != nil {
		return nil, fmt.Errorf("database error occurred while getting waste drop collection trends")
	}
	defer wasteDropRows.Close()

	// Process waste drop results
	for wasteDropRows.Next() {
		var month time.Time
		var role string
		var totalWeight float64

		if err := wasteDropRows.Scan(&month, &role, &totalWeight); err != nil {
			return nil, fmt.Errorf("error scanning waste drop trend data")
		}

		monthStr := month.Format("2006-01")
		if trends[monthStr] == nil {
			trends[monthStr] = make(map[string]float64)
		}
		trends[monthStr][role] += totalWeight
	}

	return trends, nil
}

// GetWasteDropVerifiedWeightByMonth returns monthly total verified weight from waste drop requests
func (r *WasteDropRequestItemRepository) GetWasteDropVerifiedWeightByMonth(db *gorm.DB, request *model.GovernmentDashboardRequest) (map[string]float64, error) {
	weightByMonth := make(map[string]float64)

	query := db.Model(&entity.WasteDropRequestItem{}).
		Select("DATE_TRUNC('month', waste_drop_requests.created_at) as month, COALESCE(SUM(waste_drop_request_items.verified_weight), 0) as total_weight").
		Joins("JOIN waste_drop_requests ON waste_drop_request_items.request_id = waste_drop_requests.id").
		Where("waste_drop_request_items.is_deleted = ?", false).
		Group("DATE_TRUNC('month', waste_drop_requests.created_at)").
		Order("month")

	// Apply date filters
	if request.StartMonth != "" {
		startDate, err := time.Parse("2006-01", request.StartMonth)
		if err != nil {
			return nil, fmt.Errorf("invalid start_month format: %v", err)
		}
		query = query.Where("waste_drop_requests.created_at >= ?", startDate)
	}

	if request.EndMonth != "" {
		endDate, err := time.Parse("2006-01", request.EndMonth)
		if err != nil {
			return nil, fmt.Errorf("invalid end_month format: %v", err)
		}
		endDate = endDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		query = query.Where("waste_drop_requests.created_at <= ?", endDate)
	}

	// Apply location filters if provided (filter by waste bank location)
	if request.Province != "" || request.City != "" {
		query = query.Joins("JOIN users AS waste_banks ON waste_drop_requests.waste_bank_id = waste_banks.id")

		if request.Province != "" {
			query = query.Where("waste_banks.province ILIKE ?", "%"+request.Province+"%")
		}
		if request.City != "" {
			query = query.Where("waste_banks.city ILIKE ?", "%"+request.City+"%")
		}
	}

	// Execute query
	rows, err := query.Rows()
	if err != nil {
		return nil, fmt.Errorf("database error occurred while getting waste drop verified weights")
	}
	defer rows.Close()

	for rows.Next() {
		var month time.Time
		var totalWeight float64

		if err := rows.Scan(&month, &totalWeight); err != nil {
			return nil, fmt.Errorf("error scanning waste drop verified weight data")
		}

		monthStr := month.Format("2006-01")
		weightByMonth[monthStr] = totalWeight
	}

	return weightByMonth, nil
}

func (r *WasteDropRequestItemRepository) GetWasteBankWasteDropTotals(db *gorm.DB, request *model.GovernmentDashboardRequest) (map[string]float64, error) {
	var results []struct {
		UserID      string  `json:"user_id"`
		TotalWeight float64 `json:"total_weight"`
	}

	// Use helper functions for date preparation
	startMonth, startDate := r.prepareDateParams(request.StartMonth)
	endMonth, endDate := r.prepareDateParams(request.EndMonth)

	query := db.Raw(`
		SELECT 
			wdr.waste_bank_id::text as user_id,
			COALESCE(SUM(wdri.verified_weight), 0) as total_weight
		FROM waste_drop_requests wdr
		JOIN waste_drop_request_items wdri ON wdr.id = wdri.request_id
		JOIN users u ON wdr.waste_bank_id = u.id
		WHERE wdr.is_deleted = false 
		  AND wdri.is_deleted = false
		  AND u.role IN ('waste_bank_unit', 'waste_bank_central')
		  AND (? = '' OR wdr.created_at >= ?::timestamp)
		  AND (? = '' OR wdr.created_at <= ?::timestamp)
		  AND (? = '' OR u.province ILIKE ?)
		  AND (? = '' OR u.city ILIKE ?)
		GROUP BY wdr.waste_bank_id
	`,
		startMonth, startDate,
		endMonth, endDate,
		request.Province, "%"+request.Province+"%",
		request.City, "%"+request.City+"%")

	if err := query.Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get waste drop totals: %v", err)
	}

	// Convert to map
	totals := make(map[string]float64)
	for _, result := range results {
		totals[result.UserID] = result.TotalWeight
	}

	return totals, nil
}
func (r *WasteDropRequestItemRepository) prepareDateParams(monthStr string) (string, string) {
	if monthStr == "" {
		return "", ""
	}

	if parsed, err := time.Parse("2006-01", monthStr); err == nil {
		return monthStr, parsed.Format("2006-01-02 15:04:05")
	}

	return "", ""
}
