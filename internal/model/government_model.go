package model

type GovernmentResponse struct {
	ID     string        `json:"id"`
	UserID string        `json:"user_id"`
	User   *UserResponse `json:"user,omitempty"`
}

type GovernmentRequest struct {
	UserID string `json:"user_id"`
}

type GetGovernmentRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type UpdateGovernmentRequest struct {
	ID string `json:"id" validate:"required,max=100"`
	// Add fields as needed when government profile gets more properties
}

type DeleteGovernmentRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
