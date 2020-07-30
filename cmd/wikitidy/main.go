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

func main() {
	flag.Parse()

	if *deepclean {
		// TODO(rjk): Finish this.
		log.Println("do deepclean")
		log.Println("warning: feature incomplete")

		if err := filepath.Walk(config.Basepath, article.UpdateMetadata); err != nil {
			log.Fatal("deepclean walk: ", err)
		}
		return 
	}

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
