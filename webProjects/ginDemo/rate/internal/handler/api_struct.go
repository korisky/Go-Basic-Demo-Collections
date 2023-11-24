package handler

import "time"

type ApiResponse []MarketCapItem

type MarketCapItem struct {
	Symbol              string    `json:"symbol"`
	CurrencyCode        string    `json:"currencyCode"`
	Price               float64   `json:"price"`
	MarketCap           float64   `json:"marketCap"`
	AccTradePrice24h    int64     `json:"accTradePrice24h"`
	CirculatingSupply   uint64    `json:"circulatingSupply"`
	MaxSupply           uint64    `json:"maxSupply"`
	Provider            string    `json:"provider"`
	LastUpdateTimestamp time.Time `json:"lastUpdateTimestamp"`
}
