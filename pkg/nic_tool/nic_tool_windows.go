//go:build windows
// +build windows

package nic_tool

import (
	"strconv"
	"ws-tun-vpn/pkg/util"
)

// SetCidrAndUp Set the cidr for the tun device
func (t *tool) SetCidrAndUp() string {
	return ""
}

// SetMtu Set the mtu for the tun device
func (t *tool) SetMtu() string {
	// netsh interface ipv4 set subinterface "YourTunInterfaceName" mtu=1400 store=persistent
	return execCmd("netsh", "interface", "ipv4", "set", "subinterface", t.tunName,
		"mtu="+strconv.Itoa(t.mtu), "store=persistent")
}

// SetRoute set the route for the tun device,set the distributed cidr and parse the cidr to get
// the first address as the tun device gateway and set the metric to 6
func (t *tool) SetRoute(cidr string) string {
	ip, mask := util.CidrToIPAndMask(cidr)
	gatewayIp, _ := util.CidrToIPAndMask(t.cidr)
	return execCmd("cmd", "/C", "route", "add", ip, "mask", mask, gatewayIp, "metric", "6")
}
