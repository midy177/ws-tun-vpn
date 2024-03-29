package util

import (
	"fmt"
	"net"
	"os"
	"testing"
)

func TestUtil(t *testing.T) {
	cidr := "172.28.0.0/16"

	// 解析 CIDR 地址
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Printf("Failed to parse CIDR: %v\n", err)
		os.Exit(1)
	}

	// ip 是网络地址
	fmt.Println("IP:", ip.String())
	mask := net.IP(ipnet.Mask).String()
	// Mask 是子网掩码
	fmt.Println("Mask:", mask)
}
