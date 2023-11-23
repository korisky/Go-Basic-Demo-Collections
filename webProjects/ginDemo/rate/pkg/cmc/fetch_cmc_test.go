package cmc

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

var apiKey = ""

// Test_FetchCmc is unit test for cmc exchange rate retrieving
func Test_FetchCmc(t *testing.T) {
	quote, err := FetchCmcQuote(apiKey, "3884", "2781")
	if err != nil {
		log.Fatalln(err)
		return
	}
	jsonStr, _ := json.MarshalIndent(quote, "", "  ")
	fmt.Println(string(jsonStr))
}
