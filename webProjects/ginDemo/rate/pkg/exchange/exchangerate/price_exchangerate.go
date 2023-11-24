package exchangerate

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErApiResponse struct {
	Result             string             `json:"result"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	BaseCode           string             `json:"base_code"`
	ConversionRates    map[string]float64 `json:"conversion_rates"`
}

const ErLatestPriceUrl = "https://v6.exchangerate-api.com/v6/%s/latest/%s"

// FetchExchangeRatePrice will retrieve exchange rate for IDR, KRW, SGD, THB base on USD, from ExchangeRate-api.com
func FetchExchangeRatePrice(apiKey string) (*ErApiResponse, error) {
	// request
	parsedUrl := fmt.Sprintf(ErLatestPriceUrl, apiKey, "USD")
	resp, err := http.Get(parsedUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse
	var apiResponse ErApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return &apiResponse, nil
}
