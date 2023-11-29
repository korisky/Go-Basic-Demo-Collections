package api

import (
	"github.com/shopspring/decimal"
	"strings"
)

type ApiResponse []MarketCapItem

type MarketCapItem struct {
	Symbol              string          `json:"symbol"`
	CurrencyCode        string          `json:"currencyCode"`
	Price               float64         `json:"price"`
	MarketCap           decimal.Decimal `json:"marketCap"`
	AccTradePrice24h    int64           `json:"accTradePrice24h"`
	CirculatingSupply   decimal.Decimal `json:"circulatingSupply"`
	MaxSupply           decimal.Decimal `json:"maxSupply"`
	Provider            string          `json:"provider"`
	LastUpdateTimestamp int64           `json:"lastUpdateTimestamp"`
}

// buildMarketCapItem build single MarketCapItem
func buildMarketCapItem(symbol, currency, provider string, price float64, sup decimal.Decimal, td int64) *MarketCapItem {
	return &MarketCapItem{
		Symbol:              strings.ToUpper(symbol),
		CurrencyCode:        currency,
		Price:               price,
		MarketCap:           decimal.NewFromFloat(price).Mul(sup),
		AccTradePrice24h:    0,
		CirculatingSupply:   sup,
		MaxSupply:           sup,
		Provider:            provider,
		LastUpdateTimestamp: td,
	}
}
