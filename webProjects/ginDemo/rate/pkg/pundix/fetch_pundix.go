package pundix

import (
	"encoding/json"
	"net/http"
	"own/gin/rate/pkg"
)

// FetchPundiSupply will retrieve pundix supply from the given node url
func FetchPundiSupply(nodeUrl string) (*pkg.SupplyApiResponse, error) {

	// request
	resp, err := http.Get(nodeUrl + pkg.SUPPLY_PATH)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	// parse
	var apiResponse pkg.SupplyApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if nil != err {
		return nil, err
	}
	return &apiResponse, nil
}
