//go:build windows
// +build windows

package gateway

import (
	"net"
	"strings"
)

type windowsRouteStructIPv4 struct {
	Destination string
	Netmask     string
	Gateway     string
	Interface   string
	Metric      string
}

type windowsRouteStructIPv6 struct {
	If          string
	Metric      string
	Destination string
	Gateway     string
}

func parseToWindowsRouteStructIPv4(output []byte) (windowsRouteStructIPv4, error) {
	// Windows route output format is always like this:
	// ===========================================================================
	// Interface List
	// 8 ...00 12 3f a7 17 ba ...... Intel(R) PRO/100 VE Network Connection
	// 1 ........................... Software Loopback Interface 1
	// ===========================================================================
	// IPv4 Route Table
	// ===========================================================================
	// Active Routes:
	// Network Destination        Netmask          Gateway       Interface  Metric
	//           0.0.0.0          0.0.0.0      192.168.1.1    192.168.1.100     20
	// ===========================================================================
	//
	// Windows commands are localized, so we can't just look for "Active Routes:" string
	// I'm trying to pick the active route,
	// then jump 2 lines and get the row
	// Not using regex because output is quite standard from Windows XP to 8 (NEEDS TESTING)
	lines := strings.Split(string(output), "\n")
	sep := 0
	for idx, line := range lines {
		if sep == 3 {
			// We just entered the 2nd section containing "Active Routes:"
			if len(lines) <= idx+2 {
				return windowsRouteStructIPv4{}, errNoGateway
			}

			fields := strings.Fields(lines[idx+2])
			if len(fields) < 5 {
				return windowsRouteStructIPv4{}, errCantParse
			}

			return windowsRouteStructIPv4{
				Destination: fields[0],
				Netmask:     fields[1],
				Gateway:     fields[2],
				Interface:   fields[3],
				Metric:      fields[4],
			}, nil
		}
		if strings.HasPrefix(line, "=======") {
			sep++
			continue
		}
	}
	return windowsRouteStructIPv4{}, errNoGateway
}

func parseToWindowsRouteStructIPv6(output []byte) (windowsRouteStructIPv6, error) {

	lines := strings.Split(string(output), "\n")
	sep := 0
	for idx, line := range lines {
		if sep == 3 {
			// We just entered the 2nd section containing "Active Routes:"
			if len(lines) <= idx+2 {
				return windowsRouteStructIPv6{}, errNoGateway
			}

			fields := strings.Fields(lines[idx+2])
			if len(fields) < 4 {
				return windowsRouteStructIPv6{}, errCantParse
			}

			return windowsRouteStructIPv6{
				If:          fields[0],
				Metric:      fields[1],
				Destination: fields[2],
				Gateway:     fields[3],
			}, nil
		}
		if strings.HasPrefix(line, "=======") {
			sep++
			continue
		}
	}
	return windowsRouteStructIPv6{}, errNoGateway
}

func parseWindowsGatewayIPv4(output []byte) (net.IP, error) {
	parsedOutput, err := parseToWindowsRouteStructIPv4(output)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(parsedOutput.Gateway)
	if ip == nil {
		return nil, errCantParse
	}
	return ip, nil
}

func parseWindowsGatewayIPv6(output []byte) (net.IP, error) {
	parsedOutput, err := parseToWindowsRouteStructIPv6(output)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(parsedOutput.Gateway)
	if ip == nil {
		return nil, errCantParse
	}
	return ip, nil
}
