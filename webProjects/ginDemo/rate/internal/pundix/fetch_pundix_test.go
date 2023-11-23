package pundix

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

// Test_FetchPundix is unit test for pundix supply fetching
func Test_FetchPundix(t *testing.T) {
	supply, err := FetchPundiSupply("https://fx-rest.functionx.io")
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(supply, "", "  ")
	fmt.Println(string(jsonStr))
}
