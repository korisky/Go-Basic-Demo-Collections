package supply

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"own/gin/rate/internal/load"
	"strconv"
	"strings"
)

// FetchTargetSupply will retrieve Denom circulating supply
func FetchTargetSupply(config *load.Config) (float64, error) {
	// fetch supply
	supplyResp, err := FetchSupply(config.NodeUrl)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	// extract supply
	supply, err := extractTargetSupply(supplyResp, config)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	// decimals
	return supply / math.Pow10(18), nil
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
func extractTargetSupply(supplyResp *SupplyApiResponse, config *load.Config) (float64, error) {
	for _, item := range supplyResp.Supply {
		if isRelevantDenom(item, config) {
			return strconv.ParseFloat(item.Amount, 64)
		}
	}
	return 0, errors.New("relevant denom not found")
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
