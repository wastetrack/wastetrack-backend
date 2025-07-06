package repository

import (
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

func (r *WasteTransferRequestRepository) FindByID(db *gorm.DB, wasteTransferRequest *entity.WasteTransferRequest, id string) error {
	return db.Where("id = ?", id).
		Preload("SourceUser").
		Preload("DestinationUser").
		Preload("Items").
		Preload("Items.WasteType").
		First(wasteTransferRequest).Error
}

func (r *WasteTransferRequestRepository) Search(db *gorm.DB, request *model.SearchWasteTransferRequest) ([]entity.WasteTransferRequest, int64, error) {
	var wasteTransferRequests []entity.WasteTransferRequest

	// Apply filters, pagination and preload related entities
	if err := db.Scopes(r.FilterWasteTransferRequest(request)).
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(&wasteTransferRequests).Error; err != nil {
		return nil, 0, err
	}

	// Count total records with same filters
	var total int64 = 0
	if err := db.Model(&entity.WasteTransferRequest{}).
		Scopes(r.FilterWasteTransferRequest(request)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return wasteTransferRequests, total, nil
}

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

func (r *WasteTransferRequestRepository) UpdateStatus(db *gorm.DB, id string, status string) error {
	return db.Model(&entity.WasteTransferRequest{}).
		Where("id = ?", id).
		Update("status", status).Error
}
