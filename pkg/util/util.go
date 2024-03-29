package util

import (
	"fmt"
	validate "github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"net"
	"strconv"
)

func ValidateWithFatal(field any, tag, flag string) {
	err := validate.New().Var(field, tag)
	if err != nil {
		log.Fatalf(" validation failed for flag: %s, failed on the '%s' tag", flag, tag)
	}
}

func FlagRequiredWithFatal(cmd *cobra.Command, name string) {
	err := cmd.MarkFlagRequired(name)
	if err != nil {
		log.Fatal(err)
	}
}

// GenerateTunName Randomly generate tun network card name
func GenerateTunName(n int) string {
	const charset = "0123456789"
	length := len(charset)
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(length)]
	}
	return "utun" + string(b)
}

func CidrAddrToIPAddrAndMask(cidr string) (ipAddr, mask string) {
	// 解析 CIDR 地址
	i, in, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Printf("Failed to parse CIDR: %v\n", err)
		return "", ""
	}
	ipAddr = i.String()
	mask = net.IP(in.Mask).String()
	return
}

// CidrToIPNetAndMask 将 CIDR (192.168.1.10/24)地址转换为 IPNet(192.168.1.0) 和掩码 (255.255.255.0)
func CidrToIPNetAndMask(cidr string) (ipNet, mask string) {
	_, in, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Printf("Failed to parse CIDR: %v\n", err)
		return "", ""
	}
	ipNet = in.IP.String()
	mask = net.IP(in.Mask).String()
	return
}

// CidrToIPNet 将 CIDR (192.168.1.10/24)地址转换为 IPNet和掩码(192.168.1.0/24)
func CidrToIPNet(cidr string) (ipNet string) {
	_, in, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Printf("Failed to parse CIDR: %v\n", err)
		return ""
	}
	ones, _ := in.Mask.Size()
	ipNet = in.IP.String() + "/" + strconv.Itoa(ones)
	return
}

// GetDefaultInterfaceName 获取默认网卡名称
func GetDefaultInterfaceName() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 &&
			iface.Flags&net.FlagPointToPoint == 0 &&
			iface.Flags&net.FlagRunning != 0 {
			fmt.Println("Name:", iface.Name)
			return iface.Name
		}
	}
	return ""
}
