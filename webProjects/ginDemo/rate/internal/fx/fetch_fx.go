package fx

import (
	"io"
	"net/http"
)

// FetchFxSupply will retrieve fx supply from the given node url
func FetchFxSupply(nodeUrl string) (string, error) {

	resp, err := http.Get(nodeUrl + "/cosmos/bank/v1beta1/supply")
	if nil != err {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if nil != err {
		return "", err
	}

	return string(body), nil
}
