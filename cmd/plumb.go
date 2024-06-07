package cmd

import (
	"log"

	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
)

func PlumberHelper(settings *wiki.Settings, wikilink string) {
	log.Println("PlumberHelper", wikilink)
	indexer := search.MakeWikilinkNameIndex()
	// TODO(rjk): Remove stub code.
	paths, err := indexer.Allpaths(wikilink)
	if err != nil {
		log.Printf("error plumbing %q: %v", wikilink, err)
		return
	}
	log.Println(paths)
}
