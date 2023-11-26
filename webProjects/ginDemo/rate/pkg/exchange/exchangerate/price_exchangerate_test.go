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

func Test_FetchQuotePrice(t *testing.T) {
	apiKey := ""
	fetch := ErFetcher{UsdPrice: float64(1), ApiKey: apiKey}
	prices, err := fetch.FetchConvertToQuotePrices()
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(*prices, "", "  ")
	fmt.Println(string(jsonStr))
}
