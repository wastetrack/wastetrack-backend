package model

import "time"

type WasteDropRequestItems struct {
	WasteTypeIDs []string `json:"waste_type_ids" validate:"required,min=1"`
	Quantities   []int64  `json:"quantities" validate:"required,min=1"`
}
type WasteDropRequestSimpleResponse struct {
	ID                   string            `json:"id"`
	DeliveryType         string            `json:"delivery_type"`
	CustomerID           string            `json:"customer_id"`
	UserPhoneNumber      string            `json:"user_phone_number,omitempty"`
	WasteBankID          string            `json:"waste_bank_id,omitempty"`
	AssignedCollectorID  string            `json:"assigned_collector_id,omitempty"`
	TotalPrice           int64             `json:"total_price"`
	ImageURL             string            `json:"image_url,omitempty"`
	Status               string            `json:"status"`
	AppointmentLocation  *LocationResponse `json:"appointment_location,omitempty"`
	AppointmentDate      string            `json:"appointment_date,omitempty"`
	AppointmentStartTime string            `json:"appointment_start_time,omitempty"`
	AppointmentEndTime   string            `json:"appointment_end_time,omitempty"`
	Notes                string            `json:"notes,omitempty"`
	CreatedAt            *time.Time        `json:"created_at"`
	UpdatedAt            *time.Time        `json:"updated_at"`
}

type WasteDropRequestResponse struct {
	ID                   string            `json:"id"`
	DeliveryType         string            `json:"delivery_type"`
	CustomerID           string            `json:"customer_id"`
	UserPhoneNumber      string            `json:"user_phone_number,omitempty"`
	WasteBankID          string            `json:"waste_bank_id,omitempty"`
	AssignedCollectorID  string            `json:"assigned_collector_id,omitempty"`
	TotalPrice           int64             `json:"total_price"`
	ImageURL             string            `json:"image_url,omitempty"`
	Status               string            `json:"status"`
	AppointmentLocation  *LocationResponse `json:"appointment_location,omitempty"`
	AppointmentDate      string            `json:"appointment_date,omitempty"`
	AppointmentStartTime string            `json:"appointment_start_time,omitempty"`
	AppointmentEndTime   string            `json:"appointment_end_time,omitempty"`
	Notes                string            `json:"notes,omitempty"`
	CreatedAt            *time.Time        `json:"created_at"`
	UpdatedAt            *time.Time        `json:"updated_at"`
	Customer             *UserResponse     `json:"customer"`
	WasteBank            *UserResponse     `json:"waste_bank"`
	AssignedCollector    *UserResponse     `json:"assigned_collector"`
}

type WasteDropRequestRequest struct {
	DeliveryType         string                 `json:"delivery_type" validate:"required,max=100"`
	CustomerID           string                 `json:"customer_id" validate:"required,max=100"`
	UserPhoneNumber      string                 `json:"user_phone_number,omitempty"`
	WasteBankID          string                 `json:"waste_bank_id,omitempty"`
	AssignedCollectorID  string                 `json:"assigned_collector_id,omitempty"`
	TotalPrice           int64                  `json:"total_price"`
	ImageURL             string                 `json:"image_url,omitempty"`
	AppointmentLocation  *LocationRequest       `json:"appointment_location,omitempty"`
	AppointmentDate      string                 `json:"appointment_date,omitempty"`
	AppointmentStartTime string                 `json:"appointment_start_time,omitempty"`
	AppointmentEndTime   string                 `json:"appointment_end_time,omitempty"`
	Notes                string                 `json:"notes,omitempty"`
	Items                *WasteDropRequestItems `json:"items" validate:"required"`
}

type SearchWasteDropRequest struct {
	DeliveryType         string `json:"delivery_type"`
	CustomerID           string `json:"customer_id"`
	WasteBankID          string `json:"waste_bank_id,omitempty"`
	AssignedCollectorID  string `json:"assigned_collector_id,omitempty"`
	AppointmentDate      string `json:"appointment_date,omitempty"`
	AppointmentStartTime string `json:"appointment_start_time,omitempty"`
	AppointmentEndTime   string `json:"appointment_end_time,omitempty"`
	Status               string `json:"status"`
	Page                 int    `json:"page,omitempty" validate:"min=1"`
	Size                 int    `json:"size,omitempty" validate:"min=1,max=100"`
}
type GetWasteDropRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateWasteDropRequest struct {
	ID                  string `json:"id" validate:"required,max=100"`
	DeliveryType        string `json:"delivery_type"`
	AssignedCollectorID string `json:"assigned_collector_id,omitempty"`
	Status              string `json:"status"`
}

type DeleteWasteDropRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
