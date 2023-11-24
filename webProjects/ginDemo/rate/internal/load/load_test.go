package load

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func Test_Loading(t *testing.T) {
	configuration, err := LoadConfiguration()
	if err != nil {
		log.Fatal(err)
	}
	jsonStr, _ := json.MarshalIndent(&configuration, "", "  ")
	fmt.Println(string(jsonStr))
}

func Test_Enum(t *testing.T) {
	var fxEnum ServeEnum = PundixServing
	fmt.Println(fxEnum)
}
