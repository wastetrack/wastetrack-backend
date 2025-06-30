package converter

import (
	"github.com/google/uuid"
	"github.com/wastetrack/wastetrack-backend/internal/entity"
	"github.com/wastetrack/wastetrack-backend/internal/model"
)

func WasteSubCategoryToResponse(wasteSubCategory *entity.WasteSubcategory) *model.WasteSubCategoryResponse {
	var categoryResponse *model.WasteCategoryResponse
	if wasteSubCategory.CategoryID != uuid.Nil {
		categoryResponse = WasteCategoryToResponse(&wasteSubCategory.WasteCategory)
	}
	return &model.WasteSubCategoryResponse{
		ID:            wasteSubCategory.ID.String(),
		CategoryID:    wasteSubCategory.CategoryID.String(),
		Name:          wasteSubCategory.Name,
		Description:   wasteSubCategory.Description,
		WasteCategory: categoryResponse,
	}
}
