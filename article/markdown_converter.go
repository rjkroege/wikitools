package article

import (
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark"
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/rjkroege/wikitools/article/wikiextension"
)


func NewDefaultMarkdownConverter() goldmark.Markdown {
	// TODO(rjk): what extensions do I need?
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			wikiextension.Linkminer,
			mathjax.MathJax,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	return md
}

