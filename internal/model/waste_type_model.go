package model

type WasteTypeResponse struct {
	ID            string                 `json:"id"`
	CategoryID    string                 `json:"category_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	WasteCategory *WasteCategoryResponse `json:"waste_category,omitempty"`
}

type WasteTypeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CategoryID  string `json:"category_id"`
}

type SearchWasteTypeRequest struct {
	Name       string `json:"name,omitempty"`
	CategoryID string `json:"category_id,omitempty"`
	Page       int    `json:"page,omitempty" validate:"min=1"`
	Size       int    `json:"size,omitempty" validate:"min=1,max=100"`
}
type UpdateWasteTypeRequest struct {
	ID          string `json:"id" validate:"required,max=100"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type DeleteWasteTypeRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
