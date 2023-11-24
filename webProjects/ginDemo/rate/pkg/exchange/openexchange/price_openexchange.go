package exchangerate

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type Rates struct {
	IDR float64 `json:"IDR"`
	KRW float64 `json:"KRW"`
	SGD float64 `json:"SGD"`
	THB float64 `json:"THB"`
}

type OpenExchangeResponse struct {
	Disclaimer string `json:"disclaimer"`
	License    string `json:"license"`
	Timestamp  int64  `json:"timestamp"`
	Base       string `json:"base"`
	Rates      Rates  `json:"rates"`
}

const OpenExchangeLatestPriceUrl = "https://openexchangerates.org/api/latest.json"

// FetchOpenExchangePrice will retrieve exchange rate for IDR, KRW, SGD, THB base on USD
func FetchOpenExchangePrice(apiKey string) (*OpenExchangeResponse, error) {
	// construct
	parsedUrl, _ := url.Parse(OpenExchangeLatestPriceUrl)
	params := url.Values{}
	params.Add("app_id", apiKey)
	params.Add("base", "USD")
	params.Add("symbols", "KRW,IDR,SGD,THB")
	params.Add("prettyprint", "false")
	params.Add("show_alternative", "false")
	parsedUrl.RawQuery = params.Encode()

	// request
	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	defer resp.Body.Close()

	// parse
	var apiResponse OpenExchangeResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return &apiResponse, nil
}
