//go:generate go run ./cmd/iprange

package isbot

import (
	"net/netip"
	"strings"
)

type ipRange struct {
	bot    Result
	prefix netip.Prefix
}

func botname(n string) Result {
	switch n {
	case "AWS":
		return BotRangeAWS
	case "DigitalOcean":
		return BotRangeDigitalOcean
	case "ServersCom":
		return BotRangeServersCom
	case "GoogleCloud":
		return BotRangeGoogleCloud
	case "Hetzner":
		return BotRangeHetzner
	case "Azure":
		return BotRangeAzure
	case "Alibaba":
		return BotRangeAlibaba
	case "Linode":
		return BotRangeLinode
	case "Oracle":
		return BotRangeOracle
	case "OVH":
		return BotRangeOVH
	}
	panic(n)
}

var ipRanges4 = func() map[byte][]ipRange {
	m := make(map[byte][]ipRange)
	for f := range strings.FieldsSeq(ranges4) {
		ip, name, ok := strings.Cut(f, ",")
		if !ok {
			panic(f)
		}
		prefix := netip.MustParsePrefix(ip)
		k := prefix.Addr().As4()[0]
		m[k] = append(m[k], ipRange{prefix: prefix, bot: botname(name)})
	}
	return m
}()

var ipRanges6 = func() map[[2]byte][]ipRange {
	m := make(map[[2]byte][]ipRange)
	for f := range strings.FieldsSeq(ranges6) {
		ip, name, ok := strings.Cut(f, ",")
		if !ok {
			panic(f)
		}
		prefix := netip.MustParsePrefix(ip)
		as := prefix.Addr().As16()
		k := [2]byte{as[0], as[1]}
		m[k] = append(m[k], ipRange{prefix: prefix, bot: botname(name)})
	}
	return m
}()

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

	var ranges []ipRange
	if ip.Is4() {
		ranges = ipRanges4[ip.As4()[0]]
	} else {
		as := ip.As16()
		ranges = ipRanges6[[2]byte{as[0], as[1]}]
	}

	for _, r := range ranges {
		if r.prefix.Contains(ip) {
			return r.bot
		}
	}
	return NoBotNoMatch
}
