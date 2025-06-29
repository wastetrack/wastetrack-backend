package model

type IndustryResponse struct {
	ID                  string        `json:"id"`
	UserID              string        `json:"user_id"`
	TotalWasteWeight    float64       `json:"total_waste_weight"`
	TotalRecycledWeight float64       `json:"total_recycled_weight"`
	User                *UserResponse `json:"user,omitempty"`
}

type IndustryRequest struct {
	UserID              string   `json:"user_id"`
	TotalWasteWeight    *float64 `json:"total_waste_weight,omitempty"`
	TotalRecycledWeight *float64 `json:"total_recycled_weight,omitempty"`
}

type GetIndustryRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type UpdateIndustryRequest struct {
	ID                  string   `json:"id" validate:"required,max=100"`
	TotalWasteWeight    *float64 `json:"total_waste_weight,omitempty"`
	TotalRecycledWeight *float64 `json:"total_recycled_weight,omitempty"`
}

type DeleteIndustryRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
