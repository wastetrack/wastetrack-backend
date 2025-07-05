package model

type CollectorManagementSimpleResponse struct {
	ID          string `json:"id"`
	WasteBankID string `json:"waste_bank_id"`
	CollectorID string `json:"collector_id"`
	Status      string `json:"status"`
}

type CollectorManagementResponse struct {
	ID          string        `json:"id"`
	WasteBankID string        `json:"waste_bank_id"`
	WasteBank   *UserResponse `json:"waste_bank"`
	CollectorID string        `json:"collector_id"`
	Collector   *UserResponse `json:"collector"`
	Status      string        `json:"status"`
}

type CollectorManagementRequest struct {
	WasteBankID string `json:"waste_bank_id" validate:"required,max=100"`
	CollectorID string `json:"collector_id" validate:"required,max=100"`
	Status      string `json:"status"`
}

type SearchCollectorManagementRequest struct {
	WasteBankID string `json:"waste_bank_id"`
	CollectorID string `json:"collector_id"`
	Status      string `json:"status"`
	Page        int    `json:"page,omitempty" validate:"min=1"`
	Size        int    `json:"size,omitempty" validate:"min=1,max=100"`
}
type GetCollectorManagementRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateCollectorManagementRequest struct {
	ID          string `json:"id" validate:"required,max=100"`
	WasteBankID string `json:"waste_bank_id"`
	CollectorID string `json:"collector_id"`
	Status      string `json:"status"`
}

type DeleteCollectorManagementRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
