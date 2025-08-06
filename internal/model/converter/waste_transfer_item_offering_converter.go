package converter

import (
	"math"

	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

// WasteTransferItemOffering converters

func WasteTransferItemOfferingToSimpleResponse(item *entity.WasteTransferItemOffering) *model.WasteTransferItemOfferingSimpleResponse {
	var wasteType *model.WasteTypeResponse
	if item.WasteTypeID != uuid.Nil {
		wasteType = WasteTypeToResponse(&item.WasteType)
	}

	// Calculate loss weight: if verified_weight is 0, return 0; otherwise return absolute difference
	var lossWeight float64
	if item.VerifiedWeight == 0 {
		lossWeight = 0
	} else {
		lossWeight = math.Abs(item.AcceptedWeight - item.VerifiedWeight)
	}

	return &model.WasteTransferItemOfferingSimpleResponse{
		ID:                  item.ID.String(),
		TransferFormID:      item.TransferFormID.String(),
		WasteTypeID:         item.WasteTypeID.String(),
		OfferingWeight:      item.OfferingWeight,
		OfferingPricePerKgs: item.OfferingPricePerKgs,
		AcceptedWeight:      item.AcceptedWeight,
		AcceptedPricePerKgs: item.AcceptedPricePerKgs,
		VerifiedWeight:      item.VerifiedWeight,
		LossWeight:          lossWeight,
		WasteType:           wasteType,
	}
}

func WasteTransferItemOfferingToResponse(item *entity.WasteTransferItemOffering) *model.WasteTransferItemOfferingResponse {
	var transferForm *model.WasteTransferRequestSimpleResponse
	if item.TransferFormID != uuid.Nil {
		transferForm = WasteTransferRequestToSimpleResponse(&item.TransferForm)
	}

	var wasteType *model.WasteTypeResponse
	if item.WasteTypeID != uuid.Nil {
		wasteType = WasteTypeToResponse(&item.WasteType)
	}

	// Calculate loss weight: if verified_weight is 0, return 0; otherwise return absolute difference
	var lossWeight float64
	if item.VerifiedWeight == 0 {
		lossWeight = 0
	} else {
		lossWeight = math.Abs(item.AcceptedWeight - item.VerifiedWeight)
	}

	return &model.WasteTransferItemOfferingResponse{
		ID:                  item.ID.String(),
		TransferFormID:      item.TransferFormID.String(),
		WasteTypeID:         item.WasteTypeID.String(),
		OfferingWeight:      item.OfferingWeight,
		OfferingPricePerKgs: item.OfferingPricePerKgs,
		AcceptedWeight:      item.AcceptedWeight,
		AcceptedPricePerKgs: item.AcceptedPricePerKgs,
		VerifiedWeight:      item.VerifiedWeight,
		LossWeight:          lossWeight,
		TransferForm:        transferForm,
		WasteType:           wasteType,
	}
}
