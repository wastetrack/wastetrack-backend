package model

import "time"

type WasteDropRequestListResponse struct {
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
	WasteDropRequestListResponse
	Customer          *UserResponse                   `json:"customer"`
	WasteBank         *UserResponse                   `json:"waste_bank"`
	AssignedCollector *UserResponse                   `json:"assigned_collector"`
	Items             []*WasteDropRequestItemResponse `json:"items,omitempty"`
}
