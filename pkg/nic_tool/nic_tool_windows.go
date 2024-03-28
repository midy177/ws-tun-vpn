//go:build windows
// +build windows

package nic_tool

import (
	"net"
	"strconv"
	"ws-tun-vpn/pkg/netutil"
)

// SetCidrAndUp Set the cidr for the tun device
func (t *tool) SetCidrAndUp(tunName, cidr string) string {
	return ""
}

// SetMtu Set the mtu for the tun device
func (t *tool) SetMtu(tunName string, mtu int) string {
	// netsh interface ipv4 set subinterface "YourTunInterfaceName" mtu=1400 store=persistent
	return execCmd("netsh", "interface", "ipv4", "set", "subinterface", tunName,
		"mtu="+strconv.Itoa(mtu), "store=persistent")
}

// SetRoute set the route for the tun device,set the distributed cidr and parse the cidr to get
// the first address as the tun device gateway and set the metric to 6
func (t *tool) SetRoute(tunName, cidr string) string {
	_, i, err := net.ParseCIDR(cidr)
	if err != nil {
		return err.Error()
	}
	cidrSlice := netutil.GetCidrV4SliceWithFatal(cidr)
	if len(cidrSlice) < 2 {
		return "invalid cidr"
	}
	return execCmd("cmd", "/C", "route", "add", "0.0.0.0", "mask", i.Mask.String(), cidrSlice[0], "metric", "6")
}
