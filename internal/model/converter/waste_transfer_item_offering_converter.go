package converter

import (
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
	return &model.WasteTransferItemOfferingSimpleResponse{
		ID:                  item.ID.String(),
		TransferFormID:      item.TransferFormID.String(),
		WasteTypeID:         item.WasteTypeID.String(),
		OfferingWeight:      item.OfferingWeight,
		OfferingPricePerKgs: item.OfferingPricePerKgs,
		AcceptedWeight:      item.AcceptedWeight,
		AcceptedPricePerKgs: item.AcceptedPricePerKgs,
		VerifiedWeight:      item.VerifiedWeight,
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

	return &model.WasteTransferItemOfferingResponse{
		ID:                  item.ID.String(),
		TransferFormID:      item.TransferFormID.String(),
		WasteTypeID:         item.WasteTypeID.String(),
		OfferingWeight:      item.OfferingWeight,
		OfferingPricePerKgs: item.OfferingPricePerKgs,
		AcceptedWeight:      item.AcceptedWeight,
		AcceptedPricePerKgs: item.AcceptedPricePerKgs,
		VerifiedWeight:      item.VerifiedWeight,
		TransferForm:        transferForm,
		WasteType:           wasteType,
	}
}
