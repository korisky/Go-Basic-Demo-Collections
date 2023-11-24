package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"own/gin/rate/internal/handler"
	"own/gin/rate/internal/load"
	"strconv"
)

const API_PATH = "/api/%s/info"

func main() {
	// load configuration
	config, err := load.LoadConfiguration()
	if err != nil {
		log.Fatal(err)
		return

	}
	apiPath := fmt.Sprintf(API_PATH, config.Serving)

	// start service
	router := gin.Default()
	router.GET(apiPath,
		func(context *gin.Context) {
			handler.SupplyPriceRequestHandler(context, config)
		})
	err = router.Run(":" + strconv.FormatInt(config.Port, 10))
	if err != nil {
		log.Fatal(err)
	}
}
