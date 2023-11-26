package openexchange

import (
	"encoding/json"
	"net/http"
	"net/url"
	"own/gin/rate/pkg/exchange"
)

type OxFetcher struct {
	UsdPrice float64
	ApiKey   string
}

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

const oxLatestPriceUrl = "https://openexchangerates.org/api/latest.json"

// FetchConvertToQuotePrices implementation for OpenExchange
func (o *OxFetcher) FetchConvertToQuotePrices() (*exchange.QuotePrices, error) {
	// fetch
	openExchangeResp, err := FetchOpenExchangePrice(o.ApiKey)
	if err != nil {
		return nil, err
	}
	// extract
	prices := exchange.QuotePrices{
		ToUSD: o.UsdPrice,
		ToSGD: openExchangeResp.Rates.SGD * o.UsdPrice,
		ToTHB: openExchangeResp.Rates.THB * o.UsdPrice,
		ToKRW: openExchangeResp.Rates.KRW * o.UsdPrice,
		ToIDR: openExchangeResp.Rates.IDR * o.UsdPrice,
	}
	return &prices, nil
}

// FetchOpenExchangePrice will retrieve exchange rate for IDR, KRW, SGD, THB base on USD, from OpenExchangeRates.org
func FetchOpenExchangePrice(apiKey string) (*OxApiResponse, error) {
	// construct
	parsedUrl, _ := url.Parse(oxLatestPriceUrl)
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
		return nil, err
	}
	defer resp.Body.Close()

	// parse
	var apiResponse OxApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
