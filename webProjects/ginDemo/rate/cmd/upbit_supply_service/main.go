package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"log"
	"own/gin/rate/internal/handler"
	"own/gin/rate/internal/load"
	"strconv"
	"time"
)

// ApiPath is the expose path, e.g. /api/fx/info, /api/pundix/info
const ApiPath = "/api/%s/info"

func main() {
	// 1. load configuration
	config, err, apiPath, failed := loadConfigProcess()
	if failed {
		return
	}

	// 2. init cache & run next
	c, cr := initCacheProcess(config)

	// 3. init scheduler
	err, failed = scheduleSetting(err, cr, c, config)
	if failed {
		return
	}

	// 4. start web service
	startWebServiceProcess(apiPath, config, err)
}

// loadConfigProcess for loading the configs from config.json
func loadConfigProcess() (*load.Config, error, string, bool) {
	config, err := load.LoadConfiguration()
	if err != nil {
		log.Fatal(err)
		return nil, nil, "", true
	}
	apiPath := fmt.Sprintf(ApiPath, config.NodeServing)
	return config, err, apiPath, false
}

// initCacheProcess for init the cache & scheduler
func initCacheProcess(config *load.Config) (*cache.Cache, *cron.Cron) {
	c := cache.New(65*time.Minute, 75*time.Minute)
	cr := cron.New()
	// run the fetching task, async
	go handler.SchedulingTask(c, config)
	return c, cr
}

// scheduleSetting for init the scheduler & task
func scheduleSetting(err error, cr *cron.Cron, c *cache.Cache, config *load.Config) (error, bool) {
	_, err = cr.AddFunc("@every 1h", func() {
		handler.SchedulingTask(c, config)
	})
	if err != nil {
		log.Fatal(err)
		return nil, true
	}
	cr.Start()
	return err, false
}

// startWebServiceProcess expose the API and serve
func startWebServiceProcess(apiPath string, config *load.Config, err error) {
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
