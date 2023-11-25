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

// FetchCirculatingSupply will retrieve Denom circulating supply
func FetchCirculatingSupply(config *load.Config) (float64, error) {
	// fetch supply
	supplyResp, err := fetchSupply(config.NodeUrl)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	// extract supply
	supply, err := extractTargetDenomCirculatingSupply(supplyResp, config)
	if err != nil {
		return 0, err
	}
	// decimals
	return supply / math.Pow10(18), nil
}

// extractTargetDenomCirculatingSupply help extract target denom supply
func extractTargetDenomCirculatingSupply(supplyResp *SupplyApiResponse, config *load.Config) (float64, error) {
	for _, item := range supplyResp.Supply {
		if isRelevantDenom(item, config) {
			return strconv.ParseFloat(item.Amount, 64)
		}
	}
	return 0, errors.New("relevant denom not found")
}

// isRelevantDenom filter for target denom
func isRelevantDenom(item SupplyItem, config *load.Config) bool {
	switch {
	case
		load.FxServing == config.NodeServing && strings.EqualFold("fx", item.Denom):
		return true
	case
		load.PundixServing == config.NodeServing && strings.EqualFold("ibc/55367B7B6572631B78A93C66EF9FDFCE87CDE372CC4ED7848DA78C1EB1DCDD78", item.Denom):
		return true
	default:
		return false
	}
}

// fetchSupply will retrieve denom (fx / pundix) circulating supply from the given node url
func fetchSupply(nodeUrl string) (*SupplyApiResponse, error) {
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
