package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	"log"
	"own/gin/rate/internal/load"
	"own/gin/rate/pkg/fiatexchange"
	"own/gin/rate/pkg/supply"
	"sync"
)

var (
	// for singleton liked usage
	once       sync.Once
	supplyChan chan SupplyResult
	cacheChan  chan CacheResult
)

type SupplyResult struct {
	Supply decimal.Decimal
	Err    error
}

type CacheResult struct {
	Prices fiatexchange.ToFiatPrices
	Err    error
}

// SupplyPriceRequestHandler calculate market cap, from sync supply checking + cached exchange price retrieving
func SupplyPriceRequestHandler(ctx *gin.Context, c *cache.Cache, config *load.Config) {

	// asynchronous get supply and prices
	asyncGetSupplyAndPrice(config, c)
	supplyRes := <-supplyChan
	cacheRes := <-cacheChan
	if supplyRes.Err != nil || cacheRes.Err != nil {
		log.Fatalln()
	}

	// calculating market price & form response
	apiResponse := ApiResponse{}
	provider := "Function X"
	if load.PundixServing == config.NodeServing {
		provider = "Pundi X"
	}
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "USD", provider, cacheRes.Prices.ToUSD, supplyRes.Supply, cacheRes.Prices.UpdateTimestamp))
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "SGD", provider, cacheRes.Prices.ToSGD, supplyRes.Supply, cacheRes.Prices.UpdateTimestamp))
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "THB", provider, cacheRes.Prices.ToTHB, supplyRes.Supply, cacheRes.Prices.UpdateTimestamp))
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "KRW", provider, cacheRes.Prices.ToKRW, supplyRes.Supply, cacheRes.Prices.UpdateTimestamp))
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "IDR", provider, cacheRes.Prices.ToIDR, supplyRes.Supply, cacheRes.Prices.UpdateTimestamp))
	ctx.JSON(200, apiResponse)
}

// asyncGetSupplyAndPrice will request for supply & prices asynchronous
func asyncGetSupplyAndPrice(config *load.Config, c *cache.Cache) (chan SupplyResult, chan CacheResult) {
	// only can be executed once
	once.Do(func() {
		// channel init
		supplyChan = make(chan SupplyResult)
		cacheChan = make(chan CacheResult)
	})
	// get denom circulation supply
	go func() {
		circulatingSupply, err := supply.FetchTargetSupply(config)
		supplyChan <- SupplyResult{Supply: circulatingSupply, Err: err}
	}()
	// get cache
	go func() {
		data, exists := c.Get("CACHE_PRICES")
		if !exists {
			cacheChan <- CacheResult{Err: fmt.Errorf("no cache for price")}
			return
		}
		prices, _ := data.(fiatexchange.ToFiatPrices)
		cacheChan <- CacheResult{Prices: prices, Err: nil}
	}()
	// return the constructed channel
	return supplyChan, cacheChan
}
