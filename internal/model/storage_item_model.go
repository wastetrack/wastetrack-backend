package model

import "time"

type StorageItemSimpleResponse struct {
	ID          string    `json:"id"`
	StorageID   string    `json:"storage_id"`
	WasteTypeID string    `json:"waste_type_id"`
	QuantityKgs float64   `json:"quantity_kgs"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StorageItemResponse struct {
	ID          string                 `json:"id"`
	StorageID   string                 `json:"storage_id"`
	WasteTypeID string                 `json:"waste_type_id"`
	QuantityKgs float64                `json:"quantity_kgs"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	Storage     *StorageSimpleResponse `json:"storage"`
	WasteType   *WasteTypeResponse     `json:"waste_type"`
}
type StorageItemRequest struct {
	StorageID   string  `json:"storage_id" validate:"required,max=100"`
	WasteTypeID string  `json:"waste_type_id" validate:"required,max=100"`
	QuantityKgs float64 `json:"quantity_kgs"`
}

type SearchStorageItemRequest struct {
	StorageID          string `json:"storage_id"`
	WasteTypeID        string `json:"waste_type_id"`
	OrderByQuantityKgs string `json:"order_by_quantity_kgs"`
	Page               int    `json:"page,omitempty" `
	Size               int    `json:"size,omitempty"`
}
type GetStorageItemRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateStorageItemRequest struct {
	ID     string  `json:"id" validate:"required,max=100"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type DeleteStorageItemRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
