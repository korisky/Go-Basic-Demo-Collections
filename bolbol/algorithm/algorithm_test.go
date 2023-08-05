package algorithm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
	Golang's rule for tests ->
		tests calling tests
*/

func TestCheckEveryItem(t *testing.T) {
	testAlgorithm(CheckEveryItem, t)
}

func TestBinarySearch(t *testing.T) {
	testAlgorithm(BinarySearch, t)
}

func testAlgorithm(alg Algorithm, t *testing.T) {
	items := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i := 0; i <= 9; i++ {
		assert.Equal(t, i, alg(items, i))
	}
	assert.Equal(t, -1, alg(items, 100))
}
