//go:generate ./mkip

package isbot

import (
	"net/netip"
)

type ipRange struct {
	bot    Result
	prefix netip.Prefix
}

func parseNet(ip string, b Result) ipRange {
	prefix, err := netip.ParsePrefix(ip)
	if err != nil {
		panic(err)
	}
	return ipRange{bot: b, prefix: prefix}
}

// IPRange checks if this IP address is from a range that should normally never
// send browser requests, such as AWS and other cloud providers.
func IPRange(addr string) Result {
	if addr == "" {
		return NoBotKnown
	}
	ip, err := netip.ParseAddr(addr)
	if err != nil {
		return NoBotKnown
	}

	for _, r := range ipRanges {
		if r.prefix.Contains(ip) {
			return r.bot
		}
	}
	return NoBotNoMatch
}
