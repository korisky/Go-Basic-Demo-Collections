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
	checkError(err)

	// second check the resp error, if we do forget to close it (e.g. when it return nil), the whole program would panic and stop
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	checkError(err)

	fmt.Println(string(body))
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
