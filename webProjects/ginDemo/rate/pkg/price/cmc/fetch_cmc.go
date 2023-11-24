package cmc

import (
	"encoding/json"
	"net/http"
	"net/url"
)

const CoinMarketCapQuoteUrl = "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"

// FetchCmcQuote will retrieve exchange price from CMC
func FetchCmcQuote(apiKey string, id string, convertId string) (*CmcApiResponse, error) {
	// construct
	parsedUrl, _ := url.Parse(CoinMarketCapQuoteUrl)
	params := url.Values{}
	params.Add("id", id)
	params.Add("convert_id", convertId)
	parsedUrl.RawQuery = params.Encode()

	request, _ := http.NewRequest("GET", parsedUrl.String(), nil)
	request.Header.Add("X-CMC_PRO_API_KEY", apiKey)

	// request
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse
	var apiResponse CmcApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
