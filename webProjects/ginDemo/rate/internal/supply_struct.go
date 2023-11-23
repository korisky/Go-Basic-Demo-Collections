package internal

// SupplyItem circulating supply for each denom
type SupplyItem struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

// Pagination of supply API
type Pagination struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}

// SupplyApiResponse api response
type SupplyApiResponse struct {
	Supply     []SupplyItem `json:"supply"`
	Pagination Pagination   `json:"pagination"`
}
