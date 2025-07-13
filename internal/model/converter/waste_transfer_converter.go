package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

// WasteTransferRequest converters
func WasteTransferRequestToSimpleResponse(request *entity.WasteTransferRequest) *model.WasteTransferRequestSimpleResponse {
	var startTime, endTime, appointmentDate string

	if !request.AppointmentStartTime.IsZero() {
		startTime = request.AppointmentStartTime.Format("15:04:05Z07:00")
	}

	if !request.AppointmentEndTime.IsZero() {
		endTime = request.AppointmentEndTime.Format("15:04:05Z07:00")
	}

	if !request.AppointmentDate.IsZero() {
		appointmentDate = request.AppointmentDate.Format("2006-01-02")
	}
	var location *model.LocationResponse
	if request.AppointmentLocation != nil {
		location = &model.LocationResponse{
			Latitude:  request.AppointmentLocation.Lat,
			Longitude: request.AppointmentLocation.Lng,
		}
	}
	return &model.WasteTransferRequestSimpleResponse{
		ID:                     request.ID.String(),
		SourceUserID:           request.SourceUserID.String(),
		DestinationUserID:      request.DestinationUserID.String(),
		AssignedCollectorID:    request.AssignedCollectorID.String(),
		FormType:               request.FormType,
		TotalWeight:            request.TotalWeight,
		TotalPrice:             request.TotalPrice,
		Status:                 request.Status,
		ImageURL:               request.ImageURL,
		Notes:                  request.Notes,
		SourcePhoneNumber:      request.SourcePhoneNumber,
		DestinationPhoneNumber: request.DestinationPhoneNumber,
		AppointmentDate:        appointmentDate,
		AppointmentStartTime:   startTime,
		AppointmentEndTime:     endTime,
		AppointmentLocation:    location,
		CreatedAt:              &request.CreatedAt,
		UpdatedAt:              &request.UpdatedAt,
		Distance:               request.Distance,
	}
}

func WasteTransferRequestToResponse(request *entity.WasteTransferRequest) *model.WasteTransferRequestResponse {
	var startTime, endTime, appointmentDate string

	if !request.AppointmentStartTime.IsZero() {
		startTime = request.AppointmentStartTime.Format("15:04:05Z07:00")
	}

	if !request.AppointmentEndTime.IsZero() {
		endTime = request.AppointmentEndTime.Format("15:04:05Z07:00")
	}

	if !request.AppointmentDate.IsZero() {
		appointmentDate = request.AppointmentDate.Format("2006-01-02")
	}

	var sourceUser *model.UserResponse
	if request.SourceUser.ID != uuid.Nil {
		sourceUser = UserToResponse(&request.SourceUser)
	}

	var destinationUser *model.UserResponse
	if request.DestinationUser.ID != uuid.Nil {
		destinationUser = UserToResponse(&request.DestinationUser)
	}
	var location *model.LocationResponse
	if request.AppointmentLocation != nil {
		location = &model.LocationResponse{
			Latitude:  request.AppointmentLocation.Lat,
			Longitude: request.AppointmentLocation.Lng,
		}
	}
	// Convert items
	var items []model.WasteTransferItemOfferingResponse
	for _, item := range request.Items {
		items = append(items, *WasteTransferItemOfferingToResponse(&item))
	}

	return &model.WasteTransferRequestResponse{
		ID:                     request.ID.String(),
		SourceUserID:           request.SourceUserID.String(),
		DestinationUserID:      request.DestinationUserID.String(),
		AssignedCollectorID:    request.AssignedCollectorID.String(),
		FormType:               request.FormType,
		TotalWeight:            request.TotalWeight,
		TotalPrice:             request.TotalPrice,
		Status:                 request.Status,
		ImageURL:               request.ImageURL,
		Notes:                  request.Notes,
		SourcePhoneNumber:      request.SourcePhoneNumber,
		DestinationPhoneNumber: request.DestinationPhoneNumber,
		AppointmentDate:        appointmentDate,
		AppointmentStartTime:   startTime,
		AppointmentEndTime:     endTime,
		AppointmentLocation:    location,
		CreatedAt:              &request.CreatedAt,
		UpdatedAt:              &request.UpdatedAt,
		SourceUser:             sourceUser,
		DestinationUser:        destinationUser,
		Items:                  items,
		Distance:               request.Distance,
	}
}
