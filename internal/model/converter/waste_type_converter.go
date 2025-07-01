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

	var subcategoryResponse *model.WasteSubCategorySimpleResponse
	if wasteType.SubcategoryID != uuid.Nil {
		subcategoryResponse = &model.WasteSubCategorySimpleResponse{
			ID:          wasteType.WasteSubcategory.ID.String(),
			CategoryID:  wasteType.WasteSubcategory.CategoryID.String(),
			Name:        wasteType.WasteSubcategory.Name,
			Description: wasteType.WasteSubcategory.Description,
		}
	}

	return &model.WasteTypeResponse{
		ID:               wasteType.ID.String(),
		CategoryID:       wasteType.CategoryID.String(),
		SubcategoryID:    wasteType.SubcategoryID.String(),
		Name:             wasteType.Name,
		Description:      wasteType.Description,
		WasteCategory:    categoryResponse,
		WasteSubCategory: subcategoryResponse,
	}
}
