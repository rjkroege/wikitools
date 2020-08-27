package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rjkroege/wikitools/wiki"
)

var (
	testlog = flag.Bool("testlog", false,
		"Log in the conventional way for running in a terminal. Also changes where to find the configuration file.")
)

func main() {
	flag.Parse()
	if !*testlog {
		defer wiki.LogToTemp()()
	}

	log.Println("foo bar")

	// parse. [do I want to parse or just be hacky with regexps?] If I hack it... things will soon not work.
	// update the database (can happen asynchronously)
	// modify (e.g. custom fenced blocks made into pictures), link maps
	// emit (html or markdown?)

	// pretend to be cat.
	if _, err := io.Copy(os.Stdout, os.Stdin); err != nil {
		log.Println("can't pretend to be cat:", err)
	}

	// Add a string to the doc.
	fmt.Printf("\n---\n# Hello World\n")
}
