package fx

import (
	"encoding/json"
	"net/http"
	"own/gin/rate/pkg"
)

// FetchFxSupply will retrieve fx supply from the given node url
func FetchFxSupply(nodeUrl string) (*pkg.SupplyApiResponse, error) {

	// request
	resp, err := http.Get(nodeUrl + pkg.SUPPLY_PATH)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	// here use json.NewDecoder is better for not loading whole json response into the memory
	var apiResponse pkg.SupplyApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if nil != err {
		return nil, err
	}
	return &apiResponse, nil
}
