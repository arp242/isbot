package isbot

import (
	"bufio"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"
)

var bots, notBots []string

func init() {
	if len(bots) == 0 {
		bots = readFile("bots")
	}
	if len(notBots) == 0 {
		notBots = readFile("not_bots")
	}
}

// zgo.at/gadget.Unshorten()
var unshort = strings.NewReplacer(
	"~~", "~",
	"~A", "Android",
	"~c", "Chrome/",
	"~C", "compatible",
	"~e", "Edge/",
	"~f", "Firefox/",
	"~g", "Gecko/",
	"~G", "(KHTML, like Gecko)",
	"~i", "iPhone",
	"~I", "Macintosh",
	"~a", "AppleWebKit/",
	"~L", "Linux",
	"~m", "Mobile/", "~M", "Mobile",
	"~s", "Safari/",
	"~v", "Version/",
	"~W", "Windows",
	"~Z ", "Mozilla/5.0 ")

func readFile(f string) []string {
	fp, err := os.Open("./testdata/" + f)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	i := 0
	var out []string
	for scanner.Scan() {
		i++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}

		out = append(out, unshort.Replace(line))
	}

	return out
}

func BenchmarkBot(b *testing.B) {
	r := &http.Request{Header: make(http.Header)}
	r.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:71.0) Gecko/20100101 Firefox/71.0")

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Bot(r)
	}
}

func TestBotIP(t *testing.T) {
	tests := []struct {
		in   string
		want uint8
	}{
		{"114.122.138.27", NoBotNoMatch},

		{"35.180.1.1", BotRangeAWS},
		{"2600:1fff:5000::1", BotRangeAWS},

		{"100.20.156.62", BotRangeAWS},
		{"13.237.31.191", BotRangeAWS},
		{"13.57.187.238", BotRangeAWS},
		{"13.58.249.187", BotRangeAWS},
		{"18.189.178.53", BotRangeAWS},
		{"18.191.239.50", BotRangeAWS},
		{"18.206.115.23", BotRangeAWS},
		{"18.207.119.101", BotRangeAWS},
		{"18.217.17.227", BotRangeAWS},
		{"18.224.140.0", BotRangeAWS},
		{"18.236.221.165", BotRangeAWS},
		{"3.81.56.221", BotRangeAWS},
		{"3.83.24.166", BotRangeAWS},
		{"3.94.114.22", BotRangeAWS},
		{"34.207.159.142", BotRangeAWS},
		{"34.209.26.42", BotRangeAWS},
		{"34.217.96.163", BotRangeAWS},
		{"34.221.199.187", BotRangeAWS},
		{"34.222.59.41", BotRangeAWS},
		{"34.231.157.157", BotRangeAWS},
		{"34.232.127.140", BotRangeAWS},
		{"35.174.166.183", BotRangeAWS},
		{"44.234.24.80", BotRangeAWS},
		{"44.234.66.18", BotRangeAWS},
		{"52.12.38.56", BotRangeAWS},
		{"52.34.76.65", BotRangeAWS},
		{"52.44.93.197", BotRangeAWS},
		{"52.56.255.25", BotRangeAWS},
		{"54.158.227.15", BotRangeAWS},
		{"54.159.60.243", BotRangeAWS},
		{"54.166.166.23", BotRangeAWS},
		{"54.200.108.160", BotRangeAWS},
		{"54.215.29.10", BotRangeAWS},
		{"54.226.25.34", BotRangeAWS},
		{"54.227.27.249", BotRangeAWS},
		{"54.242.93.252", BotRangeAWS},
		{"54.70.53.60", BotRangeAWS},
		{"54.71.187.124", BotRangeAWS},
		{"54.86.34.110", BotRangeAWS},
		{"54.91.251.150", BotRangeAWS},
		{"54.92.222.34", BotRangeAWS},

		{"68.183.241.134", BotRangeDigitalOcean},

		{"88.212.248.0", BotRangeServersCom},
		{"88.212.255.255", BotRangeServersCom},

		// {"2a01:4f8:162:5447::2", BotRangeHetzner},
		// {"2a01:4f8:140:21ee::2", BotRangeHetzner},

		{"88.213.0.0", NoBotNoMatch},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			r := &http.Request{Header: make(http.Header), RemoteAddr: tt.in}
			r.Header.Add("User-Agent", "Your user agent: Mozilla/5.0 (X11; Linux x86_64; rv:75.0) Gecko/20100101 Firefox/75.0")
			got := Bot(r)
			if got != tt.want {
				t.Errorf("got %d; want %d", got, tt.want)
			}
		})
	}
}

func TestBotUA(t *testing.T) {
	var fail []string
	for _, b := range bots {
		r := &http.Request{Header: make(http.Header)}
		r.Header.Add("User-Agent", b)
		if IsNot(Bot(r)) {
			fail = append(fail, b)
		}
	}
	if len(fail) > 0 {
		t.Errorf("%d failed:\n%s", len(fail), strings.Join(fail, "\n"))
	}
}

func TestNotBotUA(t *testing.T) {
	var fail []string
	for _, b := range notBots {
		r := &http.Request{Header: make(http.Header)}
		r.Header.Add("User-Agent", b)
		if Is(Bot(r)) {
			fail = append(fail, b)
		}
	}
	if len(fail) > 0 {
		t.Errorf("%d failed:\n%s", len(fail), strings.Join(fail, "\n"))
	}
}

func TestDup(t *testing.T) {
	for _, list := range [][]string{bots, notBots} {
		list = list[:]
		sort.Strings(list)
		var last string
		for i := range list {
			if list[i] == last {
				t.Errorf("duplicate: %s", list[i])
			}
			last = list[i]
		}
	}
}
