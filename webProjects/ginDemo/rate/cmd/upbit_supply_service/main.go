package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"own/gin/rate/internal/handler"
	"own/gin/rate/internal/load"
	"strconv"
)

const ApiPath = "/api/%s/info"

func main() {
	// load configuration
	config, err := load.LoadConfiguration()
	if err != nil {
		log.Fatal(err)
		return

	}
	apiPath := fmt.Sprintf(ApiPath, config.NodeServing)

	// start service
	router := gin.Default()
	router.GET(apiPath,
		func(context *gin.Context) {
			handler.SupplyPriceRequestHandler(context, config)
		})
	err = router.Run(":" + strconv.FormatUint(config.Port, 10))
	if err != nil {
		log.Fatal(err)
	}
}
