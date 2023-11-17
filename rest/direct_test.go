package rest

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func Test_directly_calling(t *testing.T) {

	resp, err := http.Get("https://fx-rest.functionx.io/cosmos/bank/v1beta1/supply")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//var jsonData json.RawMessage
	//err = json.Unmarshal(body, &jsonData)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

	fmt.Printf("Json Data: %+v\n", string(body))

}
