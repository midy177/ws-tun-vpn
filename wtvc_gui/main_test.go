package main

import (
	"os"
	"syscall"
	"testing"
)

func TestName(t *testing.T) {
	process, err := os.FindProcess(41459)
	if err != nil {
		os.Exit(1)
	}
	err = process.Signal(syscall.SIGKILL)
	if err != nil {
		os.Exit(1)
	}
}
