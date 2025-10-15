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
