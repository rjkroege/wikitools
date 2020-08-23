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


	// Default function now will relocate all files in the non-special directories.
	// TODO(rjk): Consider better command line structure. Surely there's a package
	// to do this for me.
	// This is just too dangerous to be the default action. :-)

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


	// Enumerate all of the articles in config.Newarticlespath
	articlepattern := filepath.Join(config.Basepath, config.Newarticlespath, "*")

	files, err := filepath.Glob(articlepattern)
	if err != nil {
		log.Println("Can't enumerate the available files because: ", err)
		os.Exit(1)
	}
	
	for _, f := range files {
		// TODO(rjk): Per README, handle supplemental directories.
		if filepath.Ext(f) != ".md" {
			log.Printf("Skipping %s because not Markdown", f)
			continue
		}

		fileinfo, err := os.Stat(f)
		if err != nil {
			log.Printf("Skipping %s because can't stat: %v", f, err)
			continue
		}

		// Build an article metadata representation for this article.
		article := article.MakeMetaData(filepath.Base(f), fileinfo.ModTime())
		fd, err := os.Open(f)
		if err != nil {
			log.Printf("Skipping %s because: %v", f, err)
			continue
		}
		article.RootThroughFileForMetadata(fd)
		fd.Close()

		relativearticledirectory := article.RelativeDateDirectory()
		absarticledirectory := filepath.Join(config.Basepath, relativearticledirectory)
		destarticle := filepath.Join(absarticledirectory, article.FileName())
		srcarticle := f

		if *dryrun {
			log.Printf("%s -> %s\n", srcarticle, destarticle)
			continue
		}

		if err := os.MkdirAll(absarticledirectory, 0700); err != nil {
			log.Printf("Skipping creating dir %s because: %v", absarticledirectory, err)
			continue
		}

		if err := os.Link(srcarticle, destarticle); err != nil {
			log.Printf("Can't link %s to %s because %v", srcarticle, destarticle, err)
			continue
		}

		if err := os.RemoveAll(srcarticle); err != nil {
			log.Printf("Can't remove %s because %v", srcarticle, err)
		}
	}
}
