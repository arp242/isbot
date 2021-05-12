//go:generate ./mkip

package isbot

import "net"

type ipRange struct {
	bot Result
	net *net.IPNet
}

func parseNet(ip string, b Result) ipRange {
	_, n, err := net.ParseCIDR(ip)
	if err != nil {
		panic(err)
	}
	return ipRange{bot: b, net: n}
}

// IPRange checks if this IP address is from a range that should normally never
// send browser requests, such as AWS and other cloud providers.
func IPRange(addr string) Result {
	if addr == "" {
		return 0
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		return 0
	}

	for _, r := range ipRanges {
		if r.net.Contains(ip) {
			return r.bot
		}
	}
	return 0
}
