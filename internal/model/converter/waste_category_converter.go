package converter

import (
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func WasteCategoryToResponse(wasteCategory *entity.WasteCategory) *model.WasteCategoryResponse {
	return &model.WasteCategoryResponse{
		ID:          wasteCategory.ID.String(),
		Name:        wasteCategory.Name,
		Description: wasteCategory.Description,
	}
}
