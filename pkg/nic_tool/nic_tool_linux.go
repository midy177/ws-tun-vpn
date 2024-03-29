//go:build linux
// +build linux

package nic_tool

import (
	"strconv"
	"strings"
)

// set the cidr for the tun device
func (t *tool) SetCidrAndUp() string {
	info := []string{}
	i1 := execCmd("/sbin/ip", "addr", "add", t.cidr, "dev", t.tunName)
	if i1 != "" {
		info = append(info, i1)
	}
	i2 := execCmd("/sbin/ip", "link", "set", "dev", t.tunName, "up")
	if i2 != "" {
		info = append(info, i2)
	}
	return strings.Join(info, "\n")
}

// set the mtu for the tun device
func (t *tool) SetMtu() string {
	return execCmd("/sbin/ip", "link", "set", "dev", t.tunName, "mtu", strconv.Itoa(t.mtu))
}

// set the route for the tun device
func (t *tool) SetRoute(cidr string) string {
	return execCmd("/sbin/ip", "route", "add", cidr, "dev", t.tunName)
}
