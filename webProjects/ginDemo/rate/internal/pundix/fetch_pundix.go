package pundix

import (
	"encoding/json"
	"net/http"
	"own/gin/rate/internal"
)

// FetchPundiSupply will retrieve pundix supply from the given node url
func FetchPundiSupply(nodeUrl string) (*internal.SupplyApiResponse, error) {

	// request
	resp, err := http.Get(nodeUrl + internal.SUPPLY_PATH)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	// here use json.NewDecoder is better for not loading whole json response into the memory
	var apiResponse internal.SupplyApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if nil != err {
		return nil, err
	}
	return &apiResponse, nil
}
