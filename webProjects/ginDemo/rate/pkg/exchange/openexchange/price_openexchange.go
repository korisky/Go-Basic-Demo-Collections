package openexchange

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

type OxApiResponse struct {
	Disclaimer string `json:"disclaimer"`
	License    string `json:"license"`
	Timestamp  int64  `json:"timestamp"`
	Base       string `json:"base"`
	Rates      Rates  `json:"rates"`
}

const OxLatestPriceUrl = "https://openexchangerates.org/api/latest.json"

// FetchOpenExchangePrice will retrieve exchange rate for IDR, KRW, SGD, THB base on USD, from OpenExchangeRates.org
func FetchOpenExchangePrice(apiKey string) (*OxApiResponse, error) {
	// construct
	parsedUrl, _ := url.Parse(OxLatestPriceUrl)
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
	var apiResponse OxApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return &apiResponse, nil
}
