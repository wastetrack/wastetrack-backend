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

// WasteTransferItemOfferingRepository
type WasteTransferItemOfferingRepository struct {
	Repository[entity.WasteTransferItemOffering]
	Log *logrus.Logger
}

func NewWasteTransferItemOfferingRepository(log *logrus.Logger) *WasteTransferItemOfferingRepository {
	return &WasteTransferItemOfferingRepository{
		Log: log,
	}
}

func (r *WasteTransferItemOfferingRepository) FindByID(db *gorm.DB, item *entity.WasteTransferItemOffering, id string) error {
	return db.Where("id = ?", id).
		Preload("TransferForm").
		Preload("TransferForm.SourceUser").
		Preload("TransferForm.DestinationUser").
		Preload("WasteType").
		Preload("WasteType.WasteCategory").
		First(item).Error
}

func (r *WasteTransferItemOfferingRepository) Search(db *gorm.DB, request *model.SearchWasteTransferItemOfferingRequest) ([]entity.WasteTransferItemOffering, int64, error) {
	var items []entity.WasteTransferItemOffering

	// Apply filters and pagination
	if err := db.Scopes(r.FilterWasteTransferItemOffering(request)).
		Preload("WasteType").
		Preload("WasteType.WasteCategory").
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	// Count total records with same filters
	var total int64 = 0
	if err := db.Model(&entity.WasteTransferItemOffering{}).
		Scopes(r.FilterWasteTransferItemOffering(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *WasteTransferItemOfferingRepository) FilterWasteTransferItemOffering(request *model.SearchWasteTransferItemOfferingRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if transferFormID := request.TransferFormID; transferFormID != "" {
			tx = tx.Where("transfer_request_id = ?", transferFormID)
		}

		// Filter by waste category ID - need to join waste_types table since category is in waste_type
		if wasteCategoryID := request.WasteCategoryID; wasteCategoryID != "" {
			tx = tx.Joins("JOIN waste_types ON waste_transfer_items.waste_type_id = waste_types.id").
				Where("waste_types.category_id = ?", wasteCategoryID)
		}

		if wasteTypeID := request.WasteTypeID; wasteTypeID != "" {
			tx = tx.Where("waste_type_id = ?", wasteTypeID)
		}

		return tx
	}
}

func (r *WasteTransferItemOfferingRepository) FindByTransferFormID(db *gorm.DB, transferFormID uuid.UUID) ([]entity.WasteTransferItemOffering, error) {
	var items []entity.WasteTransferItemOffering
	err := db.Where("transfer_request_id = ?", transferFormID).
		Preload("WasteType").
		Find(&items).Error
	return items, err
}

// GetWasteTransferVerifiedWeightByMonth returns monthly total verified weight from waste transfer requests
func (r *WasteTransferItemOfferingRepository) GetTransferVerifiedWeightByMonth(db *gorm.DB, request *model.GovernmentDashboardRequest) (map[string]float64, error) {
	weightByMonth := make(map[string]float64)

	query := db.Table("waste_transfer_items").
		Select("DATE_TRUNC('month', waste_transfer_requests.created_at) as month, COALESCE(SUM(waste_transfer_items.verified_weight), 0) as total_weight").
		Joins("JOIN waste_transfer_requests ON waste_transfer_items.transfer_request_id = waste_transfer_requests.id").
		Where("waste_transfer_requests.is_deleted = ?", false).
		Group("DATE_TRUNC('month', waste_transfer_requests.created_at)").
		Order("month")

	// Apply date filters
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

	// Apply location filters if provided (filter by destination user location)
	if request.Province != "" || request.City != "" {
		query = query.Joins("JOIN users AS destination_users ON waste_transfer_requests.destination_user_id = destination_users.id")

		if request.Province != "" {
			query = query.Where("destination_users.province ILIKE ?", "%"+request.Province+"%")
		}
		if request.City != "" {
			query = query.Where("destination_users.city ILIKE ?", "%"+request.City+"%")
		}
	}

	// Execute query
	rows, err := query.Rows()
	if err != nil {
		return nil, fmt.Errorf("database error occurred while getting waste transfer verified weights")
	}
	defer rows.Close()

	for rows.Next() {
		var month time.Time
		var totalWeight float64

		if err := rows.Scan(&month, &totalWeight); err != nil {
			return nil, fmt.Errorf("error scanning waste transfer verified weight data")
		}

		monthStr := month.Format("2006-01")
		weightByMonth[monthStr] = totalWeight
	}

	return weightByMonth, nil
}

// GetWasteTransferTrendsByRole returns monthly waste transfer trends grouped by destination user role
func (r *WasteTransferItemOfferingRepository) GetWasteTransferTrendsByRole(db *gorm.DB, request *model.GovernmentDashboardRequest) (map[string]map[string]float64, error) {
	// Structure: map[month]map[role]weight
	trends := make(map[string]map[string]float64)

	// Query for waste transfer items - group by month and destination user role
	wasteTransferQuery := db.Table("waste_transfer_items").
		Select("DATE_TRUNC('month', waste_transfer_requests.created_at) as month, destination_users.role, COALESCE(SUM(waste_transfer_items.verified_weight), 0) as total_weight").
		Joins("JOIN waste_transfer_requests ON waste_transfer_items.transfer_request_id = waste_transfer_requests.id").
		Joins("JOIN users AS destination_users ON waste_transfer_requests.destination_user_id = destination_users.id").
		Where("waste_transfer_requests.is_deleted = ?", false).
		Group("DATE_TRUNC('month', waste_transfer_requests.created_at), destination_users.role").
		Order("month")

	// Apply date filters for waste transfer
	if request.StartMonth != "" {
		startDate, err := time.Parse("2006-01", request.StartMonth)
		if err != nil {
			return nil, fmt.Errorf("invalid start_month format: %v", err)
		}
		wasteTransferQuery = wasteTransferQuery.Where("waste_transfer_requests.created_at >= ?", startDate)
	}

	if request.EndMonth != "" {
		endDate, err := time.Parse("2006-01", request.EndMonth)
		if err != nil {
			return nil, fmt.Errorf("invalid end_month format: %v", err)
		}
		endDate = endDate.AddDate(0, 1, -1).Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		wasteTransferQuery = wasteTransferQuery.Where("waste_transfer_requests.created_at <= ?", endDate)
	}

	// Apply location filters for waste transfer if provided
	if request.Province != "" || request.City != "" {
		if request.Province != "" {
			wasteTransferQuery = wasteTransferQuery.Where("destination_users.province ILIKE ?", "%"+request.Province+"%")
		}
		if request.City != "" {
			wasteTransferQuery = wasteTransferQuery.Where("destination_users.city ILIKE ?", "%"+request.City+"%")
		}
	}

	// Execute waste transfer query
	wasteTransferRows, err := wasteTransferQuery.Rows()
	if err != nil {
		return nil, fmt.Errorf("database error occurred while getting waste transfer collection trends")
	}
	defer wasteTransferRows.Close()

	// Process waste transfer results
	for wasteTransferRows.Next() {
		var month time.Time
		var role string
		var totalWeight float64

		if err := wasteTransferRows.Scan(&month, &role, &totalWeight); err != nil {
			return nil, fmt.Errorf("error scanning waste transfer trend data")
		}

		monthStr := month.Format("2006-01")
		if trends[monthStr] == nil {
			trends[monthStr] = make(map[string]float64)
		}
		trends[monthStr][role] += totalWeight
	}

	return trends, nil
}
