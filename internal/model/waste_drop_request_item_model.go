package model

type WasteDropRequestItemResponse struct {
	ID               string             `json:"id"`
	WasteType        *WasteTypeResponse `json:"waste_type"`
	Quantity         int64              `json:"quantity"`
	VerifiedWeight   float64            `json:"verified_weight"`
	VerifiedSubtotal int64              `json:"verified_subtotal"`
}
