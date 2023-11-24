package handler

import "github.com/gin-gonic/gin"

func SupplyPriceRequestHandler(c *gin.Context) {

	// TODO 1. retrieve

	c.JSON(200, gin.H{
		"message": "Supply price response",
	})
}
