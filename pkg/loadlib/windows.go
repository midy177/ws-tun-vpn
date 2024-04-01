//go:build windows
// +build windows

package loadlib

func loadLib() error {
	_, err := syscall.LoadLibrary("wintun.dll")
	return err
}
