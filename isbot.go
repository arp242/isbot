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
	NoBotKnown       uint8 = iota // Known to not be a bot.
	NoBotNoMatch                  // None of the rules matches, so probably not a bot.
	BotPrefetch                   // Prefetch algorithm
	BotLink                       // User-Agent contained an URL.
	BotClientLibrary              // Known client library.
	BotKnownBot                   // Known bot.
	BotBoty                       // User-Agent string looks "boty".
)

// Is this constant a bot?
func Is(r uint8) bool    { return r != NoBotKnown && r != NoBotNoMatch }
func IsNot(r uint8) bool { return !Is(r) }

// Bot checks if this HTTP request looks like a bot.
//
// It returns one of the constants as the reason we think this is a bot.
func Bot(r *http.Request) uint8 {
	h := r.Header

	// Prefetch algorithm.
	//
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Link_prefetching_FAQ
	if h.Get("X-Moz") == "prefetch" || h.Get("X-Purpose") == "prefetch" || h.Get("Purpose") == "prefetch" || h.Get("X-Purpose") == "preview" || h.Get("Purpose") == "preview" {
		return BotPrefetch
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

	for i := range knownBrowsers {
		if strings.Contains(ua, knownBrowsers[i]) {
			return NoBotKnown
		}
	}

	// Something with a link is almost always a bot.
	if strings.Contains(ua, "http://") || strings.Contains(ua, "https://") {
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
		return BotBoty
	}
	return NoBotNoMatch
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

var knownBrowsers = []string{
	"CUBOT_",
	"CUBOT ",
	"NAVER(inapp",
	"SearchCraft/",
	"StudoBrowser/",
	"YaSearchBrowser/",
	"YandexSearch/",
	"YandexSearchBrowser/",
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
	"PageAnalyzer/",
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
