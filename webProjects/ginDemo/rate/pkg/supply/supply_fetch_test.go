package supply

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func Test_FetchFx(t *testing.T) {
	fetchCalling("https://fx-rest.functionx.io")
}

func Test_FetchPundix(t *testing.T) {
	fetchCalling("https://px-rest.pundix.com")
}

func fetchCalling(nodeUrl string) {
	supply, err := FetchSupply(nodeUrl)
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(supply, "", "  ")
	fmt.Println(string(jsonStr))
}
