module github.com/rjkroege/wikitools

go 1.24.0

require (
	9fans.net/go v0.0.7
	github.com/alecthomas/kong v0.7.1
	github.com/google/go-cmp v0.6.0
	github.com/litao91/goldmark-mathjax v0.0.0-20191101121019-011def32b12f
	github.com/yuin/goldmark v1.7.1
	go.abhg.dev/goldmark/wikilink v0.5.0
	golang.org/x/sys v0.35.0
)

require golang.org/x/exp v0.0.0-20250911091902-df9299821621

replace 9fans.net/go => github.com/rjkroege/go v0.0.0-20250830183428-45f5717096f4
