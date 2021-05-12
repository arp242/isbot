// Command isbot checks if a User-Agent is a bot.
package main

import (
	"fmt"
	"os"

	"zgo.at/isbot"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: %s user-agent [user-agent...]\n", os.Args[0])
		os.Exit(1)
	}

	for _, ua := range os.Args[1:] {
		b := isbot.UserAgent(ua)
		is := isbot.Is(b)
		fmt.Printf("%t %s(%s) ← %s\n", is, map[bool]string{true: " ", false: ""}[is], b, ua)
	}
}
