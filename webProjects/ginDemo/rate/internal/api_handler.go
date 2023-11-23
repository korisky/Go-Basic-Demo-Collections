package internal

import "github.com/gin-gonic/gin"

func SupplyPriceRequestHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Supply price response",
	})
}
