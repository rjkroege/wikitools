package cmd

import (
	"log"

	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
)


// I get here on main. But I need to actually run the the response in a
// callback. On Darwin, we 


func PlumberHelper(settings *wiki.Settings, wikilink string) {
	log.Println("PlumberHelper", wikilink)

	// The spotlight indexer is going to make a new thread that runs the Mac world.
	indexer := search.MakeWikilinkNameIndex()

	// TODO(rjk): Trial code.
	paths, err := indexer.Allpaths(wikilink)
	if err != nil {
		log.Printf("indexer.Allpaths errored on %q: %v", wikilink, err)
	}
	log.Println(paths)
}
