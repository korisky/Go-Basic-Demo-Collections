package supply

import (
	"encoding/json"
	"net/http"
)

// FetchSupply will retrieve denom (fx / pundix) supply from the given node url
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
