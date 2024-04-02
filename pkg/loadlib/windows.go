//go:build windows
// +build windows

package loadlib

import "syscall"

func loadLib() error {
	_, err := syscall.LoadLibrary("wintun.dll")
	return err
}
