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
	ipAddr, mask := util.CidrToIPNetAndMask(cidr)
	gatewayAddr, _ := util.CidrAddrToIPAddrAndMask(t.cidr)
	return execCmd("cmd", "/C", "route", "add", ipAddr, "mask", mask, gatewayAddr, "metric", "6")
}

// Enable IP forwarding
func (t *tool) EnableIpForward() string {
	return ""
}

// Enable NAT forwarding for the tun device
func (t *tool) EnableNat() string {
	//ipNet := util.CidrToIPNet(t.cidr)
	//cn := util.GetDefaultInterfaceName()
	return ""
}

func (t *tool) DisableNat() string {
	return ""
}

func (t *tool) SetPrimaryDnsServer(dns string) string {
	return execCmd("netsh", "interface", "ipv4", "set", "dnsservers", t.tunName, "static", dns, "primary")
}
