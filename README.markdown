Go library to detect bots based on the HTTP request. A "bot" is defined as any
request that isn't a regular browser request initiated by the user. This
includes things like web crawlers, but also stuff like "preview" renderers and
the like.

`Bot()` accepts a `http.Request` since it looks at *all* information, not just
the `User-Agent`. You can use `UserAgent()` if you just have a `User-Agent`, but
it's highly recommended to use `Bot()`.

Import as `zgo.at/isbot`; API docs: https://pkg.go.dev/zgo.at/isbot

It's not 100% reliable, and there are some known cases where it gets things
wrong. See [`isbot_test.go`](/isbot_test.go) for a list of test cases.

The performance is pretty good; turns out that running a few `string.Contains()`
is loads faster than a `(bot|crawler|search|...)` regexp.
