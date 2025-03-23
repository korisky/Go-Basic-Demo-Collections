package streamActions

import "iter"

// Iterator is a struct for wrapping pure iter impl
type Iterator[V any] struct {
	iter iter.Seq[V]
}

// Collect is for collect back the stream
func (i Iterator[V]) Collect() []V {
	collect := make([]V, 0)
	for e := range i.iter {
		collect = append(collect, e)
	}
	return collect
}

// From help for wrapping the array into stream
func From[V any](slice []V) *Iterator[V] {
	return &Iterator[V]{
		iter: func(yield func(V) bool) {
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		},
	}
}

// Each is for mapping
func (i *Iterator[V]) Each(f func(V)) {
	for i := range i.iter {
		f(i)
	}
}

// Reverse the stream
func (i Iterator[V]) Reverse() *Iterator[V] {
	collect := i.Collect()
	counter := len(collect) - 1
	for e := range i.iter {
		collect[counter] = e
		counter--
	}
	return From(collect)
}

// Filter is for filtering
func (i *Iterator[V]) Filter(f func(V) bool) *Iterator[V] {
	cpy := i.iter
	i.iter = func(yield func(V) bool) {
		for v := range cpy {
			if f(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
	return i
}
