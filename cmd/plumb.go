package cmd

import (
	"log"

	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/rjkroege/gozen"
)

func PlumberHelper(settings *wiki.Settings, lsd, wikitext string) {
	log.Println("PlumberHelper", wikitext)

	// Remember that on darwin that we are running in a secondary Go routine.
	// TODO(rjk): Consider renaming this later.
	mapper := search.MakeWikilinkNameIndex()

	// TODO(rjk): Trial code.
	fp, err := mapper.Path(settings.Wikidir, lsd, wikitext)
	if err != nil {
		log.Printf("indexer.Path errored on %q: %v", wikitext, err)
	}

	log.Println(fp)
	gozen.Editinacme(fp)
}
