package netx

import (
	"bytes"
	"net"
	"net/http"
	"strings"
)

// struct ipRange - holds the start and end of a range of ip addresses
type ipRange struct {
	start net.IP
	end   net.IP
}

// Function inRng checks if a given ip address is within a range given
func inRange(r ipRange, ipAddr net.IP) bool {
	// strcmp type byte comparison
	return bytes.Compare(ipAddr, r.start) >= 0 && bytes.Compare(ipAddr, r.end) < 0
}

// List of common private IP ranges
var privateRanges = []ipRange{
	ipRange{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	ipRange{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	ipRange{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	ipRange{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	ipRange{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	ipRange{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// Function isPrivateSubnet checks if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
	// use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func GetClientIPFromRequest(r *http.Request) net.IP {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addrs := strings.Split(r.Header.Get(h), ",")
		// iterate over addresses from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addrs) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addrs[i])
			// header may contain spaces, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return realIP
		}
	}
	// if no ip address specified, assume loopback address (localhost)
	return net.ParseIP("::1")
}
