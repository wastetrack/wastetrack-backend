package model

type CustomerResponse struct {
	ID            string        `json:"id"`
	UserID        string        `json:"user_id"`
	CarbonDeficit int64         `json:"carbon_deficit"`
	WaterSaved    int64         `json:"water_saved"`
	BagsStored    int64         `json:"bags_stored"`
	Trees         int64         `json:"trees"`
	User          *UserResponse `json:"user,omitempty"`
}

type CustomerRequest struct {
	UserID        string `json:"user_id"`
	CarbonDeficit *int64 `json:"carbon_deficit,omitempty"`
	WaterSaved    *int64 `json:"water_saved,omitempty"`
	BagsStored    *int64 `json:"bags_stored,omitempty"`
	Trees         *int64 `json:"trees,omitempty"`
}

type GetCustomerRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type UpdateCustomerRequest struct {
	ID            string `json:"id" validate:"required,max=100"`
	CarbonDeficit *int64 `json:"carbon_deficit,omitempty"`
	WaterSaved    *int64 `json:"water_saved,omitempty"`
	BagsStored    *int64 `json:"bags_stored,omitempty"`
	Trees         *int64 `json:"trees,omitempty"`
}

type DeleteCustomerRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
