package model

type WasteCollectorResponse struct {
	ID               string        `json:"id"`
	UserID           string        `json:"user_id"`
	TotalWasteWeight float64       `json:"total_waste_weight"`
	User             *UserResponse `json:"user,omitempty"`
}

type WasteCollectorRequest struct {
	UserID           string   `json:"user_id"`
	TotalWasteWeight *float64 `json:"total_waste_weight,omitempty"`
}

type GetWasteCollectorRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type UpdateWasteCollectorRequest struct {
	ID               string   `json:"id" validate:"required,max=100"`
	TotalWasteWeight *float64 `json:"total_waste_weight,omitempty"`
}

type DeleteWasteCollectorRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
