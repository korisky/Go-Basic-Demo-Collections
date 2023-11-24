package fx

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

// Test_FetchFx is unit test for fx supply fetching
func Test_FetchFx(t *testing.T) {
	supply, err := FetchFxSupply("https://fx-rest.functionx.io")
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(supply, "", "  ")
	fmt.Println(string(jsonStr))
}
