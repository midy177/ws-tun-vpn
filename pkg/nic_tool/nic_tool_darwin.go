//go:build darwin
// +build darwin

package nic_tool

import (
	"strconv"
	"strings"
	"ws-tun-vpn/pkg/netutil"
)

// set the cidr and up for the tun device,cidr example: 10.0.0.1/24
func (t *tool) SetCidrAndUp() string {
	cidrSlice := netutil.GetCidrV4SliceWithFatal(t.cidr)
	if len(cidrSlice) < 2 {
		return "invalid cidr"
	}
	return execCmd("ifconfig", t.tunName, "inet", strings.Split(t.cidr, "/")[1], cidrSlice[0], "up")
}

// set the mtu for the tun device
func (t *tool) SetMtu() string {
	return execCmd("ifconfig", t.tunName, "mtu", strconv.Itoa(t.mtu))
}

// set the route for the tun device
func (t *tool) SetRoute(cidr string) string {
	return execCmd("route", "add", cidr, "-interface", t.tunName)
}
