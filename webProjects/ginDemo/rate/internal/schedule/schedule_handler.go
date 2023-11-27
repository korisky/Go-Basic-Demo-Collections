package schedule

import (
	"github.com/patrickmn/go-cache"
	"log"
	"own/gin/rate/internal/load"
	"own/gin/rate/pkg/cryptoexchange/coingecko"
	"own/gin/rate/pkg/fiatexchange"
	"own/gin/rate/pkg/fiatexchange/exchangerate"
	"own/gin/rate/pkg/fiatexchange/openexchange"
)

// SchedulingFeedPriceToCache will fetch price from cryptoexchange & fiatexchange, get fiatexchange rate for denom to fiat
func SchedulingFeedPriceToCache(c *cache.Cache, config *load.Config) {

	// TODO currently only support CoinGecko for cryptoexchange price provider
	//quoteProvider := config.QuoteProvider

	// fetch cryptoexchange price
	coinIds := "fx-coin"
	if config.NodeServing == load.PundixServing {
		coinIds = "pundi-x"
	}
	cgQuoteResp, err := coingecko.FetchCgQuotePrice(coinIds, "usd")
	if err != nil {
		log.Fatal(err)
		return
	}
	UsdPrice := (*cgQuoteResp)[coinIds]["usd"]

	// fetch fiat prices
	var fetcher fiatexchange.ToFiatPricesFetcher
	switch config.PriceProvider {
	case load.OpenExchange:
		fetcher = &openexchange.OxFetcher{UsdPrice: UsdPrice, ApiKey: config.PriceProviderKey}
	case load.ExchangeRate:
		fetcher = &exchangerate.ErFetcher{UsdPrice: UsdPrice, ApiKey: config.PriceProviderKey}
	default:
		log.Fatalf("Not supported price provider: %v", config.PriceProvider)
	}
	prices, err := fetcher.FetchToAllFiatPrices()
	if err != nil {
		log.Fatalf("Error on fetching fiatexchange prices: %v", err)
	}
	log.Default().Printf("[Prices] %s:USD 1:%f, %s:SGD 1:%f, %s:THB 1:%f, %s:KRW 1:%f, %s:IDR 1:%f",
		config.NodeServing, prices.ToUSD, config.NodeServing, prices.ToSGD, config.NodeServing, prices.ToTHB,
		config.NodeServing, prices.ToKRW, config.NodeServing, prices.ToIDR)
	c.Set("CACHE_PRICES", *prices, cache.DefaultExpiration)
}
