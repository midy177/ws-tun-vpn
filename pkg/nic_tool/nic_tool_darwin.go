//go:build darwin
// +build darwin

package nic_tool

import (
	"strconv"
	"ws-tun-vpn/pkg/netutil"
	"ws-tun-vpn/pkg/util"
)

// set the cidr and up for the tun device,cidr example: 10.0.0.1/24
func (t *tool) SetCidrAndUp() string {
	cidrSlice := netutil.GetCidrV4SliceWithFatal(t.cidr)
	if len(cidrSlice) < 2 {
		return "invalid cidr"
	}
	//  golang string to int
	ipa, _ := util.CidrAddrToIPAddrAndMask(t.cidr)
	return execCmd("ifconfig", t.tunName, "inet", ipa, cidrSlice[0], "up")
	//sudo ("ifconfig", iFace.Name(), "inet", ip.String(), config.ServerIP, "up")
}

// set the mtu for the tun device
func (t *tool) SetMtu() string {
	return execCmd("ifconfig", t.tunName, "mtu", strconv.Itoa(t.mtu))
}

// set the route for the tun device
func (t *tool) SetRoute(cidr string) string {
	return execCmd("route", "add", cidr, "-interface", t.tunName)
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

//func (t *tool) ReleaseDevice() string {
//	return execCmd("netsh", "interface", "set", "interface", t.tunName, "disable")
//}
