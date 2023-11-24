package openexchange

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func Test_FetchOpenExchange(t *testing.T) {
	apiKey := ""
	price, err := FetchOpenExchangePrice(apiKey)
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(price, "", "  ")
	fmt.Println(string(jsonStr))
}
