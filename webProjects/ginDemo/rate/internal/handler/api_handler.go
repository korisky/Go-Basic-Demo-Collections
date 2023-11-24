package handler

import (
	"github.com/gin-gonic/gin"
	"own/gin/rate/internal/load"
)

func SupplyPriceRequestHandler(c *gin.Context, config *load.Config) {
	denom := config.NodeServing

	// TODO 1. retrieve circulating supply (async later)

	// TODO 2. retrieve quote from cache (async later)

	// TODO 3. calculating market price & form response

	c.JSON(200, gin.H{
		"message": "Supply quote response" + denom,
	})
}

/*
	// TODO let an scheduler do below stuff (update per hour, currently)

	// TODO 1. choose the quote provider

	// TODO 2. calculate USD marketQuote

	// TODO 3. convert USD marketQuote -> KRW, IDR, SGD, THB quote

	// TODO 4. put it into local_cache
*/
