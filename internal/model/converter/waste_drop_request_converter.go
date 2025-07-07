package converter

import (
	"reflect"

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
	}

	// Extract distance from additional fields if present (when using raw SQL with distance calculation)
	if distance := extractDistanceFromEntity(wasteDropRequest); distance != nil {
		response.Distance = distance
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
	}

	// Extract distance from additional fields if present (when using raw SQL with distance calculation)
	if distance := extractDistanceFromEntity(wasteDropRequest); distance != nil {
		response.Distance = distance
	}

	return response
}

// extractDistanceFromEntity extracts distance from entity if it was calculated in SQL query
// This uses reflection to access the "distance" field that might be added by GORM
func extractDistanceFromEntity(request *entity.WasteDropRequest) *float64 {
	// Use reflection to check if there's a distance field
	val := reflect.ValueOf(request).Elem()

	// Look for additional fields that GORM might have added
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if field.Name == "Distance" || field.Tag.Get("column") == "distance" {
			distanceVal := val.Field(i)
			if distanceVal.IsValid() && !distanceVal.IsZero() {
				if distanceVal.Kind() == reflect.Float64 {
					distance := distanceVal.Float()
					return &distance
				}
				if distanceVal.Kind() == reflect.Ptr && !distanceVal.IsNil() {
					distance := distanceVal.Elem().Float()
					return &distance
				}
			}
		}
	}

	// Alternative approach: check if GORM added additional columns
	// This would require accessing GORM's additional columns map
	// For now, return nil if no distance found
	return nil
}
