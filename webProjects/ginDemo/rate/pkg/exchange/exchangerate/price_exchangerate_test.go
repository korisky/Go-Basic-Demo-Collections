package exchangerate

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func Test_FetchExchangeRatePrice(t *testing.T) {
	apiKey := ""
	price, err := FetchExchangeRatePrice(apiKey)
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(price, "", "  ")
	fmt.Println(string(jsonStr))
}
