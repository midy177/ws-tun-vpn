//go:build linux
// +build linux

package nic_tool

import (
	"strconv"
)

// set the cidr for the tun device
func (t *tool) SetCidrAndUp(tunName, cidr string) string {
	return execCmd("/sbin/ip", "addr", "add", cidr, "dev", tunName) +
		"\n" +
		execCmd("/sbin/ip", "link", "set", "dev", tunName, "up")
}

// set the mtu for the tun device
func (t *tool) SetMtu(tunName string, mtu int) string {
	return execCmd("/sbin/ip", "link", "set", "dev", tunName, "mtu", strconv.Itoa(mtu))
}

// set the route for the tun device
func (t *tool) SetRoute(tunName, dst string) string {
	return execCmd("/sbin/ip", "route", "add", dst, "dev", tunName)
}
