package exchangerate

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func Test_FetchOpenExchange(t *testing.T) {
	apiKey := ""
	quote, err := FetchOpenExchangePrice(apiKey)
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(quote, "", "  ")
	fmt.Println(string(jsonStr))
}
