package exchange

// QuotePriceFetcher is for converting external exchange prices to QuotePrices
type QuotePriceFetcher interface {
	FetchConvertToQuotePrices() (*QuotePrices, error)
}

// QuotePrices is the target prices would be cached
type QuotePrices struct {
	ToUSD           float64
	ToSGD           float64
	ToTHB           float64
	ToKRW           float64
	ToIDR           float64
	UpdateTimestamp int64
}
