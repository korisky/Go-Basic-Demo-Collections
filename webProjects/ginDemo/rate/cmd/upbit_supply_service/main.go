package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"own/gin/rate/internal/handler"
)

func main() {
	port := 20950

	router := gin.Default()
	router.GET("/api/functionx/info", handler.SupplyPriceRequestHandler)
	err := router.Run(":" + string(rune(port)))

	if err != nil {
		log.Fatalf("Could not start the service on port:%d", port)
		return
	}
}
