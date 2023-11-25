package schedule

import (
	"github.com/patrickmn/go-cache"
	"log"
	"own/gin/rate/internal/load"
	"own/gin/rate/pkg/exchange"
	"own/gin/rate/pkg/exchange/openexchange"
	"own/gin/rate/pkg/quote/coingecko"
)

// SchedulingFeedPriceToCache will fetch price from quote & exchange, get exchange rate for denom to fiat
// TODO currently only support below config + hard coding
func SchedulingFeedPriceToCache(c *cache.Cache, config *load.Config) {

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
	USDPrice := (*cgQuoteResp)[coinIds]["usd"]

	// from usd get KRW, IDR, SGD, THB
	openExchangeResp, err := openexchange.FetchOpenExchangePrice(config.PriceProviderKey)
	if err != nil {
		log.Fatal(err)
		return
	}
	SGDPrice := openExchangeResp.Rates.SGD * USDPrice
	THBPrice := openExchangeResp.Rates.THB * USDPrice
	KRWPrice := openExchangeResp.Rates.KRW * USDPrice
	IDRPrice := openExchangeResp.Rates.IDR * USDPrice
	log.Default().Printf("[Prices] %s:USD 1:%f, %s:SGD 1:%f, %s:THB 1:%f, %s:KRW 1:%f, %s:IDR 1:%f",
		config.NodeServing, USDPrice, config.NodeServing, SGDPrice, config.NodeServing, THBPrice,
		config.NodeServing, KRWPrice, config.NodeServing, IDRPrice)

	// store into cache
	prices := exchange.QuotePrices{
		ToUSD: USDPrice,
		ToSGD: SGDPrice,
		ToTHB: THBPrice,
		ToKRW: KRWPrice,
		ToIDR: IDRPrice,
	}
	c.Set("CACHE_PRICES", prices, cache.DefaultExpiration)
}
