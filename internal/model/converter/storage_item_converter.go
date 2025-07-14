package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func StorageItemToSimpleResponse(storageItem *entity.StorageItem) *model.StorageItemSimpleResponse {
	return &model.StorageItemSimpleResponse{
		ID:          storageItem.ID.String(),
		StorageID:   storageItem.StorageID.String(),
		WasteTypeID: storageItem.WasteTypeID.String(),
		WeightKgs:   storageItem.WeightKgs,
		CreatedAt:   storageItem.CreatedAt,
		UpdatedAt:   storageItem.UpdatedAt,
	}
}

func StorageItemToResponse(storageItem *entity.StorageItem) *model.StorageItemResponse {
	var storage *model.StorageSimpleResponse
	if storageItem.StorageID != uuid.Nil {
		storage = StorageToSimpleResponse(&storageItem.Storage)
	}
	var wasteType *model.WasteTypeResponse
	if storageItem.WasteTypeID != uuid.Nil {
		wasteType = WasteTypeToResponse(&storageItem.WasteType)
	}
	return &model.StorageItemResponse{
		ID:          storageItem.ID.String(),
		StorageID:   storageItem.StorageID.String(),
		WasteTypeID: storageItem.WasteTypeID.String(),
		WeightKgs:   storageItem.WeightKgs,
		CreatedAt:   storageItem.CreatedAt,
		UpdatedAt:   storageItem.UpdatedAt,
		Storage:     storage,
		WasteType:   wasteType,
	}
}
