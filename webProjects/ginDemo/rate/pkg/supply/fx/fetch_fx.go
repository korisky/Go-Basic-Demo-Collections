package fx

import (
	"encoding/json"
	"net/http"
	"own/gin/rate/pkg/supply"
)

// FetchFxSupply will retrieve fx supply from the given node url
func FetchFxSupply(nodeUrl string) (*supply.SupplyApiResponse, error) {

	// request
	resp, err := http.Get(nodeUrl + supply.SUPPLY_PATH)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// here use json.NewDecoder is better for not loading whole json response into the memory
	var apiResponse supply.SupplyApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}
	return &apiResponse, nil
}
