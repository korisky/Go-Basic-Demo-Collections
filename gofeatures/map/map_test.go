package test_map

import (
	"fmt"
	"maps"
	"strings"
	"testing"
)

func Test_clone(t *testing.T) {

	originalMap := map[string]int{"a": 1, "b": 2, "c": 3}

	// clone
	clonedMap := maps.Clone(originalMap)

	fmt.Println("Original Map:", originalMap)
	fmt.Println("Cloned Map:", clonedMap)
}

func Test_copy(t *testing.T) {
	dstMap := map[string]int{"a": 1, "b": 2}
	srcMap := map[string]int{"b": 3, "c": 4}

	// dstMap's same key element will be overwritten
	maps.Copy(dstMap, srcMap)

	fmt.Println("Destination Map:", dstMap)
}

func Test_delete(t *testing.T) {

	theMap := map[string]int{"apple": 1, "banana": 2, "cherry": 3, "date": 4}

	// the delete func -> return true -> then delete the 'true' elements
	maps.DeleteFunc(theMap, func(key string, val int) bool {
		return strings.Contains(key, "bana")
	})

	fmt.Println("After delete,", theMap)
}
