package model

import "time"

type StorageItemSimpleResponse struct {
	ID          string    `json:"id"`
	StorageID   string    `json:"storage_id"`
	WasteTypeID string    `json:"waste_type_id"`
	WeightKgs   float64   `json:"weight_kgs"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StorageItemResponse struct {
	ID          string                 `json:"id"`
	StorageID   string                 `json:"storage_id"`
	WasteTypeID string                 `json:"waste_type_id"`
	WeightKgs   float64                `json:"weight_kgs"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Storage     *StorageSimpleResponse `json:"storage"`
	WasteType   *WasteTypeResponse     `json:"waste_type"`
}
type StorageItemRequest struct {
	StorageID   string  `json:"storage_id" validate:"required,max=100"`
	WasteTypeID string  `json:"waste_type_id" validate:"required,max=100"`
	WeightKgs   float64 `json:"weight_kgs"`
}

type SearchStorageItemRequest struct {
	StorageID        string `json:"storage_id"`
	WasteTypeID      string `json:"waste_type_id"`
	OrderByWeightKgs string `json:"order_by_weight_kgs"`
	Page             int    `json:"page,omitempty" `
	Size             int    `json:"size,omitempty"`
}
type GetStorageItemRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateStorageItemRequest struct {
	ID        string  `json:"id" validate:"required,max=100"`
	UserID    string  `json:"user_id" validate:"required,max=100"`
	StorageID string  `json:"storage_id" validate:"required,max=100"`
	Weight    float64 `json:"weight_kgs"`
}

type DeleteStorageItemRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
