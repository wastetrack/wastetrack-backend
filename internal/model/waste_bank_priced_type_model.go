package model

type WasteBankPricedTypeSimpleResponse struct {
	ID                string             `json:"id"`
	WasteBankID       string             `json:"waste_bank_id"`
	WasteTypeID       string             `json:"waste_type_id"`
	CustomPricePerKgs int64              `json:"custom_price_per_kgs"`
	CreatedAt         string             `json:"created_at"`
	UpdatedAt         string             `json:"updated_at"`
	WasteType         *WasteTypeResponse `json:"waste_type"`
}
type WasteBankPricedTypeResponse struct {
	ID                string             `json:"id"`
	WasteBankID       string             `json:"waste_bank_id"`
	WasteTypeID       string             `json:"waste_type_id"`
	CustomPricePerKgs int64              `json:"custom_price_per_kgs"`
	CreatedAt         string             `json:"created_at"`
	UpdatedAt         string             `json:"updated_at"`
	WasteBank         *UserResponse      `json:"waste_bank"`
	WasteType         *WasteTypeResponse `json:"waste_type"`
}

type WasteBankPricedTypeRequest struct {
	WasteBankID       string `json:"waste_bank_id"`
	WasteTypeID       string `json:"waste_type_id"`
	CustomPricePerKgs int64  `json:"custom_price_per_kgs"`
}

type WasteBankPricedTypeBatchRequest struct {
	Items []WasteBankPricedTypeRequest `json:"items" validate:"required,dive"`
}

type SearchWasteBankPricedTypeRequest struct {
	WasteBankID string `json:"waste_bank_id"`
	WasteTypeID string `json:"waste_type_id"`
	Page        int    `json:"page,omitempty" validate:"min=1"`
	Size        int    `json:"size,omitempty" validate:"min=1,max=100"`
}
type GetWasteBankPricedTypeRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateWasteBankPricedTypeRequest struct {
	ID                string `json:"id" validate:"required,max=100"`
	CustomPricePerKgs int64  `json:"custom_price_per_kgs"`
}

type DeleteWasteBankPricedTypeRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
