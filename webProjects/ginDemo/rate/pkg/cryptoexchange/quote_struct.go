package cryptoexchange

// ToUsdPriceFetcher interface for fetching denom to usd price
type ToUsdPriceFetcher interface {
	FetchToUsdPrice() (*ToUsdPrice, error)
}

type ToUsdPrice struct {
	DenomToUsd float64
}
