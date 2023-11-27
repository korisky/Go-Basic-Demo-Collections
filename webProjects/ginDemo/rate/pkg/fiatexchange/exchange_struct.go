package fiatexchange

// ToFiatPricesFetcher is for converting external fiatexchange prices to ToFiatPrices
type ToFiatPricesFetcher interface {
	FetchToAllFiatPrices() (*ToFiatPrices, error)
}

// ToFiatPrices is the target prices would be cached
type ToFiatPrices struct {
	ToUSD           float64
	ToSGD           float64
	ToTHB           float64
	ToKRW           float64
	ToIDR           float64
	UpdateTimestamp int64
}
