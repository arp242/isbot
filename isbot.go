// Package isbot attempts to detect HTTP bots.
//
// A "bot" is defined as any request that isn't a regular browser request
// initiated by the user. This includes things like web crawlers, but also stuff
// like "preview" renderers and the like.
package isbot

import (
	"net/http"
	"strings"
)

const (
	NoBot         uint8 = iota // Not a bot
	Prefetch                   // Prefetch algorithm
	Link                       // User-Agent contained an URL.
	ClientLibrary              // Known client library.
	KnownBot                   // Known bot.
	Boty                       // User-Agent string looks "boty".
)

func Is(r uint8) bool    { return r != NoBot }
func IsNot(r uint8) bool { return r == NoBot }

// Bot checks if this HTTP request looks like a bot.
//
// It returns one of the constants as the reason we think this is a bot.
func Bot(r *http.Request) uint8 {
	h := r.Header

	// Prefetch algorithm.
	//
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Link_prefetching_FAQ
	if h.Get("X-Moz") == "prefetch" || h.Get("X-Purpose") == "prefetch" || h.Get("Purpose") == "prefetch" || h.Get("X-Purpose") == "preview" || h.Get("Purpose") == "preview" {
		return Prefetch
	}

	return UserAgent(r.UserAgent())
}

// UserAgent checks if this User-Agent header looks like a bot.
//
// It returns one of the constants as the reason we think this is a bot.
func UserAgent(ua string) uint8 {
	// TODO: it's not uncommon to not have a User-Agent at all ... not sure what
	// we want to do with that; a quick looks reveals they *may* be regular
	// users who cleared it? Not sure...

	// Something with a link is almost always a bot.
	if strings.Contains(ua, "http://") || strings.Contains(ua, "https://") {
		return Link
	}

	for i := range clientLibraries {
		if strings.Contains(ua, clientLibraries[i]) {
			return ClientLibrary
		}
	}

	for i := range knownBots {
		if strings.Contains(ua, knownBots[i]) {
			return KnownBot
		}
	}

	// Boty words.
	ua = strings.ToLower(ua)
	if strings.Contains(ua, "bot") ||
		strings.Contains(ua, "crawler") ||
		strings.Contains(ua, "spider") ||
		strings.Contains(ua, "spyder") ||
		strings.Contains(ua, "search") ||
		strings.Contains(ua, "worm") ||
		strings.Contains(ua, "fetch") ||
		strings.Contains(ua, "nutch") {
		return Boty
	}
	return NoBot
}

var clientLibraries = []string{
	"Apache-HttpClient/",
	"Go-http-client/",
	"HTTPClient/",
	"Java/",
	"PycURL/",
	"Python-urllib/",
	"Robosourcer/",
	"Ruby",
	"Wget/",
	"Wget/",
	"WinHttp.WinHttpRequest.5",
	"curl/",
	"python-requests/",
}

var knownBots = []string{
	"ADmantX",
	"AlexaToolbar/",
	"BingPreview/",
	"Chrome-Lighthouse",
	"DumpRenderTree/",
	"Faraday v",
	"GigablastOpenSource/",
	"Google Web Preview",
	"Google favicon",
	"Google-Ads-Overview",
	"Google-Site-Verification",
	"GoogleSecurityScanner",
	"Google_Analytics_Snippet_Validator",
	"HeadlessChrome/",
	"Netcraft Web Server Survey",
	"NetcraftSurveyAgent/",
	"Owler/",
	"PageAnalyzer/1.1",
	"ScopeContentAG-HTTP-Client",
	"Survey/",
	"Synapse",
	"Wappalyzer",
	"WhatWeb/",
	"WinInet",
	"WordPress.com",
	"burpcollaborator.net/",
	"okhttp/",
	"panscient.com",
	"tracemyfile/",
	"wsr-agent/",
}
