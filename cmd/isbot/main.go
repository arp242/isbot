// Command isbot checks if a User-Agent is a bot.
package main

import (
	"fmt"
	"os"
	"strings"

	"zgo.at/isbot"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [ user-agent | ip:addr ] [ ... ] \n", os.Args[0])
		os.Exit(1)
	}

	for _, ua := range os.Args[1:] {
		var b isbot.Result
		if strings.HasPrefix(ua, "ip:") {
			b = isbot.IPRange(ua[3:])
		} else {
			b = isbot.UserAgent(ua)
		}
		is := isbot.Is(b)
		fmt.Printf("%t %s(%s) ← %s\n", is, map[bool]string{true: " ", false: ""}[is], b, ua)
	}
}
