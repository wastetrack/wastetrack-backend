package model

type StorageSimpleResponse struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type StorageResponse struct {
	ID     string  `json:"id"`
	UserID string  `json:"user_id"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	User   *UserResponse
}
type StorageRequest struct {
	UserID string  `json:"user_id"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type SearchStorageRequest struct {
	UserID string  `json:"user_id"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Page   int     `json:"page,omitempty" validate:"min=1"`
	Size   int     `json:"size,omitempty" validate:"min=1,max=100"`
}
type GetStorageRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateStorageRequest struct {
	ID     string  `json:"id" validate:"required,max=100"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type DeleteStorageRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
