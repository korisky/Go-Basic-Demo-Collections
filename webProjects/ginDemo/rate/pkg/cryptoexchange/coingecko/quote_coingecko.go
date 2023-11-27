package coingecko

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"own/gin/rate/pkg/cryptoexchange"
)

type CgFetcher struct {
	Ids        string
	Currencies string
}

type CgApiResponse map[string]map[string]float64

const CoinGeckoQuoteUrl = "https://api.coingecko.com/api/v3/simple/price"

// FetchToUsdPrice only retrieve denom:usd exchange rate
func (c *CgFetcher) FetchToUsdPrice() (*cryptoexchange.ToUsdPrice, error) {
	// Fetch for response
	cgApiResponse, err := FetchCgQuotePrice(c.Ids, c.Currencies)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	// extract to usd exchange rate
	toUsdPrice := cryptoexchange.ToUsdPrice{DenomToUsd: (*cgApiResponse)[c.Ids]["usd"]}
	return &toUsdPrice, nil
}

// FetchCgQuotePrice will retrieve denom:usd exchange rate from CoinGecko
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
