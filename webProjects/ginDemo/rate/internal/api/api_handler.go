package api

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"log"
	"own/gin/rate/internal/load"
	"own/gin/rate/pkg/exchange"
	"own/gin/rate/pkg/supply"
	"time"
)

// SupplyPriceRequestHandler calculate market cap, from sync supply checking + cached quote price retrieving
func SupplyPriceRequestHandler(ctx *gin.Context, c *cache.Cache, config *load.Config) {

	// get denom circulation supply
	circulatingSupply, err := supply.FetchCirculatingSupply(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	// get cache
	data, exists := c.Get("CACHE_PRICES")
	if !exists {
		log.Fatal("No Cache")
		return
	}
	prices, _ := data.(exchange.QuotePrices)

	// calculating market price & form response
	apiResponse := ApiResponse{}
	provider := "Function X"
	if load.PundixServing == config.NodeServing {
		provider = "Pundi X"
	}
	now := time.Now()
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "USD", provider, prices.ToUSD, circulatingSupply, now))
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "SGD", provider, prices.ToSGD, circulatingSupply, now))
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "THB", provider, prices.ToTHB, circulatingSupply, now))
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "KRW", provider, prices.ToKRW, circulatingSupply, now))
	apiResponse = append(apiResponse, *buildMarketCapItem(string(config.NodeServing), "IDR", provider, prices.ToIDR, circulatingSupply, now))
	ctx.JSON(200, apiResponse)
}
