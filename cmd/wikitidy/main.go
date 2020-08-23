package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/config"
)

var dryrun = flag.Bool("n", false, "Don't actually move the files, just show what would happen")
var deepclean = flag.Bool("deepclean", false, "Rewrite the metadata, move files into improved directories")
var reportflag = flag.Bool("report", false, "Generate the metadata status report.")

func main() {
	flag.Parse()

	// TODO(rjk): This mainline can be refactored nicely..
	if *reportflag {
		log.Println("do report")
		report, err := article.MakeMetadataReporter()
		if err != nil {
			log.Fatal("No BatchCleaner:", err)
		}

		if err := filepath.Walk(config.Basepath, func(path string, info os.FileInfo, err error) error {
			return report.EachFile(path, info, err)
		}); err != nil {
			log.Fatal("report walk: ", err)
		}

		if err := report.Summary(); err != nil {
			log.Fatal("report Summary: ", err)
		}
		os.Exit(0)
	}

	if *deepclean {
		// TODO(rjk): Finish this.
		log.Println("do deepclean")
		update, err := article.MakeMetadataUpdater()
		if err != nil {
			log.Fatal("No MetadataUpdater:", err)
		}

		if err := filepath.Walk(config.Basepath, func(path string, info os.FileInfo, err error) error {
			return update.EachFile(path, info, err)
		}); err != nil {
			log.Fatal("deepclean walk: ", err)
		}
		if err := update.Summary(); err != nil {
			log.Fatal("report Summary: ", err)
		}
		os.Exit(0)
	}

	// Default function now will relocate all files in the non-special
	// directories. TODO(rjk): Consider better command line structure. Surely
	// there's a package to do this for me.

	mover, err := article.MakeFilemover(*dryrun)
	if err != nil {
		log.Fatal("No MakeFilemover:", err)
	}
	if err := filepath.Walk(config.Basepath, func(path string, info os.FileInfo, err error) error {
		return mover.EachFile(path, info, err)
	}); err != nil {
		log.Fatal("mover walk: ", err)
	}
	if err := mover.Summary(); err != nil {
		log.Fatal("report Summary: ", err)
	}
	os.Exit(0)

}
