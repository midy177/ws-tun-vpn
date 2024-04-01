//go:build linux
// +build linux

package nic_tool

import (
	"strconv"
	"strings"
	"ws-tun-vpn/pkg/util"
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

// Enable IP forwarding
func (t *tool) EnableIpForward() string {
	return execCmd("/sbin/sysctl", "-w", "net.ipv4.ip_forward=1")
}

// Enable NAT forwarding for the tun device
func (t *tool) EnableNat() string {
	ipNet := util.CidrToIPNet(t.cidr)
	cn := util.GetDefaultInterfaceName()
	return execCmd("/sbin/iptables", "-t", "nat", "-A", "POSTROUTING", "-s", ipNet, "-o", cn, "-j", "MASQUERADE")
}

//func (t *tool) ReleaseDevice() string {
//	return execCmd("netsh", "interface", "set", "interface", t.tunName, "disable")
//}
