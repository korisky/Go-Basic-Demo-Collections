package coinmarketcap

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"own/gin/rate/pkg/cryptoexchange"
)

const reqUrl = "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest"

type CmcFetcher struct {
	Id        string
	ConvertId string
	ApiKey    string
}

func (c *CmcFetcher) FetchToUsdPrice() (*cryptoexchange.ToUsdPrice, error) {
	// fetch
	resp, err := FetchCmcQuote(c.ApiKey, c.Id, c.ConvertId)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	// extract
	cryptoData := resp.Data[c.Id]
	quoteData := cryptoData.Quote[c.ConvertId]
	toUsdPrice := cryptoexchange.ToUsdPrice{DenomToUsd: quoteData.Price}
	return &toUsdPrice, nil
}

// FetchCmcQuote will retrieve denom : fiat from CMC
func FetchCmcQuote(apiKey string, id string, convertId string) (*CmcApiResponse, error) {
	// construct
	parsedUrl, _ := url.Parse(reqUrl)
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
