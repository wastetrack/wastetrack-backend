package repository

import (
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
