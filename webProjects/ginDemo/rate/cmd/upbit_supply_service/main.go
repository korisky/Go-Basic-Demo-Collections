package main

import (
	"github.com/gin-gonic/gin"
	"own/gin/rate/internal"
)

func main() {
	router := gin.Default()
	router.GET("/api/functionx/info", internal.SupplyPriceRequestHandler)
	err := router.Run(":20950")
	if err != nil {
		return
	}
}
