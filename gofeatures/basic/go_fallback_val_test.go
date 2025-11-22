package basic

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

// TestFallbackValCMPUsage 使用cmp进行环境参数获取&默认配置获取
func TestFallbackValCMPUsage(t *testing.T) {
	// cmp.Or() 逻辑是, 从前到后获取到第一个非nil的值, 则进行返回
	port := cmp.Or(getPortFromEnv(), getPortFromFlag(), "8080")
	fmt.Println(port)
}

func getPortFromFlag() string {
	return "9098"
}

func getPortFromEnv() string {
	return os.Getenv("PORT")
}

func TestAnotherInputScanning(t *testing.T) {
	name := cmp.Or(getUserNameFromInput(), "Anonymous")
	fmt.Printf("Hello, %s\n", name)
}

func getUserNameFromInput() string {
	fmt.Print("Enter your name, 'Anonymous' as default")
	var name string
	if _, err := fmt.Scanln(&name); err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(name)
}
