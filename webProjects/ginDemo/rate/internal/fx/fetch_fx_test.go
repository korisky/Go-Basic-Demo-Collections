package fx

import (
	"fmt"
	"testing"
)

// Test_FetchFx is unit test for fx supply fetching
func Test_FetchFx(t *testing.T) {
	supply, err := FetchFxSupply("https://fx-rest.functionx.io")
	if err == nil {
		fmt.Println(supply)
	}
}
