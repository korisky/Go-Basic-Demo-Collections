package coingecko

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

// Test_FetchCgQuote is unit test for fetch CoinGecko's fiatexchange cryptoexchange
func Test_FetchCgQuote(t *testing.T) {
	quote, err := FetchCgQuotePrice("fx-coin,pundi-x", "usd")
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(quote, "", "  ")
	fmt.Println(string(jsonStr))
}
