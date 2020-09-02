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
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var (
	// Writes the output to a file, logs to stdout for debugging convenience.
	debug = flag.Bool("debug", false,
		"Writes the output to a file, logs to stdout for debugging convenience.")
)

func main() {
	flag.Parse()
	if !*debug {
		defer wiki.LogToTemp()()
	}

	destfd := os.Stdout
	if *debug {
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
		),
//		goldmark.WithParserOptions(
//			parser.WithAutoHeadingID(),
//		),
//		goldmark.WithRendererOptions(
//			html.WithHardWraps(),
//			html.WithXHTML(),
//		),
	)
	context := parser.NewContext()

	// 2. Convert, update shared state, etc.
	if err := md.Convert(mdf, destfd); err != nil {
  		  log.Fatalf("markdown Convert failed: %v", err)
	}

	// experiment
	// TODO(rjk): metadata extraction failed
	// TODO(rjk): parse math data
	// TODO(rjk): extract the outbound links
	metaData := meta.Get(context)
	title := metaData["title"]
	log.Println(title)

	
}
