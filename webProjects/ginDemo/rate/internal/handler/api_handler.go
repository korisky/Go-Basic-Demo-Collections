package handler

import (
	"github.com/gin-gonic/gin"
	"own/gin/rate/internal/load"
)

func SupplyPriceRequestHandler(c *gin.Context, config *load.Config) {
	denom := config.NodeServing

	// TODO 1. choose the denom supply

	// TODO 2. choose the quote provider

	// TODO 3. calculate USD marketCap

	// TODO 4. convert USD marketCap -> KRW, IDR, SGD, THB quote

	c.JSON(200, gin.H{
		"message": "Supply quote response" + denom,
	})
}
