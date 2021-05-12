Go library to detect bots based on the HTTP request. A "bot" is defined as any
request that isn't a regular browser request initiated by the user. This
includes things like web crawlers, but also stuff like "preview" renderers and
the like.

`Bot()` accepts a `http.Request` since it looks at *all* information, not just
the `User-Agent`. You can use `UserAgent()` if you just have a `User-Agent`, but
it's highly recommended to use `Bot()`.

Import as `zgo.at/isbot`; API docs: https://pkg.go.dev/zgo.at/isbot

There is a command-line tool in `cmd/isbot` to check if User-Agents are bots:

    $ isbot 'Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0' 'Wget/1.13.4 (linux-gnu)'
    false (1: NoBotNoMatch) ← Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0
    true  (4: BotClientLibrary) ← Wget/1.13.4 (linux-gnu)

It's not 100% reliable, and there are some known cases where it gets things
wrong. See [`isbot_test.go`](/isbot_test.go) for a list of test cases.

The performance is pretty good; turns out that running a few `string.Contains()`
is loads faster than a `(bot|crawler|search|...)` regexp.
