package isbot

import "strings"

// UserAgent checks if this User-Agent header looks like a bot.
//
// It returns one of the constants as the reason we think this is a bot.
func UserAgent(ua string) Result {
	// TODO: it's not uncommon to not have a User-Agent at all ... not sure what
	// we want to do with that; a quick looks reveals they *may* be regular
	// users who cleared it? Not sure...

	// Anything without a slash or space is almost certainly a bot.
	// TODO: don't need 2 containsRune/loops over string; copy and modify code.
	if len(ua) < 10 || !strings.ContainsRune(ua, ' ') || !strings.ContainsRune(ua, '/') {
		return BotShort
	}

	// Nothing after first closing )
	// if strings.IndexByte(ua, ')') == len(ua)-1 {
	// 	return BotBoty
	// }

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

	// Boty words.
	// TODO: avoid ToLower() allocation.
	ua = strings.ToLower(ua)
	if strings.Contains(ua, "bot") ||
		strings.Contains(ua, "crawler") ||
		strings.Contains(ua, "spider") {
		return BotBoty
	}

	return NoBotNoMatch
}

var clientLibraries = []string{
	"Go-http-client/",
	"HttpClient/",
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
	"libwww-perl/",
}

var knownBrowsers = []string{
	"CUBOT_",
	"CUBOT ",
	"StudoBrowser/",
}

var knownBots = []string{
	"ADmantX",
	"AlexaToolbar/",
	"BingPreview/",
	"Chrome-Lighthouse",
	"DumpRenderTree/",
	"Faraday v",
	"GigablastOpenSource/",
	// TODO: Just checking for "Google" might actually work; hmm...
	"Google Web Preview",
	"Google favicon",
	"Google-Ad",
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
	"burpcollaborator.net/", // Burp security analyzer
	"okhttp/",
	"panscient.com",
	"tracemyfile/",
	"wsr-agent/",
	"RuxitRecorder/",     // Dynatrace performance monitor.
	"RuxitSynthetic/",    // Dynatrace performance monitor.
	"TrendsmapResolver/", // ?
	"ubermetrics-technologies.com",
	"zgrab/",                     //  https://github.com/zmap/zgrab2/search?q=user-agent
	"nbertaupete95(at)gmail.com", // Not sure what this belongs to
	"Dataprovider.com",
	"wkhtmltoimage", "wkhtmltopdf",
	"SlimerJS",
}
