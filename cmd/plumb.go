package cmd

import (
	"log"

	"github.com/rjkroege/wikitools/wiki"
)

func PlumberHelper(settings *wiki.Settings, wikilink string) {
	log.Println("PlumberHelper", wikilink)
}
