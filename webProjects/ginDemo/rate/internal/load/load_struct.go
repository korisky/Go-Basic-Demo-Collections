package load

type ServeEnum string
type PriceProviderEnum string

// ServeEnum for denom supply, FX or Pundix
const (
	FxServing     ServeEnum = "fx"
	PundixServing           = "pundix"
)

// PriceProviderEnum for quote provider, CoinMarketCap, CoinGecko
const (
	CmcPrice PriceProviderEnum = "coinmarketcap"
	CgPrice                    = "coingecko"
)

type Config struct {
	Port             uint64            `json:"port"`
	NodeServing      ServeEnum         `json:"node_serving"`
	NodeUrl          string            `json:"node_url"` // rest url
	PriceProvider    PriceProviderEnum `json:"price_provider"`
	PriceProviderUrl string            `json:"price_provider_url"`
	PriceProviderKey string            `json:"price_provider_key"`
}
