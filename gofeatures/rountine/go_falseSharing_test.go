package main

// CountersWithoutPadding will cause 'false-sharing'
type CountersWithoutPadding struct {
	requests uint64
	errors   uint64
	latency  uint64
}

// CountersWithPadding since with extra padding, fill the whole cache line
// then it will not cause 'false-sharing', and result in a greater performance
type CountersWithPadding struct {
	requests uint64
	_        [120]byte
	errors   uint64
	_        [128]byte
	latency  uint64
	_        [128]byte
}

type Counters interface {
	RequestsPtr() *uint64
	ErrorPtr() *uint64
	LatencyPtr() *uint64
}
