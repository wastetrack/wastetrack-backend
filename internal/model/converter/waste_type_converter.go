package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func WasteTypeToResponse(wasteType *entity.WasteType) *model.WasteTypeResponse {
	var categoryResponse *model.WasteCategoryResponse
	if wasteType.CategoryID != uuid.Nil {
		categoryResponse = WasteCategoryToResponse(&wasteType.WasteCategory)
	}

	return &model.WasteTypeResponse{
		ID:            wasteType.ID.String(),
		CategoryID:    wasteType.CategoryID.String(),
		Name:          wasteType.Name,
		Description:   wasteType.Description,
		WasteCategory: categoryResponse,
	}
}
