package supply

import (
	"encoding/json"
	"fmt"
	"log"
	"own/gin/rate/internal/load"
	"testing"
)

func Test_FetchTarget(t *testing.T) {
	configuration, _ := load.LoadConfiguration("../../config/config.json")
	supply, _ := FetchTargetSupply(configuration)
	jsonStr, _ := json.MarshalIndent(supply, "", "  ")
	fmt.Println(string(jsonStr))
}

func Test_Fetch(t *testing.T) {
	fetchCalling("https://fx-rest.functionx.io")
}

func Test_FetchP(t *testing.T) {
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
