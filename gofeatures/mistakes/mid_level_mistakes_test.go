package mistakes

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func Test_http_close(t *testing.T) {
	// first check the request error
	resp, err := http.Get("https://www.google.com")
	if resp != nil {
		// is resp is not nil, do not forget the close it later
		// if forget to close it (e.g. when it return nil), the whole program would panic and stop
		defer resp.Body.Close()
	}
	checkError(err)

	body, err := io.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

func Test_recover(t *testing.T) {
	defer func() {
		fmt.Println("recovered: ", recover())
	}()
	panic("Not good")
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
