package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func WasteBankPricedTypeToSimpleResponse(wasteBankPricedType *entity.WasteBankPricedType) *model.WasteBankPricedTypeSimpleResponse {
	var wasteType *model.WasteTypeResponse
	if wasteBankPricedType.WasteTypeID != uuid.Nil {
		wasteType = WasteTypeToResponse(&wasteBankPricedType.WasteType)
	}
	return &model.WasteBankPricedTypeSimpleResponse{
		ID:                wasteBankPricedType.ID.String(),
		WasteBankID:       wasteBankPricedType.WasteBankID.String(),
		WasteTypeID:       wasteBankPricedType.WasteTypeID.String(),
		CustomPricePerKgs: wasteBankPricedType.CustomPricePerKgs,
		CreatedAt:         wasteBankPricedType.CreatedAt.String(),
		UpdatedAt:         wasteBankPricedType.UpdatedAt.String(),
		WasteType:         wasteType,
	}
}

func WasteBankPricedTypeToResponse(wasteBankPricedType *entity.WasteBankPricedType) *model.WasteBankPricedTypeResponse {
	var wasteBank *model.UserResponse
	if wasteBankPricedType.WasteBankID != uuid.Nil {
		wasteBank = UserToResponse(&wasteBankPricedType.WasteBank)
	}

	var wasteType *model.WasteTypeResponse
	if wasteBankPricedType.WasteTypeID != uuid.Nil {
		wasteType = WasteTypeToResponse(&wasteBankPricedType.WasteType)
	}
	return &model.WasteBankPricedTypeResponse{
		ID:                wasteBankPricedType.ID.String(),
		WasteBankID:       wasteBankPricedType.WasteBankID.String(),
		WasteTypeID:       wasteBankPricedType.WasteTypeID.String(),
		CustomPricePerKgs: wasteBankPricedType.CustomPricePerKgs,
		CreatedAt:         wasteBankPricedType.CreatedAt.String(),
		UpdatedAt:         wasteBankPricedType.UpdatedAt.String(),
		WasteBank:         wasteBank,
		WasteType:         wasteType,
	}
}
