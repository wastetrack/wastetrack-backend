package model

import "time"

type WasteTransferRequestItems struct {
	WasteTypeIDs         []string  `json:"waste_type_ids" validate:"required,min=1"`
	OfferingWeights      []float64 `json:"offering_weights" validate:"required,min=1"`
	OfferingPricesPerKgs []int64   `json:"offering_prices_per_kgs" validate:"required,min=1"`
}

type WasteTransferRequestRequest struct {
	SourceUserID           string                     `json:"source_user_id" validate:"required,max=100"`
	DestinationUserID      string                     `json:"destination_user_id" validate:"required,max=100"`
	FormType               string                     `json:"form_type" validate:"required,max=100"`
	Status                 string                     `json:"status" `
	ImageURL               string                     `json:"image_url,omitempty"`
	Notes                  string                     `json:"notes,omitempty"`
	SourcePhoneNumber      string                     `json:"source_phone_number" validate:"required,max=100"`
	DestinationPhoneNumber string                     `json:"destination_phone_number" validate:"required,max=100"`
	AppointmentDate        string                     `json:"appointment_date" validate:"required"`
	AppointmentStartTime   string                     `json:"appointment_start_time,omitempty"`
	AppointmentEndTime     string                     `json:"appointment_end_time,omitempty"`
	AppointmentLocation    *LocationRequest           `json:"appointment_location,omitempty"`
	Items                  *WasteTransferRequestItems `json:"items" validate:"required"`
}

type AssignCollectorByWasteTypeRequest struct {
	ID                  string                            `json:"id" validate:"required,max=100"`
	AssignedCollectorID string                            `json:"assigned_collector_id"`
	WasteTypes          []AssignCollectorWasteTypeRequest `json:"items" validate:"required,min=1"`
}

type AssignCollectorWasteTypeRequest struct {
	WasteTypeID         string  `json:"waste_type_id" validate:"required,max=100"`
	AcceptedWeight      float64 `json:"accepted_weight" validate:"required,min=0"`
	AcceptedPricePerKgs int64   `json:"accepted_price_per_kgs" validate:"required,min=0"`
}

type CompleteWasteTransferRequestItems struct {
	WasteTypeIDs []string  `json:"waste_type_ids" validate:"required,min=1"`
	Weights      []float64 `json:"weights" validate:"required,min=1"`
}

type CompleteWasteTransferRequest struct {
	ID    string                             `json:"id" validate:"required,max=100"`
	Items *CompleteWasteTransferRequestItems `json:"items" validate:"required"`
}

type RecycleWasteTransferRequestItems struct {
	WasteTypeIDs []string  `json:"waste_type_ids" validate:"required,min=1"`
	Weights      []float64 `json:"weights" validate:"required,min=1"`
}

type RecycleWasteTransferRequest struct {
	ID    string                            `json:"id" validate:"required,max=100"`
	Items *RecycleWasteTransferRequestItems `json:"items" validate:"required"`
}

// Response models
type WasteTransferRequestSimpleResponse struct {
	ID                     string            `json:"id"`
	SourceUserID           string            `json:"source_user_id"`
	DestinationUserID      string            `json:"destination_user_id"`
	AssignedCollectorID    string            `json:"assigned_collector_id,omitempty"`
	FormType               string            `json:"form_type"`
	TotalWeight            float64           `json:"total_weight"`
	TotalPrice             int64             `json:"total_price"`
	Status                 string            `json:"status"`
	ImageURL               string            `json:"image_url,omitempty"`
	Notes                  string            `json:"notes,omitempty"`
	SourcePhoneNumber      string            `json:"source_phone_number"`
	DestinationPhoneNumber string            `json:"destination_phone_number"`
	AppointmentDate        string            `json:"appointment_date,omitempty"`
	AppointmentStartTime   string            `json:"appointment_start_time,omitempty"`
	AppointmentEndTime     string            `json:"appointment_end_time,omitempty"`
	AppointmentLocation    *LocationResponse `json:"appointment_location,omitempty"`
	CreatedAt              *time.Time        `json:"created_at"`
	UpdatedAt              *time.Time        `json:"updated_at"`
	Distance               *float64          `json:"distance,omitempty"`
	IsDeleted              bool              `json:"is_deleted"`
}

type WasteTransferRequestResponse struct {
	ID                     string                              `json:"id"`
	SourceUserID           string                              `json:"source_user_id"`
	DestinationUserID      string                              `json:"destination_user_id"`
	AssignedCollectorID    string                              `json:"assigned_collector_id,omitempty"` // NEW
	FormType               string                              `json:"form_type"`
	TotalWeight            float64                             `json:"total_weight"`
	TotalPrice             int64                               `json:"total_price"`
	Status                 string                              `json:"status"`
	ImageURL               string                              `json:"image_url,omitempty"`
	Notes                  string                              `json:"notes,omitempty"`
	SourcePhoneNumber      string                              `json:"source_phone_number"`
	DestinationPhoneNumber string                              `json:"destination_phone_number"`
	AppointmentDate        string                              `json:"appointment_date,omitempty"`
	AppointmentStartTime   string                              `json:"appointment_start_time,omitempty"`
	AppointmentEndTime     string                              `json:"appointment_end_time,omitempty"`
	AppointmentLocation    *LocationResponse                   `json:"appointment_location,omitempty"`
	CreatedAt              *time.Time                          `json:"created_at"`
	UpdatedAt              *time.Time                          `json:"updated_at"`
	SourceUser             *UserResponse                       `json:"source_user"`
	DestinationUser        *UserResponse                       `json:"destination_user"`
	AssignedCollector      *UserResponse                       `json:"assigned_collector,omitempty"` // NEW
	Items                  []WasteTransferItemOfferingResponse `json:"items"`
	Distance               *float64                            `json:"distance,omitempty"`
	IsDeleted              bool                                `json:"is_deleted"`
}

// Search and operation models
type SearchWasteTransferRequest struct {
	SourceUserID         string   `json:"source_user_id"`
	DestinationUserID    string   `json:"destination_user_id"`
	AssignedCollectorID  string   `json:"assigned_collector_id"` // NEW
	FormType             string   `json:"form_type"`
	Status               string   `json:"status"`
	AppointmentDate      string   `json:"appointment_date,omitempty"`
	AppointmentStartTime string   `json:"appointment_start_time,omitempty"`
	AppointmentEndTime   string   `json:"appointment_end_time,omitempty"`
	Latitude             *float64 `json:"latitude,omitempty"`
	Longitude            *float64 `json:"longitude,omitempty"`
	OrderDir             string   `json:"order_dir,omitempty"`
	IsDeleted            *bool    `json:"is_deleted"`
	Page                 int      `json:"page,omitempty" validate:"min=1"`
	Size                 int      `json:"size,omitempty" validate:"min=1,max=100"`
}

type GetWasteTransferRequest struct {
	ID        string   `json:"id" validate:"required,max=100"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type UpdateWasteTransferRequest struct {
	ID                   string `json:"id" validate:"required,max=100"`
	FormType             string `json:"form_type"`
	Status               string `json:"status"`
	AppointmentDate      string `json:"appointment_date,omitempty"`
	AppointmentStartTime string `json:"appointment_start_time,omitempty"`
	AppointmentEndTime   string `json:"appointment_end_time,omitempty"`
}

type DeleteWasteTransferRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
