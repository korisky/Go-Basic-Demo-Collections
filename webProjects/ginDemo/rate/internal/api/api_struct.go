package api

import (
	"strings"
	"time"
)

type ApiResponse []MarketCapItem

type MarketCapItem struct {
	Symbol              string  `json:"symbol"`
	CurrencyCode        string  `json:"currencyCode"`
	Price               float64 `json:"price"`
	MarketCap           float64 `json:"marketCap"`
	AccTradePrice24h    int64   `json:"accTradePrice24h"`
	CirculatingSupply   uint64  `json:"circulatingSupply"`
	MaxSupply           uint64  `json:"maxSupply"`
	Provider            string  `json:"provider"`
	LastUpdateTimestamp int64   `json:"lastUpdateTimestamp"`
}

// buildMarketCapItem build single MarketCapItem
func buildMarketCapItem(symbol, currency, provider string, price, supply float64, td time.Time) *MarketCapItem {
	return &MarketCapItem{
		Symbol:       strings.ToUpper(symbol),
		CurrencyCode: currency,
		Price:        price,
		MarketCap:    price * supply,
		//AccTradePrice24h:    0,
		CirculatingSupply:   uint64(supply),
		MaxSupply:           uint64(supply),
		Provider:            provider,
		LastUpdateTimestamp: td.UnixMilli(),
	}
}
