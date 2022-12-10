package cmd

import (
	"log"
	"os"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/corpus/tidy"
	"github.com/rjkroege/wikitools/wiki"
)

func Tidy(settings *wiki.Settings, dryrun, deepclean, reportflag bool) {
	// The default Tidying implementation can always be created without error.
	// TODO(rjk): Improve the selection of the operations: construct them in
	tidying, err := tidy.NewFilemover(settings, dryrun)
	switch {
	case reportflag:
		tidying, err = tidy.NewMetadataReporter(settings)
		if err != nil {
			log.Fatal("No MetadataReporter:", err)
		}
	case deepclean:
		tidying, err = tidy.NewMetadataUpdater()
		if err != nil {
			log.Fatal("No MetadataUpdater:", err)
		}
	}

	// Default function now will relocate all files in the non-special
	// directories. TODO(rjk): Consider better command line structure. Surely
	// there's a package to do this for me.
	// That's now partially done.

	if err := corpus.Everyfile(settings, tidying); err != nil {
		log.Fatalf("walking all the files: %v", err)
	}

	if err := tidying.Summary(); err != nil {
		log.Fatal("report Summary: ", err)
	}

	// TODO(rjk): Do I need?
	os.Exit(0)
}
