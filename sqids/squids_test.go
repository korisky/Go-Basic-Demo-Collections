package squids

import (
	"fmt"
	"github.com/sqids/sqids-go"
	"testing"
)

func Test_Squids(t *testing.T) {
	// init
	s, _ := sqids.New(sqids.Options{MinLength: 10})

	// encode
	id, _ := s.Encode([]uint64{1, 3, 1234234}) // "86Rf07"
	fmt.Println(id)

	// decode
	numbers := s.Decode(id) // [1, 2, 3]
	fmt.Println(numbers)
}
