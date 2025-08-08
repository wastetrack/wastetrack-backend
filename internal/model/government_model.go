package model

type GovernmentResponse struct {
	ID     string        `json:"id"`
	UserID string        `json:"user_id"`
	User   *UserResponse `json:"user,omitempty"`
}

type GovernmentRequest struct {
	UserID string `json:"user_id"`
}

type GetGovernmentRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

type UpdateGovernmentRequest struct {
	ID string `json:"id" validate:"required,max=100"`
	// Add fields as needed when government profile gets more properties
}

type DeleteGovernmentRequest struct {
	ID string `json:"id" validate:"required,max=100"`
}

// GovernmentDashboardRequest represents the request for government dashboard data
type GovernmentDashboardRequest struct {
	StartMonth string `json:"start_month" ` // Format: "2024-01"
	EndMonth   string `json:"end_month"`    // Format: "2024-12"
	Province   string `json:"province,omitempty"`
	City       string `json:"city,omitempty"`
}

// GovernmentDashboardResponse represents the complete dashboard response
type GovernmentDashboardResponse struct {
	TotalBankSampah  int64                   `json:"totalBankSampah"`
	TotalOfftaker    int64                   `json:"totalOfftaker"`
	TotalCollected   float64                 `json:"totalCollected"`
	CollectionTrends []CollectionTrendByRole `json:"collectionTrends"`
	TopOfftakers     []TopOfftaker           `json:"topOfftakers"`
	LargestBanks     []LargestBank           `json:"largestBanks"`
}

// CollectionTrend represents collection data over time
type CollectionTrendByRole struct {
	Month              string  `json:"month"`
	WasteBankUnit      float64 `json:"waste_bank_unit"`
	WasteBankCentral   float64 `json:"waste_bank_central"`
	Industry           float64 `json:"industry"`
	CollectionRequests float64 `json:"collection_requests_weight"`
	TransferRequests   float64 `json:"transfer_requests_weight"`
	TotalAmount        float64 `json:"total_requests_weight"`
}

// TopOfftaker represents top performing offtakers
type TopOfftaker struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Institution string  `json:"institution"`
	City        string  `json:"city"`
	Province    string  `json:"province"`
	TotalWeight float64 `json:"total_weight"`
	TotalPrice  int64   `json:"total_price"`
}

// LargestBank represents largest waste banks
type LargestBank struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Institution string  `json:"institution"`
	City        string  `json:"city"`
	Province    string  `json:"province"`
	TotalWeight float64 `json:"total_weight"` // Total waste processed (waste drops + transfers)
	Volume      float64 `json:"volume"`       // Storage volume (length * width * height)
}

// WasteBankCountRequest for counting waste banks in a time period
type WasteBankCountRequest struct {
	StartMonth string `json:"start_month"`
	EndMonth   string `json:"end_month"`
	Province   string `json:"province,omitempty"`
	City       string `json:"city,omitempty"`
}
