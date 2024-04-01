package loadlib

import (
	"runtime"
	"sync/atomic"
)

var once uint32

func LoadTunLib() error {
	if runtime.GOOS == "windows" && atomic.LoadUint32(&once) == 0 {
		return loadLib()
		atomic.StoreUint32(&once, 1)
		//defer syscall.FreeLibrary(h)
	}
	return nil
}
