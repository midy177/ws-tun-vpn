//go:build darwin
// +build darwin

package nic_tool

import (
	"strconv"
	"strings"
	"ws-tun-vpn/pkg/netutil"
)

func newTool(tunName, cidr string, mtu int) *tool {
	return &tool{tunName: tunName, cidr: cidr, mtu: mtu}
}

// set the cidr and up for the tun device,cidr example: 10.0.0.1/24
func SetCidrAndUp(tunName, cidr string) string {
	cidrSlice := netutil.GetCidrV4SliceWithFatal(cidr)
	if len(cidrSlice) < 2 {
		return "invalid cidr"
	}
	return execCmd("ifconfig", tunName, "inet", strings.Split(cidr, "/")[1], cidrSlice[0], "up")
}

// set the mtu for the tun device
func SetMtu(tunName string, mtu int) string {
	return execCmd("ifconfig", tunName, "mtu", strconv.Itoa(mtu))
}

// set the route for the tun device
func SetRoute(tunName, dst string) string {
	return execCmd("route", "add", dst, "-interface", tunName)
}
