package converter

import (
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func WasteDropRequestToSimpleResponse(wasteDropRequest *entity.WasteDropRequest) *model.WasteDropRequestSimpleResponse {
	var startTime, endTime, appointmentDate string

	if !wasteDropRequest.AppointmentStartTime.IsZero() {
		startTime = wasteDropRequest.AppointmentStartTime.Format("15:04:05Z07:00")
	}

	if !wasteDropRequest.AppointmentEndTime.IsZero() {
		endTime = wasteDropRequest.AppointmentEndTime.Format("15:04:05Z07:00")
	}
	if !wasteDropRequest.AppointmentDate.IsZero() {
		appointmentDate = wasteDropRequest.AppointmentDate.Format("2006-01-02")
	}

	var location *model.LocationResponse
	if wasteDropRequest.AppointmentLocation != nil {
		location = &model.LocationResponse{
			Latitude:  wasteDropRequest.AppointmentLocation.Lat,
			Longitude: wasteDropRequest.AppointmentLocation.Lng,
		}
	}

	// Handle potentially nil UUID pointers
	var wasteBankID, assignedCollectorID string
	if wasteDropRequest.WasteBankID != nil {
		wasteBankID = wasteDropRequest.WasteBankID.String()
	}
	if wasteDropRequest.AssignedCollectorID != nil {
		assignedCollectorID = wasteDropRequest.AssignedCollectorID.String()
	}

	response := &model.WasteDropRequestSimpleResponse{
		ID:                   wasteDropRequest.ID.String(),
		DeliveryType:         wasteDropRequest.DeliveryType,
		CustomerID:           wasteDropRequest.CustomerID.String(),
		UserPhoneNumber:      wasteDropRequest.UserPhoneNumber,
		WasteBankID:          wasteBankID,
		AssignedCollectorID:  assignedCollectorID,
		TotalPrice:           wasteDropRequest.TotalPrice,
		ImageURL:             wasteDropRequest.ImageURL,
		Status:               wasteDropRequest.Status,
		AppointmentLocation:  location,
		AppointmentDate:      appointmentDate,
		AppointmentStartTime: startTime,
		AppointmentEndTime:   endTime,
		Notes:                wasteDropRequest.Notes,
		CreatedAt:            &wasteDropRequest.CreatedAt,
		UpdatedAt:            &wasteDropRequest.UpdatedAt,
		Distance:             wasteDropRequest.Distance, // Now directly accessible
		IsDeleted:            wasteDropRequest.IsDeleted,
	}

	return response
}

func WasteDropRequestToResponse(wasteDropRequest *entity.WasteDropRequest) *model.WasteDropRequestResponse {
	var startTime, endTime, appointmentDate string

	if !wasteDropRequest.AppointmentStartTime.IsZero() {
		startTime = wasteDropRequest.AppointmentStartTime.Format("15:04:05Z07:00")
	}

	if !wasteDropRequest.AppointmentEndTime.IsZero() {
		endTime = wasteDropRequest.AppointmentEndTime.Format("15:04:05Z07:00")
	}
	if !wasteDropRequest.AppointmentDate.IsZero() {
		appointmentDate = wasteDropRequest.AppointmentDate.Format("2006-01-02")
	}

	var location *model.LocationResponse
	if wasteDropRequest.AppointmentLocation != nil {
		location = &model.LocationResponse{
			Latitude:  wasteDropRequest.AppointmentLocation.Lat,
			Longitude: wasteDropRequest.AppointmentLocation.Lng,
		}
	}

	var assignedCollector *model.UserResponse
	if wasteDropRequest.AssignedCollector != nil {
		assignedCollector = UserToResponse(wasteDropRequest.AssignedCollector)
	}

	// Handle potentially nil UUID pointers
	var wasteBankID, assignedCollectorID string
	if wasteDropRequest.WasteBankID != nil {
		wasteBankID = wasteDropRequest.WasteBankID.String()
	}
	if wasteDropRequest.AssignedCollectorID != nil {
		assignedCollectorID = wasteDropRequest.AssignedCollectorID.String()
	}

	response := &model.WasteDropRequestResponse{
		ID:                   wasteDropRequest.ID.String(),
		DeliveryType:         wasteDropRequest.DeliveryType,
		CustomerID:           wasteDropRequest.CustomerID.String(),
		UserPhoneNumber:      wasteDropRequest.UserPhoneNumber,
		WasteBankID:          wasteBankID,
		AssignedCollectorID:  assignedCollectorID,
		TotalPrice:           wasteDropRequest.TotalPrice,
		ImageURL:             wasteDropRequest.ImageURL,
		Status:               wasteDropRequest.Status,
		AppointmentLocation:  location,
		AppointmentDate:      appointmentDate,
		AppointmentStartTime: startTime,
		AppointmentEndTime:   endTime,
		Notes:                wasteDropRequest.Notes,
		CreatedAt:            &wasteDropRequest.CreatedAt,
		UpdatedAt:            &wasteDropRequest.UpdatedAt,
		Customer:             UserToResponse(&wasteDropRequest.Customer),
		WasteBank:            UserToResponse(&wasteDropRequest.WasteBank),
		AssignedCollector:    assignedCollector,
		Distance:             wasteDropRequest.Distance, // Now directly accessible
		IsDeleted:            wasteDropRequest.IsDeleted,
	}

	return response
}
