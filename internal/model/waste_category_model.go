package model

type WasteCategoryResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type WasteCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type SearchWasteCategoryRequest struct {
	Name string `json:"name"`
	Page int    `json:"page,omitempty" validate:"min=1"`
	Size int    `json:"size,omitempty" validate:"min=1,max=100"`
}
type GetWasteCategoryRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateWasteCategoryRequest struct {
	ID          string `json:"id" validate:"required,max=100"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type DeleteWasteCategoryRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
