package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/wiki"
	"github.com/rjkroege/wikitools/corpus/tidy"
)

var dryrun = flag.Bool("n", false, "Don't actually move the files, just show what would happen")
var deepclean = flag.Bool("deepclean", false, "Rewrite the metadata, move files into improved directories")
var reportflag = flag.Bool("report", false, "Generate the metadata status report.")

func main() {
	flag.Parse()
	
	// The default Tidying implementation can always be created without error.
	tidying, err := tidy.NewFilemover(*dryrun)
	switch {
	case *reportflag:
		tidying, err = tidy.NewMetadataReporter()
		if err != nil {
			log.Fatal("No MetadataReporter:", err)
		}
	case *deepclean:
		tidying, err = tidy.NewMetadataUpdater()
		if err != nil {
			log.Fatal("No MetadataUpdater:", err)
		}
	}

	// Default function now will relocate all files in the non-special
	// directories. TODO(rjk): Consider better command line structure. Surely
	// there's a package to do this for me.

	if err := filepath.Walk(wiki.Basepath, func(path string, info os.FileInfo, err error) error {
		return tidying.EachFile(path, info, err)
	}); err != nil {
		log.Fatal("mover walk: ", err)
	}
	if err := tidying.Summary(); err != nil {
		log.Fatal("report Summary: ", err)
	}
	os.Exit(0)

}
