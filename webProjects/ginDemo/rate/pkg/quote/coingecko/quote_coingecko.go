package coingecko

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type CgApiResponse map[string]map[string]float64

const CoinGeckoQuoteUrl = "https://api.coingecko.com/api/v3/simple/price"

// FetchCgQuotePrice will retrieve exchange quote from CoinGecko
func FetchCgQuotePrice(ids, vsCurrencies string) (*CgApiResponse, error) {
	// construct
	parsedUrl, _ := url.Parse(CoinGeckoQuoteUrl)
	params := url.Values{}
	params.Add("ids", ids)
	params.Add("vs_currencies", vsCurrencies)
	parsedUrl.RawQuery = params.Encode()
	s := parsedUrl.String()

	// request
	resp, err := http.Get(s)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse
	var apiResponse CgApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
