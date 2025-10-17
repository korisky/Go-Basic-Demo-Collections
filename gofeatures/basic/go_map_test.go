package basic

import (
	"fmt"
	"testing"
)

type TheVertex struct {
	Lat, Long float64
}

// Embedded Initialization
var m = map[string]TheVertex{
	"Bell labs": {
		Lat:  40.5823,
		Long: -75.3242,
	},
	"Google": {
		Lat:  26.5920,
		Long: 129.2342,
	},
}

func TestMapInit(t *testing.T) {
	fmt.Println(m)
}

// TestMapContaining 可以看出，类似interface中我们可以使用断言确认类型并获取值
// map中同样可以, 简单直接同时获取值+是否存在
func TestMapContaining(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "123"
	m["b"] = "456"

	sa, ok := m["a"]
	fmt.Println(sa, ok)

	sc, ok := m["c"]
	fmt.Println(sc, ok)

	scc := m["c"]
	fmt.Println(scc)
}
