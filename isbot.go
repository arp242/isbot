// Package isbot attempts to detect HTTP bots.
//
// A "bot" is defined as any request that isn't a regular browser request
// initiated by the user. This includes things like web crawlers, but also stuff
// like "preview" renderers and the like.
package isbot

import (
	"net/http"
	"strconv"
)

type Result uint8

func (r Result) String() string {
	return strconv.Itoa(int(r)) + ": " + map[Result]string{
		0:   "NoBotKnown",
		1:   "NoBotNoMatch",
		2:   "BotPrefetch",
		3:   "BotLink",
		4:   "BotClientLibrary",
		5:   "BotKnownBot",
		6:   "BotBoty",
		7:   "BotShort",
		8:   "BotRangeAWS",
		9:   "BotRangeDigitalOcean",
		10:  "BotRangeServersCom",
		11:  "BotRangeGoogleCloud",
		12:  "BotRangeHetzner",
		13:  "BotRangeAzure",
		14:  "BotRangeAlibaba",
		15:  "BotRangeLinode",
		150: "BotJSPhanton",
		151: "BotJSNightmare",
		152: "BotJSSelenium",
		153: "BotJSWebDriver",
	}[r]
}

// Not bots.
const (
	NoBotKnown   = 0 // Known to not be a bot.
	NoBotNoMatch = 1 // None of the rules matches, so probably not a bot.
)

// Bots identified by User-Agent.
const (
	BotPrefetch      = 2 // Prefetch algorithm
	BotLink          = 3 // User-Agent contained an URL.
	BotClientLibrary = 4 // Known client library.
	BotKnownBot      = 5 // Known bot.
	BotBoty          = 6 // User-Agent string looks "boty".
	BotShort         = 7 // User-Agent is short of strangely formatted.
)

// Bots identified by IP.
const (
	BotRangeAWS          = 8  // AWS cloud
	BotRangeDigitalOcean = 9  // Digital Ocean
	BotRangeServersCom   = 10 // servers.com
	BotRangeGoogleCloud  = 11 // Google Cloud
	BotRangeHetzner      = 12 // hetzner.de
	BotRangeAzure        = 13 // Azure Cloud
	BotRangeAlibaba      = 14 // Alibaba cloud
	BotRangeLinode       = 15 // Linode
)

// These are never set by isbot, but can be used to send signals from JS; for
// example:
//
//	var is_bot = function() {
//	    var w = window, d = document
//	    if (w.callPhantom || w._phantom || w.phantom)
//	        return 150
//	    if (w.__nightmare)
//	        return 151
//	    if (d.__selenium_unwrapped || d.__webdriver_evaluate || d.__driver_evaluate)
//	        return 152
//	    if (navigator.webdriver)
//	        return 153
//	    return 0
//	}
const (
	BotJSPhanton   = 150 // Phantom headless browser.
	BotJSNightmare = 151 // Nightmare headless browser.
	BotJSSelenium  = 152 // Selenium headless browser.
	BotJSWebDriver = 153 // Generic WebDriver-based headless browser.
)

// Is this constant a bot?
func Is(r Result) bool { return r != NoBotKnown && r != NoBotNoMatch }

// IsNot is the inverse of Is().
func IsNot(r Result) bool { return !Is(r) }

// IsUserAgent reports if this is considered a bot because of the User-Agent
// header.
func IsUserAgent(r Result) bool {
	return r == BotLink || r == BotClientLibrary || r == BotKnownBot || r == BotBoty || r == BotShort
}

// Bot checks if this HTTP request looks like a bot.
//
// It returns one of the constants as the reason we think this is a bot.
//
// This assumes that r.RemoteAddr is set to the real IP and does not check
// X-Forwarded-For or X-Real-IP.
//
// Note that both 0 and 1 may indicate that it's *not* a bot; use Is() and
// IsNot() to check.
func Bot(r *http.Request) Result {
	if Prefetch(r.Header) {
		return BotPrefetch
	}

	i := UserAgent(r.UserAgent())
	if i > 1 {
		return i
	}

	return IPRange(r.RemoteAddr)
}

// Prefetch checks if this request is a browser "pre-fetch" request.
//
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Link_prefetching_FAQ
func Prefetch(h http.Header) bool {
	return h.Get("X-Moz") == "prefetch" || h.Get("X-Purpose") == "prefetch" ||
		h.Get("Purpose") == "prefetch" || h.Get("X-Purpose") == "preview" ||
		h.Get("Purpose") == "preview"
}
