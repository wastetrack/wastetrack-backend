package model

type StorageSimpleResponse struct {
	ID                    string  `json:"id"`
	UserID                string  `json:"user_id"`
	Length                float64 `json:"length"`
	Width                 float64 `json:"width"`
	Height                float64 `json:"height"`
	IsForRecycledMaterial bool    `json:"is_for_recycled_material"`
}

type StorageResponse struct {
	ID                    string  `json:"id"`
	UserID                string  `json:"user_id"`
	Length                float64 `json:"length"`
	Width                 float64 `json:"width"`
	Height                float64 `json:"height"`
	User                  *UserResponse
	IsForRecycledMaterial bool `json:"is_for_recycled_material"`
}
type StorageRequest struct {
	UserID                string  `json:"user_id"`
	Length                float64 `json:"length"`
	Width                 float64 `json:"width"`
	Height                float64 `json:"height"`
	IsForRecycledMaterial bool    `json:"is_for_recycled_material"`
}

type SearchStorageRequest struct {
	UserID                string
	IsForRecycledMaterial *bool
	MinLength             *float64
	MaxLength             *float64
	MinWidth              *float64
	MaxWidth              *float64
	MinHeight             *float64
	MaxHeight             *float64
	Page                  int
	Size                  int
}
type GetStorageRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateStorageRequest struct {
	ID                    string  `json:"id" validate:"required,max=100"`
	UserID                string  `json:"user_id"`
	Length                float64 `json:"length"`
	Width                 float64 `json:"width"`
	Height                float64 `json:"height"`
	IsForRecycledMaterial *bool   `json:"is_for_recycled_material"`
}

type DeleteStorageRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
