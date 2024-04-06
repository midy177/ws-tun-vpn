package util

import (
	"fmt"
	"net"
	"testing"
	"ws-tun-vpn/pkg/gateway"
)

func TestUtil(t *testing.T) {

	pv4, err := gateway.DiscoverGatewayIPv4()
	if err != nil {
		return
	}
	fmt.Println("Gateway IPv4:", pv4)
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

func TestName(t *testing.T) {
	addresses := []string{
		"192.168.1.1:8080",
		"example.com:8080",
		"192.168.1.1",
		"example.com",
		"invalid-address",
	}

	for _, address := range addresses {
		if IsValidAddress(address) {
			fmt.Printf("%s 是有效的地址\n", address)
		} else {
			fmt.Printf("%s 是无效的地址\n", address)
		}
	}
}
