package model

type WasteBankResponse struct {
	ID               string        `json:"id"`
	UserID           string        `json:"user_id"`
	User             *UserResponse `json:"user,omitempty"`
	TotalWasteWeight int64         `json:"total_waste_weight"`
	TotalWorkers     int64         `json:"total_workers"`
	OpenTime         string        `json:"open_time"`
	CloseTime        string        `json:"close_time"`
}

type WasteBankRequest struct {
	UserID           string  `json:"user_id"`
	TotalWasteWeight *int64  `json:"total_waste_weight,omitempty"`
	TotalWorkers     *int64  `json:"total_workers,omitempty"`
	OpenTime         *string `json:"open_time,omitempty"`
	CloseTime        *string `json:"close_time,omitempty"`
}

type WasteBankUpdateRequest struct {
	ID               string  `json:"id"validate:"required,max=100"`
	TotalWasteWeight *int64  `json:"total_waste_weight,omitempty"`
	TotalWorkers     *int64  `json:"total_workers,omitempty"`
	OpenTime         *string `json:"open_time,omitempty"`
	CloseTime        *string `json:"close_time,omitempty"`
}

type DeleteWasteBankRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
