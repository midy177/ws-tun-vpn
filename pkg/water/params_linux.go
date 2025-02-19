//go:build linux
// +build linux

package water

// DevicePermissions determines the owner and group owner for the newly created
// interface.
type DevicePermissions struct {
	// Owner is the ID of the user which will be granted ownership of the
	// device.  If set to a negative value, the owner value will not be
	// changed.  By default, Linux sets the owner to -1, which allows any user.
	Owner uint

	// Group is the ID of the group which will be granted access to the device.
	// If set to a negative value, the group value will not be changed.  By
	// default, Linux sets the group to -1, which allows any group.
	Group uint
}

// PlatformSpecificParams defines parameters in Config that are specific to
// Linux. A zero-value of such type is valid, yielding an interface
// with OS defined name.
type PlatformSpecificParams struct {
	// Name is the name to be set for the interface to be created. This overrides
	// the default name assigned by OS such as tap0 or tun0. A zero-value of this
	// field, i.e. an empty string, indicates that the default name should be
	// used.
	Name string

	// Persist specifies whether persistence mode for the interface device
	// should be enabled or disabled.
	Persist bool

	// Permissions, if non-nil, specifies the owner and group owner for the
	// interface.  A zero-value of this field, i.e. nil, indicates that no
	// changes to owner or group will be made.
	Permissions *DevicePermissions

	// MultiQueue specifies whether the multiqueue flag should be set on the
	// interface.  From version 3.8, Linux supports multiqueue tuntap which can
	// uses multiple file descriptors (queues) to parallelize packets sending
	// or receiving.
	MultiQueue bool
	// Network is required when creating a TUN interface. The library will call
	// net.ParseCIDR() to parse this string into LocalIP, RemoteNetaddr,
	// RemoteNetmask. The underlying driver will need those to generate ARP
	// response to Windows kernel, to emulate an TUN interface.
	// Please note that it cannot perceive the IP changes caused by DHCP, user
	// configuration to the adapter and etc,. If IP changed, please reconfigure
	// the adapter using syscall, just like openDev().
	// For detail, please refer
	// https://github.com/OpenVPN/tap-windows6/blob/master/src/device.c#L431
	// and https://github.com/songgao/water/pull/13#issuecomment-270341777
	Network []string
	// Mtu is the maximum transmission unit of the interface.
	Mtu int
}

func defaultPlatformSpecificParams() PlatformSpecificParams {
	return PlatformSpecificParams{}
}
