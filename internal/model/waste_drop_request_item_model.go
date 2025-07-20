package model

type WasteDropRequestItems struct {
	WasteTypeIDs []string `json:"waste_type_ids" validate:"required,min=1"`
	Quantities   []int64  `json:"quantities" validate:"required,min=1"`
}

type CompleteWasteDropRequestItems struct {
	WasteTypeIDs []string  `json:"waste_type_ids" validate:"required,min=1"`
	Weights      []float64 `json:"weights" validate:"required,min=1"`
}

type CompleteWasteDropRequest struct {
	ID    string                         `json:"id" validate:"required,max=100"`
	Items *CompleteWasteDropRequestItems `json:"items" validate:"required"`
}
type WasteDropRequestItemSimpleResponse struct {
	ID                  string  `json:"id"`
	RequestID           string  `json:"request_id"`
	WasteTypeID         string  `json:"waste_type_id"`
	Quantity            int64   `json:"quantity"`
	VerifiedWeight      float64 `json:"verified_weight"`
	VerifiedPricePerKgs int64   `json:"verified_price_per_kgs"`
	VerifiedSubtotal    int64   `json:"verified_subtotal"`
}
type WasteDropRequestItemResponse struct {
	ID                  string  `json:"id"`
	RequestID           string  `json:"request_id"`
	WasteTypeID         string  `json:"waste_type_id"`
	Quantity            int64   `json:"quantity"`
	VerifiedWeight      float64 `json:"verified_weight"`
	VerifiedPricePerKgs int64   `json:"verified_price_per_kgs"`
	VerifiedSubtotal    int64   `json:"verified_subtotal"`
	Request             *WasteDropRequestSimpleResponse
	WasteType           *WasteTypeResponse
}

type WasteDropRequestItemRequest struct {
	RequestID           string  `json:"request_id"`
	WasteTypeID         string  `json:"waste_type_id"`
	Quantity            int64   `json:"quantity"`
	VerifiedWeight      float64 `json:"verified_weight"`
	VerifiedPricePerKgs float64 `json:"verified_price_per_kgs"`
	VerifiedSubtotal    int64   `json:"verified_subtotal"`
}

type SearchWasteDropRequestItemRequest struct {
	RequestID   string `json:"request_id"`
	WasteTypeID string `json:"waste_type_id"`
	Page        int    `json:"page,omitempty" validate:"min=1"`
	Size        int    `json:"size,omitempty" validate:"min=1,max=100"`
}

type GetWasteDropRequestItemRequest struct {
	ID string `json:"id"`
}

type UpdateWasteDropRequestItemRequest struct {
	ID                  string  `json:"id"`
	RequestID           string  `json:"request_id"`
	WasteTypeID         string  `json:"waste_type_id"`
	Quantity            int64   `json:"quantity"`
	VerifiedWeight      float64 `json:"verified_weight"`
	VerifiedPricePerKgs int64   `json:"verified_price_per_kgs"`
	VerifiedSubtotal    int64   `json:"verified_subtotal"`
}

type DeleteWasteDropRequestItemRequest struct {
	ID string `json:"id"`
}
