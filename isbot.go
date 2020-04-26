// Package isbot attempts to detect HTTP bots.
//
// A "bot" is defined as any request that isn't a regular browser request
// initiated by the user. This includes things like web crawlers, but also stuff
// like "preview" renderers and the like.
package isbot

import (
	"net"
	"net/http"
	"strings"
)

const (
	NoBotKnown       uint8 = iota // Known to not be a bot.
	NoBotNoMatch                  // None of the rules matches, so probably not a bot.
	BotPrefetch                   // Prefetch algorithm
	BotLink                       // User-Agent contained an URL.
	BotClientLibrary              // Known client library.
	BotKnownBot                   // Known bot.
	BotBoty                       // User-Agent string looks "boty".
	BotShort                      // User-Agent is short of strangely formatted.

	BotRangeAWS          // AWS cloud
	BotRangeDigitalOcean // Digital Ocean
)

// These are never set by isbot, but can be used to send signals from JS.
const (
	BotJSPhanton uint8 = iota + 150
	BotJSNightmare
	BotJSSelenium
	BotJSWebDriver
)

// Is this constant a bot?
func Is(r uint8) bool    { return r != NoBotKnown && r != NoBotNoMatch }
func IsNot(r uint8) bool { return !Is(r) }

// Bot checks if this HTTP request looks like a bot.
//
// It returns one of the constants as the reason we think this is a bot.
//
// Note: this assumes that r.RemoteAddr is set to the real IP, and does not
// check X-Forwarded-For or X-Real-IP.
func Bot(r *http.Request) uint8 {
	h := r.Header

	// Prefetch algorithm.
	//
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Link_prefetching_FAQ
	if h.Get("X-Moz") == "prefetch" || h.Get("X-Purpose") == "prefetch" || h.Get("Purpose") == "prefetch" || h.Get("X-Purpose") == "preview" || h.Get("Purpose") == "preview" {
		return BotPrefetch
	}

	i := IPRange(r.RemoteAddr)
	if i > 0 {
		return i
	}

	return UserAgent(r.UserAgent())
}

// IPRange checks if this IP address is from a range that should normally never
// send browser requests, such as AWS and other cloud providers.
func IPRange(addr string) uint8 {
	if addr == "" {
		return 0
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		return 0
	}

	if ip.To4() != nil { // TODO: can probably do more efficient check.
		if containsIP(awsEC2Ranges4, ip) {
			return BotRangeAWS
		}
	} else {
		if containsIP(awsEC2Ranges6, ip) {
			return BotRangeAWS
		}
	}
	if containsIP(digitalOceanRanges, ip) {
		return BotRangeDigitalOcean
	}

	return 0
}

func containsIP(ranges []*net.IPNet, ip net.IP) bool {
	for i := range ranges {
		if ranges[i].Contains(ip) {
			return true
		}
	}
	return false
}

func parseNet(ip string) *net.IPNet {
	_, n, err := net.ParseCIDR(ip)
	if err != nil {
		panic(err)
	}
	return n
}

// UserAgent checks if this User-Agent header looks like a bot.
//
// It returns one of the constants as the reason we think this is a bot.
func UserAgent(ua string) uint8 {
	// TODO: it's not uncommon to not have a User-Agent at all ... not sure what
	// we want to do with that; a quick looks reveals they *may* be regular
	// users who cleared it? Not sure...

	// Anything without a slash or space is almost certainly a bot.
	// TODO: don't need 2 containsRune/loops over string; copy and modify code.
	if len(ua) < 10 || !strings.ContainsRune(ua, ' ') || !strings.ContainsRune(ua, '/') {
		return BotShort
	}

	for i := range knownBrowsers {
		if strings.Contains(ua, knownBrowsers[i]) {
			return NoBotKnown
		}
	}

	// Something with a link is almost always a bot.
	if strings.Contains(ua, "://") {
		return BotLink
	}

	for i := range clientLibraries {
		if strings.Contains(ua, clientLibraries[i]) {
			return BotClientLibrary
		}
	}

	for i := range knownBots {
		if strings.Contains(ua, knownBots[i]) {
			return BotKnownBot
		}
	}

	// TODO: avoid ToLower() allocation.
	// BenchmarkBot-2            724268              1574 ns/op               0 B/op          0 allocs/op
	// BenchmarkBot-2            541042              1996 ns/op              80 B/op          1 allocs/op

	// Boty words.
	ua = strings.ToLower(ua)
	if strings.Contains(ua, "bot") ||
		strings.Contains(ua, "crawler") ||
		strings.Contains(ua, "spider") {
		return BotBoty
	}
	return NoBotNoMatch
}
