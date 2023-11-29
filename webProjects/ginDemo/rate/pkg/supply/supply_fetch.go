package supply

import (
	"encoding/json"
	"errors"
	"github.com/shopspring/decimal"
	"log"
	"math/big"
	"net/http"
	"own/gin/rate/internal/load"
	"strings"
)

// FetchTargetSupply will retrieve Denom circulating supply
func FetchTargetSupply(config *load.Config) (decimal.Decimal, error) {
	// fetch supply
	supplyResp, err := FetchSupply(config.NodeUrl)
	if err != nil {
		log.Fatal(err)
		return decimal.Decimal{}, err
	}
	// extract supply
	rawSupply, err := extractTargetSupply(supplyResp, config)
	if err != nil {
		log.Fatal(err)
		return decimal.Decimal{}, err
	}
	divisor := decimal.NewFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), 0)
	supply := rawSupply.Div(divisor)
	return supply, nil
}

// FetchSupply will retrieve denom (fx / pundix) circulating supply from the given node url
func FetchSupply(nodeUrl string) (*SupplyApiResponse, error) {
	// request
	resp, err := http.Get(nodeUrl + SUPPLY_PATH)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// here use json.NewDecoder is better for not loading whole json response into the memory
	var apiResponse SupplyApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse, nil
}

// extractTargetSupply help extract target denom supply
func extractTargetSupply(supplyResp *SupplyApiResponse, config *load.Config) (decimal.Decimal, error) {
	for _, item := range supplyResp.Supply {
		if isRelevantDenom(item, config) {
			return decimal.NewFromString(item.Amount)
		}
	}
	return decimal.Decimal{}, errors.New("relevant denom not found")
}

// isRelevantDenom filter for target denom
func isRelevantDenom(item SupplyItem, config *load.Config) bool {
	switch config.NodeServing {
	case load.FxServing:
		return strings.EqualFold(string(load.FxServing), item.Denom)
	case load.PundixServing:
		return strings.EqualFold("ibc/", item.Denom) // TODO replace to target ibc
	default:
		return false
	}
}
