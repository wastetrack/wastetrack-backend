package model

type WasteSubCategoryResponse struct {
	ID            string                 `json:"id"`
	CategoryID    string                 `json:"category_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	WasteCategory *WasteCategoryResponse `json:"waste_category,omitempty"`
}

type WasteSubCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CategoryID  string `json:"category_id"`
}

type SearchWasteSubCategoryRequest struct {
	Name       string `json:"name,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	Page       int    `json:"page,omitempty" validate:"min=1"`
	Size       int    `json:"size,omitempty" validate:"min=1,max=100"`
}
type UpdateWasteSubCategoryRequest struct {
	ID          string `json:"id" validate:"required,max=100"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type DeleteWasteSubCategoryRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
