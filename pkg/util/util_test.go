package util

import (
	"fmt"
	"net"
	"testing"
)

func TestUtil(t *testing.T) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 &&
			iface.Flags&net.FlagPointToPoint == 0 &&
			iface.Flags&net.FlagRunning != 0 {
			fmt.Println("Name:", iface.Name)
		}
	}
	ipNet := GetDefaultInterfaceName()
	fmt.Println(ipNet)
}
