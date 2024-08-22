//go:build windows
// +build windows

package water

import (
	"net/netip"

	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
)

type wintun struct {
	dev tun.Device
}

func (w *wintun) Close() error {
	return w.dev.Close()
}

func (w *wintun) Write(b []byte) (int, error) {
	return w.dev.Write(b, 0)
}

func (w *wintun) Read(b []byte) (int, error) {
	return w.dev.Read(b, 0)
}

func openDev(config Config) (ifce *Interface, err error) {
	if config.DeviceType == TAP {
		return nil, err
	}
	//id := &windows.GUID{
	//	0x0000000,
	//	0xFFFF,
	//	0xFFFF,
	//	[8]byte{0xFF, 0xe9, 0x76, 0xe5, 0x8c, 0x74, 0x06, 0x3e},
	//}
	dev, err := tun.CreateTUN(config.PlatformSpecificParams.Name, config.PlatformSpecificParams.Mtu)
	//dev, err := tun.CreateTUNWithRequestedGUID(config.PlatformSpecificParams.Name, id, config.PlatformSpecificParams.Mtu)
	if err != nil {
		return nil, err
	}
	nativeTunDevice := dev.(*tun.NativeTun)
	link := winipcfg.LUID(nativeTunDevice.LUID())
	networks := config.PlatformSpecificParams.Network
	if len(networks) == 0 {
		panic("network is empty")
	}
	// set ip addresses
	ipPrefix := []netip.Prefix{}
	for _, n := range networks {
		ip, err := netip.ParsePrefix(n)
		if err != nil {
			panic(err)
		}
		ipPrefix = append(ipPrefix, ip)
	}
	err = link.SetIPAddresses(ipPrefix)
	if err != nil {
		panic(err)
	}
	// set dns
	servers := []netip.Addr{}
	for index, ns := range config.PlatformSpecificParams.NameServer {
		if index > 2 {
			break
		}
		nsAddr, err := netip.ParseAddr(ns)
		if err != nil {
			return nil, err
		}
		servers = append(servers, nsAddr)
	}
	domains := []string{"wintun.dns"}
	err = link.SetDNS(windows.AF_INET, servers, domains)
	if err != nil {
		panic(err)
	}

	tunName, err := dev.Name()
	ifce = &Interface{isTAP: config.DeviceType == TAP, ReadWriteCloser: &wintun{dev: dev}, name: tunName}
	return ifce, nil
}
