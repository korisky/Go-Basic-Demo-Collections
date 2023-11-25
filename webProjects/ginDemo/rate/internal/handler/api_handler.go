package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"log"
	"own/gin/rate/internal/load"
	"own/gin/rate/pkg/exchange"
	"own/gin/rate/pkg/supply"
	"strconv"
	"strings"
	"time"
)

// SupplyPriceRequestHandler calculate market cap, from sync supply checking + cached quote price retrieving
func SupplyPriceRequestHandler(ctx *gin.Context, c *cache.Cache, config *load.Config) {

	// get denom supply
	circulatingSupply, failed := getCirculatingSupply(config)
	if failed {
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

// getCirculatingSupply will retrieve Denom circulating supply
func getCirculatingSupply(config *load.Config) (float64, bool) {
	// fetch supply
	supplyResp, err := supply.FetchSupply(config.NodeUrl)
	if err != nil {
		log.Fatal(err)
		return 0, true
	}

	// extract supply TODO compact into request util
	var circulatingSupply float64
	for _, item := range supplyResp.Supply {
		if load.FxServing == config.NodeServing && strings.EqualFold("fx", item.Denom) {
			circulatingSupply, err = strconv.ParseFloat(item.Amount, 64)
			break
		}
		if load.PundixServing == config.NodeServing && strings.EqualFold("ibc/55367B7B6572631B78A93C66EF9FDFCE87CDE372CC4ED7848DA78C1EB1DCDD78", item.Denom) {
			circulatingSupply, err = strconv.ParseFloat(item.Amount, 64)
			break
		}
	}
	if err != nil {
		log.Fatal(err)
		return 0, true
	}
	return circulatingSupply, false
}
