//go:build linux
// +build linux

package gateway

import (
	"log"
	"net"
	"os/exec"
	"strings"
)

func discoverGatewayOSSpecificIPv4() (ip net.IP, err error) {
	ipstr := execCmd("sh", "-c", "route -n | grep 'UG[ \t]' | awk 'NR==1{print $2}'")
	ipv4 := net.ParseIP(ipstr)
	if ipv4 == nil {
		return nil, errCantParse
	}
	return ipv4, nil
}

func discoverGatewayOSSpecificIPv6() (ip net.IP, err error) {
	ipstr := execCmd("sh", "-c", "route -6 -n | grep 'UG[ \t]' | awk 'NR==1{print $2}'")
	ipv6 := net.ParseIP(ipstr)
	if ipv6 == nil {
		return nil, errCantParse
	}
	return ipv6, nil
}

// execCmd executes a command and returns its output as a string.
func execCmd(c string, args ...string) string {
	cmd := exec.Command(c, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Println("failed to exec cmd:", err)
	}
	if len(out) == 0 {
		return ""
	}
	s := string(out)
	return strings.ReplaceAll(s, "\n", "")
}
