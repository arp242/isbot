// Package isbot attempts to detect HTTP bots.
//
// A "bot" is defined as any request that isn't a regular browser request
// initiated by the user. This includes things like web crawlers, but also stuff
// like "preview" renderers and the like.
package isbot

import (
	"net/http"
)

const (
	NoBotKnown           = 0  // Known to not be a bot.
	NoBotNoMatch         = 1  // None of the rules matches, so probably not a bot.
	BotPrefetch          = 2  // Prefetch algorithm
	BotLink              = 3  // User-Agent contained an URL.
	BotClientLibrary     = 4  // Known client library.
	BotKnownBot          = 5  // Known bot.
	BotBoty              = 6  // User-Agent string looks "boty".
	BotShort             = 7  // User-Agent is short of strangely formatted.
	BotRangeAWS          = 8  // AWS cloud
	BotRangeDigitalOcean = 9  // Digital Ocean
	BotRangeServersCom   = 10 // servers.com
	BotRangeGoogleCloud  = 11 // Google Cloud
	BotRangeHetzner      = 12 // hetzner.de
)

// These are never set by isbot, but can be used to send signals from JS; for
// example:
//
//    var is_bot = function() {
//        var w = window, d = document
//        if (w.callPhantom || w._phantom || w.phantom)
//            return 150
//        if (w.__nightmare)
//            return 151
//        if (d.__selenium_unwrapped || d.__webdriver_evaluate || d.__driver_evaluate)
//            return 152
//        if (navigator.webdriver)
//            return 153
//        return 0
//    }
const (
	BotJSPhanton   = 150 // Phantom headless browser.
	BotJSNightmare = 151 // Nightmare headless browser.
	BotJSSelenium  = 152 // Selenium headless browser.
	BotJSWebDriver = 153 // Generic WebDriver-based headless browser.
)

// Is this constant a bot?
func Is(r uint8) bool { return r != NoBotKnown && r != NoBotNoMatch }

// IsNot is the inverse of Is().
func IsNot(r uint8) bool { return !Is(r) }

// Bot checks if this HTTP request looks like a bot.
//
// It returns one of the constants as the reason we think this is a bot.
//
// Note: this assumes that r.RemoteAddr is set to the real IP, and does not
// check X-Forwarded-For or X-Real-IP.
func Bot(r *http.Request) uint8 {
	if Prefetch(r.Header) {
		return BotPrefetch
	}

	i := IPRange(r.RemoteAddr)
	if i > 0 {
		return i
	}

	return UserAgent(r.UserAgent())
}

// Prefetch checks if this request is a browser "pre-fetch" request.
//
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Link_prefetching_FAQ
func Prefetch(h http.Header) bool {
	return h.Get("X-Moz") == "prefetch" || h.Get("X-Purpose") == "prefetch" ||
		h.Get("Purpose") == "prefetch" || h.Get("X-Purpose") == "preview" ||
		h.Get("Purpose") == "preview"
}
