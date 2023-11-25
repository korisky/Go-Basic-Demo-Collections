package schedule

import (
	"github.com/patrickmn/go-cache"
	"own/gin/rate/internal/load"
	"testing"
	"time"
)

func Test_ScheduleTask(t *testing.T) {
	config, _ := load.LoadConfiguration("../../config/config.json")
	c := cache.New(65*time.Minute, 75*time.Minute)
	SchedulingFeedPriceToCache(c, config)
}
