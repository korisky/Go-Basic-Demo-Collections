package handler

import (
	"github.com/gin-gonic/gin"
	"own/gin/rate/internal/load"
)

func SupplyPriceRequestHandler(c *gin.Context, config *load.Config) {
	serving := config.Serving
	c.JSON(200, gin.H{
		"message": "Supply price response" + serving,
	})
}
