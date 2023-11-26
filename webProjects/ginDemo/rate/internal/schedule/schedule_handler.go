package schedule

import (
	"github.com/patrickmn/go-cache"
	"log"
	"own/gin/rate/internal/load"
	"own/gin/rate/pkg/exchange/openexchange"
	"own/gin/rate/pkg/quote/coingecko"
)

// SchedulingFeedPriceToCache will fetch price from quote & exchange, get exchange rate for denom to fiat
func SchedulingFeedPriceToCache(c *cache.Cache, config *load.Config) {

	// TODO currently only support below config + hard coding
	quoteProvider := config.QuoteProvider
	priceProvider := config.PriceProvider
	if quoteProvider != load.CgPrice || priceProvider != load.OpenExchange {
		log.Fatal("Currently only support CoinGecko as quote provider, OpenExchange as exchange price provider")
		return
	}

	// fetch quote price
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
	fetcher := openexchange.OxFetcher{
		UsdPrice: UsdPrice,
		ApiKey:   config.PriceProviderKey}
	prices, err := fetcher.FetchConvertToQuotePrices()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Default().Printf("[Prices] %s:USD 1:%f, %s:SGD 1:%f, %s:THB 1:%f, %s:KRW 1:%f, %s:IDR 1:%f",
		config.NodeServing, prices.ToUSD, config.NodeServing, prices.ToSGD, config.NodeServing, prices.ToTHB,
		config.NodeServing, prices.ToKRW, config.NodeServing, prices.ToIDR)
	c.Set("CACHE_PRICES", *prices, cache.DefaultExpiration)
}
