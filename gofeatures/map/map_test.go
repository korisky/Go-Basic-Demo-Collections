package test_map

import (
	"fmt"
	"maps"
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
