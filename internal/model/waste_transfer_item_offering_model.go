package model

// Waste Transfer Item Offering models
type WasteTransferItemOfferingSimpleResponse struct {
	ID                  string  `json:"id"`
	TransferFormID      string  `json:"transfer_form_id"`
	WasteTypeID         string  `json:"waste_type_id"`
	OfferingWeight      float64 `json:"offering_weight"`
	OfferingPricePerKgs int64   `json:"offering_price_per_kgs"`
	AcceptedWeight      float64 `json:"accepted_weight"`
	AcceptedPricePerKgs int64   `json:"accepted_price_per_kgs"`
	VerifiedWeight      float64 `json:"verified_weight"`
}

type WasteTransferItemOfferingResponse struct {
	ID                  string                              `json:"id"`
	TransferFormID      string                              `json:"transfer_form_id"`
	WasteTypeID         string                              `json:"waste_type_id"`
	OfferingWeight      float64                             `json:"offering_weight"`
	OfferingPricePerKgs int64                               `json:"offering_price_per_kgs"`
	AcceptedWeight      float64                             `json:"accepted_weight"`
	AcceptedPricePerKgs int64                               `json:"accepted_price_per_kgs"`
	VerifiedWeight      float64                             `json:"verified_weight"`
	TransferForm        *WasteTransferRequestSimpleResponse `json:"transfer_form,omitempty"`
	WasteType           *WasteTypeResponse                  `json:"waste_type,omitempty"`
}

type WasteTransferItemOfferingRequest struct {
	TransferFormID      string  `json:"transfer_form_id" validate:"required"`
	WasteTypeID         string  `json:"waste_type_id" validate:"required"`
	OfferingWeight      float64 `json:"offering_weight" validate:"required,min=0"`
	OfferingPricePerKgs int64   `json:"offering_price_per_kgs" validate:"required,min=0"`
	AcceptedWeight      float64 `json:"accepted_weight"`
	AcceptedPricePerKgs int64   `json:"accepted_price_per_kgs"`
}

type SearchWasteTransferItemOfferingRequest struct {
	TransferFormID string `json:"transfer_form_id"`
	WasteTypeID    string `json:"waste_type_id"`
	Page           int    `json:"page,omitempty" validate:"min=1"`
	Size           int    `json:"size,omitempty" validate:"min=1,max=100"`
}

type GetWasteTransferItemOfferingRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type UpdateWasteTransferItemOfferingRequest struct {
	ID                  string  `json:"id" validate:"required,max=100"`
	OfferingWeight      float64 `json:"offering_weight"`
	OfferingPricePerKgs int64   `json:"offering_price_per_kgs"`
	AcceptedWeight      float64 `json:"accepted_weight"`
	AcceptedPricePerKgs int64   `json:"accepted_price_per_kgs"`
}

type DeleteWasteTransferItemOfferingRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
