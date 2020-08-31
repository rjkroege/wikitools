package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/rjkroege/wikitools/wiki"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var (
	// true for now for testing. But put it back.
	testlog = flag.Bool("testlog", false,
		"Log in the conventional way for running in a terminal. Also changes where to find the configuration file.")
)

func main() {
	flag.Parse()
	if !*testlog {
		defer wiki.LogToTemp()()
	}

	// for testing
	log.SetOutput(os.Stderr)

	log.Println("foo bar")

	// 0. read
	mdf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println("can't read Stdin:", err)
		if _, err := io.WriteString(os.Stdout, "NOCUSTOM"); err != nil {
			// We've done our best so give up.
			log.Fatalf("giving up, can't write to Stdout: %v", err)
		}
	}

	// 1. Make a parser.
	// TODO(rjk): what extensions do I need?
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	context := parser.NewContext()

	// 2. Parse
	reader := text.NewReader(mdf)
	docroot := md.Parser().Parse(reader, parser.WithContext(context))

	// experiment
	// TODO(rjk): metadata extraction failed
	// TODO(rjk): parse math data
	// TODO(rjk): extract the outbound links
	metaData := meta.Get(context)
	title := metaData["title"]
	log.Println(title)

	// 3. Update the document and stored state.
	// TODO(rjk): Writing an extension is probably the right way.
	ast.Walk(docroot, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		log.Println("node:", n.Type(), "kind:", n.Kind())
		return ast.WalkContinue, nil
	})

	// 4. Generate HTML
	if err := md.Renderer().Render(os.Stdout, mdf, docroot); err != nil {
		log.Fatalf("can't write: %v", err)
	}
}
