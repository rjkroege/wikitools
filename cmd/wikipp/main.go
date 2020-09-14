package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/yuin/goldmark/parser"
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

	// -1. Read metadata first and only proceed for valid article types.


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
	md := article.NewDefaultMarkdownConverter()

	// 2. Make a context.
	context := parser.NewContext()
	// TODO(rjk): I'll have to put other state in context for accessing the graph data.
	// In particular, I expect that I'll want a NewDefaultMarkdownConversionContext

	// 3. Convert, update shared state, etc.
	if err := md.Convert(mdf, destfd, parser.WithContext(context)); err != nil {
		log.Fatalf("markdown Convert failed: %v", err)
	}

	// TODO(rjk): Figure out why the metadata extraction is failing.
	// Conversely, I could just use my own metadata extracter because it
	// handles all of my articles with bad metdata (where I should fallback
	// to just copying the content to Stdout and let Marked sort it out.
}
