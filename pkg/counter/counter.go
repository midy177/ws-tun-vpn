package counter

import (
	"sync/atomic"

	"github.com/inhies/go-bytesize"
)

// totalReadBytes is the total number of bytes read
var _totalReadBytes uint64 = 0

// totalWrittenBytes is the total number of bytes written
var _totalWrittenBytes uint64 = 0

// IncrReadBytes increments the number of bytes read
func IncrReadBytes(n int) {
	atomic.AddUint64(&_totalReadBytes, uint64(n))
}

// IncrWrittenBytes increments the number of bytes written
func IncrWrittenBytes(n int) {
	atomic.AddUint64(&_totalWrittenBytes, uint64(n))
}

// GetReadBytes returns the number of bytes read
func GetReadBytes() uint64 {
	return atomic.LoadUint64(&_totalReadBytes)
}

// GetWrittenBytes returns the number of bytes written
func GetWrittenBytes() uint64 {
	return atomic.LoadUint64(&_totalWrittenBytes)
}

// PrintBytes returns the bytes info
func PrintBytes(serverMode bool) (download, upload string) {
	if serverMode {
		download = bytesize.New(float64(GetWrittenBytes())).String()
		upload = bytesize.New(float64(GetReadBytes())).String()

	} else {
		download = bytesize.New(float64(GetReadBytes())).String()
		upload = bytesize.New(float64(GetWrittenBytes())).String()
	}
	return
}

// ResetBytes resets the bytes counters
func ResetBytes() {
	atomic.StoreUint64(&_totalReadBytes, 0)
	atomic.StoreUint64(&_totalWrittenBytes, 0)
}
