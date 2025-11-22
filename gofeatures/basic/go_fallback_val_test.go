package basic

import (
	"cmp"
	"fmt"
	"os"
	"testing"
)

// TestFallbackValCMPUsage 使用cmp进行环境参数获取&默认配置获取
func TestFallbackValCMPUsage(t *testing.T) {
	// cmp.Or() 逻辑是, 从前到后获取到第一个非nil的值, 则进行返回
	port := cmp.Or(getPortFromEnv(), getPortFromFlag())
	fmt.Println(port)
}

func getPortFromFlag() string {
	return "9098"
}

func getPortFromEnv() string {
	return os.Getenv("PORT")
}
