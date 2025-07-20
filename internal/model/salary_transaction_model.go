package model

type SalaryTransactionSimpleResponse struct {
	ID              string `json:"id"`
	SenderID        string `json:"sender_id"`
	ReceiverID      string `json:"receiver_id"`
	TransactionType string `json:"transaction_type"`
	Amount          int64  `json:"amount"`
	IsDeleted       bool   `json:"is_deleted"`
	CreatedAt       string `json:"created_at"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
}

type SalaryTransactionResponse struct {
	ID              string `json:"id"`
	SenderID        string `json:"sender_id"`
	ReceiverID      string `json:"receiver_id"`
	TransactionType string `json:"transaction_type"`
	Amount          int64  `json:"amount"`
	CreatedAt       string `json:"created_at"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
	IsDeleted       bool   `json:"is_deleted"`
	Sender          *UserResponse
	Receiver        *UserResponse
}
type SalaryTransactionRequest struct {
	SenderID        string `json:"-"`
	ReceiverID      string `json:"receiver_id"`
	TransactionType string `json:"transaction_type"`
	Amount          int64  `json:"amount"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
}

type SearchSalaryTransactionRequest struct {
	SenderID        string `json:"sender_id"`
	ReceiverID      string `json:"receiver_id"`
	TransactionType string `json:"transaction_type"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
	IsDeleted       *bool  `json:"is_deleted"`
	Page            int    `json:"page,omitempty" validate:"min=1"`
	Size            int    `json:"size,omitempty" validate:"min=1,max=100"`
}
type GetSalaryTransactionRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
type UpdateSalaryTransactionRequest struct {
	ID              string `json:"id" validate:"required,max=100"`
	TransactionType string `json:"transaction_type"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
}

type DeleteSalaryTransactionRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}
