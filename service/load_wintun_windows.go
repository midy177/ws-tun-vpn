//go:build windows
// +build windows

package service

import (
	"fmt"
	"syscall"
)

func init() {
	dll := "wintun.dll"
	_, err := syscall.LoadLibrary(dll)
	if err != nil {
		fmt.Printf("Error loading %s DLL: %v\n", dll, err)
		return
	}
	//defer syscall.FreeLibrary(h)
	fmt.Printf("%s loaded successfully.\n", dll)
}
