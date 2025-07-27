package model

type PointConversionSimpleResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Amount    int64  `json:"amount"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
}

type PointConversionResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Amount    int64  `json:"amount"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	IsDeleted bool   `json:"is_deleted"`
	User      *UserResponse
}

type PointConversionRequest struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
	Status string `json:"status"`
}

type SearchPointConversionRequest struct {
	UserID    string `json:"user_id"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	IsDeleted *bool  `json:"is_deleted"`
	OrderBy   string `json:"order_by"`
	OrderDir  string `json:"order_dir,omitempty"`
	Page      int    `json:"page,omitempty" validate:"min=1"`
	Size      int    `json:"size,omitempty" validate:"min=1,max=100"`
}
type GetPointConversionRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdatePointConversionRequest struct {
	ID        string `json:"id" validate:"required,max=100"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	IsDeleted *bool  `json:"is_deleted"`
}

type DeletePointConversionRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
