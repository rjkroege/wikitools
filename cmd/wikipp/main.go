package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	//   "github.com/alecthomas/chroma/formatters/html"
	//    "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/renderer/html"
)

var (
	// Writes the output to a file, logs to stdout for debugging convenience.
	debug = flag.Bool("debug", false,
		"Writes the output to a file, logs to stdout for debugging convenience.")
)

func main() {
	flag.Parse()

	destfd := os.Stdout
	if !*debug {
		// Marked expects Stdout to be where it will find the generated html. But
		// Goldmark's debugging helpers print to Stdout. So replace Stdout so
		// that I can separate debugging output from the actual output and send
		// that to a log file.
		defer wiki.ReplaceStdout()()
	}
	if *debug {
		// In debug mode, wikipp dumps debugging to Stdout and writes the
		// generated output to out.html
		fd, err := os.Create("out.html")
		if err != nil {
			log.Fatalf("can't write to out.html: %v", err)
		}
		destfd = fd
	}

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

	// TODO(rjk): skip running this on files with bad metadata?

	// 1. Make a converter
	// TODO(rjk): what extensions do I need?
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			meta.Meta,
			Linkminer,
			mathjax.MathJax,
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		//			highlighting.NewHighlighting(
		//               highlighting.WithStyle("monokai"),
		//               highlighting.WithFormatOptions(
		//                   html.WithLineNumbers(true),
		//              ),
		//           ),
	)

	// 2. Make a context.
	context := parser.NewContext()
	// TODO(rjk): I'll have to put other state in context for accessing the graph data.

	// 3. Convert, update shared state, etc.
	if err := md.Convert(mdf, destfd, parser.WithContext(context)); err != nil {
		log.Fatalf("markdown Convert failed: %v", err)
	}

	// TODO(rjk): Figure out why the metadata extraction is failing.
	// Conversely, I could just use my own metadata extracter because it
	// handles all of my articles with bad metdata (where I should fallback
	// to just copying the content to Stdout and let Marked sort it out.
	metaData := meta.Get(context)
	title := metaData["title"]
	log.Println("title from meta", title)

}
