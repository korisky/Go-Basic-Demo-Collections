package pundix

import (
	"encoding/json"
	"net/http"
	"own/gin/rate/pkg/supply"
)

// FetchPundiSupply will retrieve pundix supply from the given node url
func FetchPundiSupply(nodeUrl string) (*supply.SupplyApiResponse, error) {

	// request
	resp, err := http.Get(nodeUrl + supply.SUPPLY_PATH)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse
	var apiResponse supply.SupplyApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
