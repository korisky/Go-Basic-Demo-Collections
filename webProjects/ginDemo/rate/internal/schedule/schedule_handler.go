package schedule

import (
	"github.com/patrickmn/go-cache"
	"log"
	"own/gin/rate/internal/load"
	"own/gin/rate/pkg/cryptoexchange"
	"own/gin/rate/pkg/cryptoexchange/coingecko"
	"own/gin/rate/pkg/cryptoexchange/coinmarketcap"
	"own/gin/rate/pkg/fiatexchange"
	"own/gin/rate/pkg/fiatexchange/exchangerate"
	"own/gin/rate/pkg/fiatexchange/openexchange"
)

// SchedulingFeedPriceToCache will fetch price from quote-provider & price-provider, get exchange rate for denom to fiat
func SchedulingFeedPriceToCache(c *cache.Cache, config *load.Config) {

	// fetch exchange price
	cgUseId, cmcUseId := determineCoinId(config)
	var quoteFetcher cryptoexchange.ToUsdPriceFetcher
	switch config.QuoteProvider {
	case load.CgPrice:
		quoteFetcher = &coingecko.CgFetcher{Ids: cgUseId, Currencies: "usd"}
	case load.CmcPrice:
		quoteFetcher = &coinmarketcap.CmcFetcher{Id: cmcUseId, ConvertId: "2781", ApiKey: config.QuoteProviderKey}
	default:
		log.Fatalf("Not supported quote provider: %v", config.QuoteProvider)
	}
	usdPrice, err := quoteFetcher.FetchToUsdPrice()
	if err != nil {
		log.Fatal(err)
		return
	}

	// fetch fiat prices
	var fetcher fiatexchange.ToFiatPricesFetcher
	switch config.PriceProvider {
	case load.OpenExchange:
		fetcher = &openexchange.OxFetcher{UsdPrice: usdPrice.DenomToUsd, ApiKey: config.PriceProviderKey}
	case load.ExchangeRate:
		fetcher = &exchangerate.ErFetcher{UsdPrice: usdPrice.DenomToUsd, ApiKey: config.PriceProviderKey}
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

	// save in cache
	c.Set("CACHE_PRICES", *prices, cache.DefaultExpiration)
}

// determineCoinId is for giving correct coin-id according to the node serving, (cg, cmc)
func determineCoinId(config *load.Config) (string, string) {
	if config.NodeServing == load.PundixServing {
		return "pundi-x-2", "9004"
	}
	return "fx-coin", "3884"
}
