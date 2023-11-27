package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	"log"
	"own/gin/rate/internal/api"
	"own/gin/rate/internal/load"
	"own/gin/rate/internal/schedule"
	"strconv"
	"time"
)

// ApiPath is the expose path, e.g. /api/fx/info, /api/pundix/info
const ApiPath = "/api/%s/info"

func main() {
	// 1. load configuration
	config, apiPath, failed := loadConfigProcess()
	if failed {
		return
	}
	// 2. init cache & run next
	c, cr := initCacheProcess(config)
	// 3. init scheduler
	failed = scheduleSetting(c, cr, config)
	if failed {
		return
	}
	// 4. start web service
	startWebServiceProcess(c, apiPath, config)
}

// loadConfigProcess for loading the configs from config.json
func loadConfigProcess() (*load.Config, string, bool) {
	config, err := load.LoadConfiguration("config/config.json")
	if err != nil {
		log.Fatal(err)
		return nil, "", true
	}
	apiPath := fmt.Sprintf(ApiPath, config.NodeServing)
	return config, apiPath, false
}

// initCacheProcess for init the cache & scheduler
func initCacheProcess(config *load.Config) (*cache.Cache, *cron.Cron) {
	c := cache.New(65*time.Minute, 75*time.Minute)
	cr := cron.New()
	// run the fetching task, async
	go schedule.SchedulingFeedPriceToCache(c, config)
	return c, cr
}

// scheduleSetting for init the scheduler & task
func scheduleSetting(c *cache.Cache, cr *cron.Cron, config *load.Config) bool {
	_, err := cr.AddFunc("@every 1h", func() {
		schedule.SchedulingFeedPriceToCache(c, config)
	})
	if err != nil {
		log.Fatal(err)
		return true
	}
	cr.Start()
	return false
}

// startWebServiceProcess expose the API and serve
func startWebServiceProcess(c *cache.Cache, apiPath string, config *load.Config) {
	router := gin.Default()
	router.GET(apiPath,
		func(context *gin.Context) {
			api.SupplyPriceRequestHandler(context, c, config)
		})
	err := router.Run(":" + strconv.FormatUint(config.Port, 10))
	if err != nil {
		log.Fatal(err)
	}
}
