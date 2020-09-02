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
	if err := md.Convert(mdf, os.Stdout); err != nil {
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
