package isbot

import (
	"net/http"
	"strings"
	"testing"
)

func BenchmarkBot(b *testing.B) {
	r := &http.Request{Header: make(http.Header)}
	r.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:71.0) Gecko/20100101 Firefox/71.0")

	// Comparison of strings.Contains() vs. the regexp mssola/user_agent uses:
	// var reBot = regexp.MustCompile("(?i)(bot|crawler|sp(i|y)der|search|worm|fetch|nutch)")
	//
	// re:          BenchmarkBot-2             55642             20843 ns/op
	// re no capt:  BenchmarkBot-2             91999             14231 ns/op
	// contains:    BenchmarkBot-2           1826011               653 ns/op
	for n := 0; n < b.N; n++ {
		Bot(r)
	}
}

func TestBot(t *testing.T) {
	var fail []string
	for _, b := range bots {
		r := &http.Request{Header: make(http.Header)}
		r.Header.Add("User-Agent", b)
		if IsNot(Bot(r)) {
			fail = append(fail, b)
		}
	}
	if len(fail) > 0 {
		t.Errorf("failed:\n%s", strings.Join(fail, "\n"))
	}
}

func TestNotBot(t *testing.T) {
	var fail []string
	for _, b := range notBots {
		r := &http.Request{Header: make(http.Header)}
		r.Header.Add("User-Agent", b)
		if Is(Bot(r)) {
			fail = append(fail, b)
		}
	}
	if len(fail) > 0 {
		t.Errorf("failed:\n%s", strings.Join(fail, "\n"))
	}
}

func TestDup(t *testing.T) {
	for _, list := range [][]string{bots, notBots} {
		for i := range list {
			for j := range list {
				if i != j && list[i] == list[j] {
					t.Errorf("duplicate: %s", list[i])
				}
			}
		}
	}
}

var bots = []string{
	// Baidu bot
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible; Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 (compatible; Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)",
	"Baiduspider+(+http://www.baidu.com/search/spider.htm???);googlebot|baiduspider|baidu|spider|sogou|bingbot|bot|yahoo|soso|sosospider|360spider|youdaobot",
	"Mozilla/5.0+(compatible;+Baiduspider/2.0;++http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible%3B Baiduspider-render/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible; Baiduspider/2.0;+http://www.baidu.com/search/spider.html）",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36 Mozilla/5.0 (Linux;u;Android 4.2.2;zh-cn;) AppleWebKit/534.46 (KHTML,like Gecko) Version/5.1 Mobile Safari/10600.6.3 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",

	// BingBot
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534+ (KHTML, like Gecko) BingPreview/1.0b",
	"Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)",
	"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm) Safari/537.36",

	// GoogleBot
	"GoogleBot 1.0",
	"Googlebot/2.1 (+http://www.google.com/bot.html)",
	"Googlebot/2.1 (+http://www.googlebot.com/bot.html)",
	"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36 Googlebot/2.1 (+http://www.googlebot.com/bot.html)",
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html",
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 6_0 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Mobile/10A5376e Safari/8536.25 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_1 like Mac OS X; en-us) AppleWebKit/532.9 (KHTML, like Gecko) Version/4.0.5 Mobile/8B117 Safari/6531.22.7 (compatible; Googlebot-Mobile/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Chrome/80.0.3987.122 Safari/537.36",
	"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko; compatible; Googlebot/2.1; +http://www.google.com/bot.html) Safari/537.36",
	"Mozilla/5.0 Googlebot/2.1",

	// Other Google stuff
	"APIs-Google (+https://developers.google.com/webmasters/APIs-Google.html)",
	"AdsBot-Google (+http://www.google.com/adsbot.html)",
	"Google-Ads-Overview Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.118 Safari/537.36",
	"Google_Analytics_Snippet_Validator",
	"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2272.96 Mobile Safari/537.36 (compatible; Google-Shopping-Quality +http://www.google.com/merchants/tos/extend/US/tos.html)",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.9; rv:28.0) Gecko/20100101 Firefox/28.0 (FlipboardProxy/1.6; +http://flipboard.com/browserproxy)",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.116 Safari/537.36 AppEngine-Google; (+http://code.google.com/appengine; appid: s~wanglion204)",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0 AppEngine-Google; (+http://code.google.com/appengine; appid: s~cnwogodl4)",
	"Mozilla/5.0 (Windows NT 6.1; rv:6.0) Gecko/20110814 Firefox/6.0 Google (+https://developers.google.com/+/web/snippet/)",
	"Mozilla/5.0 (Windows NT 6.1; rv:6.0) Gecko/20110814 Firefox/6.0 Google favicon",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.9.1.2) Gecko/20090729 Firefox/3.5.2 (.NET CLR 3.5.30729; Diffbot/0.1; IpsyCrawler; +http://www.diffbot.com)",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 Google (+https://developers.google.com/+/web/snippet/)",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 AppEngine-Google; (+http://code.google.com/appengine; appid: s~rylan-yan01)",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko; Google Web Preview) Chrome/27.0.1453 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; GoogleSecurityScanner Google-HTTP-Java-Client) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36 (scan 5098682184892416, run 5053249132036096)",
	"Mozilla/5.0 (compatible; Google-Site-Verification/1.0)",
	"Mozilla/5.0 (compatible; Google-Structured-Data-Testing-Tool +https://search.google.com/structured-data/testing-tool)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 (compatible; AdsBot-Google-Mobile; +http://www.google.com/mobile/adsbot.html)",
	"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.92 Safari/537.36 (compatible; Google-Shopping-Quality +http://www.google.com/merchants/tos/extend/US/tos.html)",

	// User-Agents from various client libraries that are rarely a browsers.
	"Apache-HttpClient/4.2.3 (java 1.5)",
	"Apache-HttpClient/4.3 (java 1.5)",
	"Apache-HttpClient/4.3.3 (java 1.5)",
	"Go-http-client/1.1",
	"HTTPClient/1.0 (2.3.4.1, ruby 1.9.3 (2013-06-27))",
	"HTTPClient/1.0 (2.4.0, ruby 1.9.3 (2013-06-27))",
	"Java/1.6.0_29",
	"Java/1.6.0_45",
	"Java/1.7.0_09",
	"Java/1.7.0_21",
	"Java/1.7.0_40",
	"Java/1.7.0_60-ea",
	"Java/1.7.0_65",
	"Mozilla/4.0 (compatible; Win32; WinHttp.WinHttpRequest.5)",
	"PycURL/7.23.1",
	"Python-urllib/1.17",
	"Python-urllib/2.6",
	"Python-urllib/2.7",
	"Python-urllib/3.4",
	"Robosourcer/1.0",
	"Ruby",
	"Wget/1.12 (linux-gnu)",
	"Wget/1.13.4 (linux-gnu)",
	"curl/7.19.7 (x86_64-redhat-linux-gnu) libcurl/7.19.7 NSS/3.13.1.0 zlib/1.2.3 libidn/1.18 libssh2/1.2.2",
	"curl/7.19.7 (x86_64-redhat-linux-gnu) libcurl/7.19.7 NSS/3.14.0.0 zlib/1.2.3 libidn/1.18 libssh2/1.4.2",
	"curl/7.19.7 (x86_64-redhat-linux-gnu) libcurl/7.19.7 NSS/3.15.3 zlib/1.2.3 libidn/1.18 libssh2/1.4.2",
	"curl/7.35.0",
	"python-requests/1.1.0 CPython/2.7.4 Linux/3.8.0-19-generic",
	"python-requests/1.2.0 CPython/2.7.4 Linux/3.8.0-33-generic",
	"python-requests/2.2.1 CPython/2.7.6 Linux/3.13.0-24-generic",

	// Headless Chrome
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/60.0.3095.0 Safari/537.36 SeoSiteCheckup (https://seositecheckup.com)",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/68.0.3440.106 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/69.0.3497.81 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/77.0.3835.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/78.0.3882.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/80.0.3987.132 Safari/537.36 Prerender (+https://github.com/prerender/prerender)",

	// Various analyzer tools one can run on a website.
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3694.0 Safari/537.36 Chrome-Lighthouse",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/61.0.3163.100 Chrome/61.0.3163.100 Safari/537.36 PingdomPageSpeed/1.0 (pingbot/2.0; +http://www.pingdom.com/)",
	"Mozilla/5.0 (compatible; Wappalyzer)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1 SeoSiteCheckup (https://seositecheckup.com)",
	"W3C_Validator/1.3 http://validator.w3.org/services",

	// Arguably a "real" browser, since a person (probably) sees the image; we
	// count it as a bot though.
	"Slack-ImgProxy (+https://api.slack.com/robots)",
	"Slackbot 1.0 (+https://api.slack.com/robots)",
	"Slackbot-LinkExpanding 1.0 (+https://api.slack.com/robots)",

	// Uncategorized
	"ADmantX Platform Semantic Analyzer - ADmantX Inc. - www.admantx.com - support@admantx.com",
	"CATExplorador/1.0beta (sistemes at domini dot cat; http://domini.cat/catexplorador.html)",
	"COMODOSpider/Nutch-1.2",
	"Comodo Spider 1.2",
	"Comodo-Webinspector-Crawler 2.1",
	"DemandbaseSiteAnalyzer/0.1 (http://www.demandbase.com; info@demandbase.com)",
	"Doximity-Diffbot",
	"Facebot",
	"Faraday v0.8.9",
	"GigablastOpenSource/1.0",
	"GumGum-Bot/1.0 (http://gumgum.com; support@gumgum.com)",
	"Http://ar5ecpsuzokh4t5pff0zjpvp7gd91zpsrgm3cr1.burpcollaborator.net/5.0 (Linux; Android 7.0; PLUS Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36",
	"LinkedInBot/1.0 (compatible; Mozilla/5.0; Apache-HttpClient +http://www.linkedin.com)",
	"MAZBot/1.0 (http://mazdigital.com)",
	"MJ12bot/v1.0.8 (http://majestic12.co.uk/bot.php?+)",
	//"Mozilla/2.0 (compatible; crw)",
	//"Mozilla/3.0 (compatible; Indy Library)",
	//"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; InfoPath.2)",
	"Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.1; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; MDDR; .NET4.0C; .NET4.0E; .NET CLR 1.1.4322; Tablet PC 2.0); 360Spider",
	"Mozilla/4.0 (compatible; Netcraft Web Server Survey)",
	"Mozilla/4.0 (compatible; Synapse)",
	"Mozilla/4.0 (compatible; http://search.thunderstone.com/texis/websearch/about.html)",
	"Mozilla/5.0 (Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.100 Safari/537.36 (compatible; Statically-Screenshot; +https://statically.io/screenshot)",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:41.0) Gecko/20100101 Firefox/55.0 BrandVerity/1.0 (http://www.brandverity.com/why-is-brandverity-visiting-me)",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/600.2.5 (KHTML, like Gecko) Version/8.0.2 Safari/600.2.5 (Applebot/0.1; +http://www.apple.com/go/applebot)",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64, SimpleAnalyticsBot/1.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1547.57 Safari/537.36 AlexaToolbar/alxg-3.1",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.89 Safari/537.1; 360Spider",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.89 Safari/537.1; 360Spider(compatible; HaosouSpider; http://www.haosou.com/help/help_3_2.html)",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.79 Safari/537.36 (https://shrinktheweb.com)",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) KomodiaBot/1.0",
	"Mozilla/5.0 (Windows NT 6.2; WOW64) Runet-Research-Crawler (itrack.ru/research/cmsrate; rating@itrack.ru)",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36 (compatible; PagePeeker/3.0; +https://pagepeeker.com/robots/)",
	"Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.120 Safari/537.36 (compatible; PagePeeker/3.0; +https://pagepeeker.com/robots/)",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; en; rv:1.9.0.13) Gecko/2009073022 Firefox/3.5.2 (.NET CLR 3.5.30729) Survey/2.3 (fr.wsdata.com)",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; en; rv:1.9.0.13) Gecko/2009073022 Firefox/3.5.2 (.NET CLR 3.5.30729) SurveyBot/2.3 (DomainTools)",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; )  Firefox/1.5.0.11; 360Spider",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.8.0.11)  Firefox/1.5.0.11; 360Spider",
	"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.8.0.11) Gecko/20070312 Firefox/1.5.0.11; 360Spider",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/536.11 (KHTML, like Gecko) DumpRenderTree/0.0.0.0 Safari/536.11",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko; compatible; BuiltWith/1.0; +http://builtwith.com/biup) Chrome/60.0.3112.50 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:10.0.12) Gecko/20100101 Firefox/21.0 WordPress.com mShots",
	"Mozilla/5.0 (compatible; AhrefsBot/4.0; +http://ahrefs.com/robot/)",
	"Mozilla/5.0 (compatible; AhrefsBot/6.1; +http://ahrefs.com/robot/)",
	"Mozilla/5.0 (compatible; Aprc/2.9.15-24; +https://aprc.it/) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3528.4 Safari/537.36",
	"Mozilla/5.0 (compatible; Cincraw/1.0; +http://cincrawdata.net/bot/)",
	"Mozilla/5.0 (compatible; EchoboxBot/1.0; hash/w4mwnpbXf3MFAbxOkJRw; +http://www.echobox.com)",
	"Mozilla/5.0 (compatible; Embedly/0.2; +http://support.embed.ly/)",
	"Mozilla/5.0 (compatible; Embedly/0.2; snap; +http://support.embed.ly/)",
	"Mozilla/5.0 (compatible; IstellaBot/1.18.81 +http://www.tiscali.it/)",
	"Mozilla/5.0 (compatible; MJ12bot/v1.2.4; http://www.majestic12.co.uk/bot.php?+)",
	"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 5.1) (http://name911.com)",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0); 360Spider",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0); 360Spider(compatible; HaosouSpider; http://www.haosou.com/help/help_3_2.html)",
	"Mozilla/5.0 (compatible; Mappy/1.0; Warning:UserAgent will be changed by Feb 2020; +http://mappydata.net/bot/)",
	"Mozilla/5.0 (compatible; NetcraftSurveyAgent/1.0; +info@netcraft.com)",
	"Mozilla/5.0 (compatible; ONBbot/3.3.0 +https://webarchiv.onb.ac.at/robot.html)",
	"Mozilla/5.0 (compatible; Onespot-ScraperBot/1.0; +https://www.onespot.com/identifying-traffic.html)",
	"Mozilla/5.0 (compatible; Owler/0.4; +; )",
	"Mozilla/5.0 (compatible; PageAnalyzer/1.1;)",
	"Mozilla/5.0 (compatible; Wappalyzer; https://www.wappalyzer.com)",
	"Mozilla/5.0 (compatible; XML Sitemaps Generator; http://www.xml-sitemaps.com) Gecko XML-Sitemaps/1.0",
	"Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)",
	"Mozilla/5.0 (compatible; archive.org_bot +http://www.archive.org/details/archive.org_bot)",
	"Mozilla/5.0 (compatible; oBot/2.3.1; +http://www.xforce-security.com/crawler/)",
	"Mozilla/5.0 (compatible; special_archiver/3.1.1 +http://www.archive.org/details/archive.org_bot)",
	"Mozilla/5.0 (compatible; tracemyfile/1.0)",
	"Mozilla/5.0 (compatible; woorankreview/2.0; +https://www.woorank.com/)",
	"Mozilla/5.0 (compatible;Impact Radius Compliance Bot)",
	"Mozilla/5.0 (compatible;contxbot/1.0)",
	"Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1 (compatible; woorankreview/2.0; +https://www.woorank.com/)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1 (compatible; woorankreview/2.0; +https://www.woorank.com/)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Mobile/11A465 Twitter for iPhone BrandVerity/1.0 (http://www.brandverity.com/why-is-brandverity-visiting-me)",
	"Mozilla/5.0 AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36 (Refindbot/1.0)",
	"Mozilla/5.0 PhantomJS (compatible; Seznam screenshot-generator 2.1; +http://fulltext.sblog.cz/screenshot/)",
	"Mozilla/5.0(compatible; Sosospider/2.0; +http://help.soso.com/webspider.htm)",
	"Mozilla/5.0(compatible;Sosospider/2.0;+http://help.soso.com/webspider.htm)",
	"NativeAIBot",
	"NetNewsWire (RSS Reader; https://ranchero.com/netnewswire/)",
	"NutchCVS/0.8-dev (Nutch; http://lucene.apache.org/nutch/bot.html; nutch-agent@lucene.apache.org)",
	"Porkbun/Mustache (Website Analysis; http://porkbun.com; tech@porkbun.com)",
	"Quora-Bot",
	"ScopeContentAG-HTTP-Client www.thescope.com/0.1",
	"Sosospider+(+http://help.soso.com/webspider.htm)",
	"TelegramBot (like TwitterBot)",
	"Twitterbot",
	"WebTarantula.com Crawler",
	"WhatWeb/0.4.8-dev",
	"Who.is Bot",
	"WinInet Test",
	"WordPress/4.9.11; http://119.3.12.182/wordpress",
	"WordPress/4.9.12; http://119.3.12.182/wordpress",
	"WordPress/5.3.2; https://radiosonar.net",
	"YisouSpider",
	"arquivo-web-crawler (compatible; heritrix/3.3.0-SNAPSHOT-2019-08-26T10:34:48Z +http://arquivo.pt)",
	"bitlybot/4.0 (+http://bit.ly/)",
	"facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)",
	"facebookexternalhit/1.1 (compatible; Blueno/1.0; +http://naver.me/scrap)",
	"facebookplatform/1.0 (+http://developers.facebook.com)",
	"http://w850tb9gga13lfmbw1hl0bcbo2uvil6bu6hw5l.burpcollaborator.net/5.0 (Linux; Android 7.0; PLUS Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36",
	"https://m9wqu1a6h02tm5n1xrib11d1psvljb72vxin6c.burpcollaborator.net/5.0 (Linux; Android 7.0; PLUS Build/NRD90M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.98 Mobile Safari/537.36",
	"ip-web-crawler.com",
	"lotayabot",
	"ltx71 - (http://ltx71.com/)",
	"lyticsbot-external",
	"msnbot/2.0b (+http://search.msn.com/msnbot.htm)",
	"okhttp/3.12.1",
	"panscient.com",
	"screeenly-bot 2.0",
	"spotinfluence/Nutch-1.4 (Spot Influence crawler; http://spotinfluence.com; hello at spotinfluence dot com)",
	"visaduhoc.info Crawler",
	"wsr-agent/1.0",

	// TODO: determine what this is
	//"Mozilla/5.0 (compatible; MSIE 10.0; Windows NT 6.2; Trident/6.0; dbot)",
	//"Instapaper/7.7.1.2 CFNetwork/978.0.7 Darwin/18.7.0",

	// TODO: not sure if this is a bot?
	// https://github.com/tesseract-ocr/tesseract
	//471 | Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X; Tesseract/1.0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.4 Mobile/15E148 Safari/604.1
	//468 | Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X; Tesseract/1.0) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1
}

var notBots = []string{
	// Popular
	"Mozilla/5.0 (Android 9; Mobile; rv:68.0) Gecko/68.0 Firefox/68.0",
	"Mozilla/5.0 (Linux; Android 4.2.1; en-us; Nexus 5 Build/JOP40D) AppleWebKit/535.19 (KHTML, like Gecko; googleweblight) Chrome/38.0.1025.166 Mobile Safari/535.19",
	"Mozilla/5.0 (Linux; Android 6.0.1; SAMSUNG SM-G800H) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/11.1 Chrome/75.0.3770.143 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 9; SAMSUNG SM-A730F) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/11.1 Chrome/75.0.3770.143 Mobile Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.132 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:67.0) Gecko/20100101 Firefox/67.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:72.0) Gecko/20100101 Firefox/72.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:71.0) Gecko/20100101 Firefox/71.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 12_4_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.2 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.4 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.5 Mobile/15E148 Safari/604.1",

	// User-Agents with a + in them
	"MT6735_TD/V1 Linux/3.18.19+ Android/6.0 Release/03.03.2015 Browser/AppleWebKit537.36 Chrome/39.0.0.0 Mobile Safari/537.36 System/Android 6.0;",
	"Mozilla/5.0 (BB10; Kbd) AppleWebKit/537.35+ (KHTML, like Gecko) Version/10.3.3.3216 Mobile Safari/537.35+",
	"Mozilla/5.0 (Linux; Android 4.4.2; V3+ Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/30.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 7.0; Micromax Q402+ Build/NRD90M; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/64.0.3282.137 Mobile Safari/537.36",
	"Mozilla/5.0 (X11; Linux i686 on x86_64; rv:45.0) Gecko/20100101 Firefox/45.0 Ordissimo/3.8.6.6+svn37147",
	"Mozilla/5.0 (X11; Linux i686 on x86_64; rv:52.0) Gecko/20100101 Firefox/52.0 webissimo3/3.8.12+svn35125",
	"Mozilla/5.0 (X11; U; Linux armv7l like Android; en-us) AppleWebKit/531.2+ (KHTML, like Gecko) Version/5.0 Safari/533.2+ Kindle/3.0+",
	"Mozilla/5.066704189 Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",

	// User-Agents with "search" in them but are NOT search crawlers.
	"Mozilla/5.0 (Linux; Android 10; VOG-AL10 Build/HUAWEIVOG-AL10; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.186 Mobile Safari/537.36 SearchCraft/3.6.4 (Baidu; P1 10)",
	"Mozilla/5.0 (Linux; Android 5.1.1; SM-J320F Build/LMY47V; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/55.0.2883.91 Mobile Safari/537.36 YandexSearch/8.05 YandexSearchBrowser/8.05",
	"Mozilla/5.0 (Linux; Android 5.1.1; SM-J320F Build/LMY47V; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.87 Mobile Safari/537.36 YandexSearch/7.15",
	"Mozilla/5.0 (Linux; arm; Android 9; SM-A105F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 YaApp_Android/10.45 YaSearchBrowser/10.45 BroPP/1.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0 Mobile/15E148 Safari/605.1 NAVER(inapp; search; 690; 10.16.4; 8)",

	// Contains a link, but is not a bot.
	"Mozilla/5.0 (Linux; Android 10; Android SDK built for x86 Build/QSR1.190920.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.185 Mobile Safari/537.36 StudoBrowser/3.17.6-android-studo.staging (https://studo.co/studo-browser-information-and-contact)",
	"Mozilla/5.0 (Linux; Android 10; ONEPLUS A6003 Build/QKQ1.190716.003; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/80.0.3987.99 Mobile Safari/537.36 StudoBrowser/3.17.6-android-studo.staging (https://studo.co/studo-browser-information-and-contact)",
	"Mozilla/5.0 (Linux; Android 7.0; HUAWEI VNS-L31 Build/HUAWEIVNS-L31; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/79.0.3945.93 Mobile Safari/537.36 StudoBrowser/3.16.0-android-studo (https://studo.co/studo-browser-information-and-contact)",
	"Mozilla/5.0 (iPad; CPU OS 12_4_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 StudoBrowser/3.16.0-ios-studo (https://studo.co/studo-browser-information-and-contact)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 StudoBrowser/3.16.0-ios-studo.staging (https://studo.co/studo-browser-information-and-contact)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 StudoBrowser/3.16.0-ios-studo (https://studo.co/studo-browser-information-and-contact)",

	// Contains "BOT" but not a bot.
	"Mozilla/5.0 (Linux; Android 6.0; CUBOT MAX) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.92 Mobile Safari/537.36",
	"Mozilla/5.0 (Linux; Android 8.0.0; CUBOT_X18_Plus Build/O00623; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/80.0.3987.132 Mobile Safari/537.36 Instagram 132.0.0.26.134 Android (26/8.0.0; 480dpi; 1080x2016; CUBOT; CUBOT_X18_Plus; CUBOT_X18_Plus; mt6755; en_GB; 202766609)",

	// TODO: fails, and they're detected as a bot; they're quite rare (about 100
	// hits out of 6.5 million) so not a disaster.
	// "Mozilla/5.0 (Linux; Android 6.0; B BOT 550) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.93 Mobile Safari/537.36",
	// "Mozilla/5.0 (Linux; Android 6.0; VR BOT 552) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Mobile Safari/537.36",
	// "Mozilla/5.0 (Linux; Android 7.0; M bot 60) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Mobile Safari/537.36",
	// "Mozilla/5.0 (Linux; Android 7.1.2; M_bot_tab_71 Build/NHG47K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.83 Safari/537.36",
	// "Mozilla/5.0 (Linux; Android 8.1.0; XBot Junior Build/O11019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/65.0.3325.109 Mobile Safari/537.36",

	// Embeded browsers for various services that also commonly run bots.
	"Instagram 128.0.0.19.119 (iPhone12,3; iOS 13_3_1; en_US; en-US; scale=3.00; 1125x2436; 197357527) AppleWebKit/420+",
	"Mozilla/5.0 (Linux; Android 9; LG-H870 Build/PKQ1.190522.001; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/80.0.3987.119 Mobile Safari/537.36 [Pinterest/Android]",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 131.0.0.21.117 (iPhone11,8; iOS 13_3_1; en_US; en-US; scale=2.00; 828x1792; 201058581)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 131.0.0.21.117 (iPhone12,1; iOS 13_3_1; en_US; en-US; scale=2.00; 828x1792; 201058581)",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 LightSpeed [FBAN/MessengerLiteForiOS;FBAV/256.0.1.26.113;FBBV/203261359;FBDV/iPhone8,4;FBMD/iPhone;FBSN/iOS;FBSV/13.3.1;FBSS/2;FBCR/;FBID/phone;FBLC/en_GE;FBOP/0]",

	// Outlook
	"Outlook-iOS/719.3711406.prod.iphone (4.22.0)",
	"Outlook-iOS/719.3750190.prod.iphone (4.23.0)",
}
