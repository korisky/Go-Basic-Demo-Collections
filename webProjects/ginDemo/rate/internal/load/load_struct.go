package load

type ServeEnum string
type QuoteProviderEnum string
type PriceProviderEnum string

// ServeEnum for denom supply, FX or Pundix
const (
	FxServing     ServeEnum = "fx"
	PundixServing           = "pundix"
)

// QuoteProviderEnum for quote provider, CoinMarketCap, CoinGecko
const (
	CmcPrice QuoteProviderEnum = "coinmarketcap"
	CgPrice                    = "coingecko"
)

// PriceProviderEnum for price provider, ExchangeRate, OpenExchange
const (
	ExchangeRate PriceProviderEnum = "exchangerate"
	OpenExchange                   = "openexchange"
)

type Config struct {
	Port             uint64            `json:"port"`
	NodeServing      ServeEnum         `json:"node_serving"`
	NodeUrl          string            `json:"node_url"` // rest url
	QuoteProvider    QuoteProviderEnum `json:"quote_provider"`
	QuoteProviderUrl string            `json:"quote_provider_url"`
	QuoteProviderKey string            `json:"quote_provider_key"`
	PriceProvider    PriceProviderEnum `json:"price_provider"`
	PriceProviderKey string            `json:"price_provider_key"`
}
