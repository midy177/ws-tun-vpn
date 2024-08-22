//go:build windows
// +build windows

package nic_tool

import (
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
	"net/netip"
	"strconv"
	"syscall"
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
	//ipAddr, mask := util.CidrToIPNetAndMask(cidr)
	gatewayAddr, _ := util.CidrAddrToIPAddrAndMask(t.cidr)

	dev := t.iFac.GetDev()
	if dev == nil {
		return "获取设备失败"
	}
	nextHop, err := netip.ParseAddr(gatewayAddr)
	if err != nil {
		return err.Error()
	}

	link := winipcfg.LUID(dev.LUID())

	err = link.AddRoute(netip.MustParsePrefix(cidr), nextHop, 6)
	if err != nil {
		return err.Error()
	}
	return "add route success"
	//return execCmd("cmd", "/C", "route", "add", ipAddr, "mask", mask, gatewayAddr, "metric", "6")
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
	dev := t.iFac.GetDev()
	if dev == nil {
		return "获取设备失败"
	}
	link := winipcfg.LUID(dev.LUID())
	err := link.FlushDNS(syscall.AF_INET)
	if err != nil {
		return err.Error()
	}
	ns, err := netip.ParseAddr(dns)
	if err != nil {
		return err.Error()
	}
	domains := []string{"wintun.dns"}
	err = link.SetDNS(windows.AF_INET, []netip.Addr{ns}, domains)
	if err != nil {
		return err.Error()
	}
	return "change dns success"
}
