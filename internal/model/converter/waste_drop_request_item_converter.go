package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func WasteDropRequestItemToSimpleResponse(wasteDropRequestItem *entity.WasteDropRequestItem) *model.WasteDropRequestItemSimpleResponse {
	return &model.WasteDropRequestItemSimpleResponse{
		ID:                  wasteDropRequestItem.ID.String(),
		RequestID:           wasteDropRequestItem.RequestID.String(),
		WasteTypeID:         wasteDropRequestItem.WasteTypeID.String(),
		Quantity:            wasteDropRequestItem.Quantity,
		VerifiedWeight:      wasteDropRequestItem.VerifiedWeight,
		VerifiedPricePerKgs: wasteDropRequestItem.VerifiedPricePerKgs,
		VerifiedSubtotal:    wasteDropRequestItem.VerifiedSubtotal,
	}
}

func WasteDropRequestItemToResponse(wasteDropRequestItem *entity.WasteDropRequestItem) *model.WasteDropRequestItemResponse {
	var wasteDropRequest *model.WasteDropRequestSimpleResponse
	if wasteDropRequestItem.RequestID != uuid.Nil {
		wasteDropRequest = WasteDropRequestToSimpleResponse(&wasteDropRequestItem.Request)
	}
	var wasteType *model.WasteTypeResponse
	if wasteDropRequestItem.WasteTypeID != uuid.Nil {
		wasteType = WasteTypeToResponse(&wasteDropRequestItem.WasteType)
	}
	return &model.WasteDropRequestItemResponse{
		ID:                  wasteDropRequestItem.ID.String(),
		RequestID:           wasteDropRequestItem.RequestID.String(),
		WasteTypeID:         wasteDropRequestItem.WasteTypeID.String(),
		Quantity:            wasteDropRequestItem.Quantity,
		VerifiedWeight:      wasteDropRequestItem.VerifiedWeight,
		VerifiedPricePerKgs: wasteDropRequestItem.VerifiedPricePerKgs,
		VerifiedSubtotal:    wasteDropRequestItem.VerifiedSubtotal,
		Request:             wasteDropRequest,
		WasteType:           wasteType,
	}
}
