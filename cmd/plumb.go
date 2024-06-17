package cmd

import (
	"log"

	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
)


// I get here on main. But I need to actually run the the response in a
// callback.

// Proposed behaviour for 


func PlumberHelper(settings *wiki.Settings, wikilink string) {
	log.Println("PlumberHelper", wikilink)

	// Remember that on darwin that we are running in a secondary Go routine.
	indexer := search.MakeWikilinkNameIndex()

	// TODO(rjk): Trial code.
	paths, err := indexer.Allpaths(wikilink)
	if err != nil {
		log.Printf("indexer.Allpaths errored on %q: %v", wikilink, err)
	}
	log.Println(paths)

	// TODO(rjk): Eventually here, we need to
	// a. do the path searching for the result of Allpaths
	// b do the disambiguation for the 1 or more paths.
	
	// TODO(rjk): api needs to change because a wikitext path needs to know
}
