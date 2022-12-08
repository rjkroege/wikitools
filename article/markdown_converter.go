package article

import (
	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/rjkroege/wikitools/article/wikiextension"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

func NewDefaultMarkdownConverter(settings *wiki.Settings) goldmark.Markdown {
	// TODO(rjk): what extensions do I need?
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			wikiextension.NewLinkminer(settings),
			mathjax.MathJax,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
	return md
}
