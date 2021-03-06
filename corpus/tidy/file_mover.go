package tidy

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/wiki"
)

type fileMover struct {
	removeddirectories map[string]struct{}
	dryrun             bool
}

// NewFilemover creates a Tidying implementation that positions files in
// the right wiki directories
func NewFilemover(dryrun bool) (corpus.Tidying, error) {
	return &fileMover{
		removeddirectories: make(map[string]struct{}),
		dryrun:             dryrun,
	}, nil
}

func (fm *fileMover) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}

	// TODO(rjk): I could do less work if I returned "skip this directory" for
	// templates and generated.
	if wiki.IsWikiArticle(path, info) {
		return nil
	}

	// TODO(rjk): Make this block into a utility function.
	// need the date
	d, err := os.Stat(path)
	if err != nil {
		log.Println("updateMetadata Stat error", err)
		return fmt.Errorf("can't FileMover Stat %s: %v", path, err)
	}

	// get the metadata
	ifd, err := os.Open(path)
	if err != nil {
		log.Println("updateMetadata Open error", err)
		return fmt.Errorf("can't FileMover Open %s: %v", path, err)
	}
	defer ifd.Close()
	fd := bufio.NewReader(ifd)

	// verify that this is the right path.
	md := article.MakeMetaData(filepath.Base(path), d.ModTime())
	md.RootThroughFileForMetadata(fd)

	// Determine the correct directory for the article.
	relativearticledirectory := md.RelativeDateDirectory()
	absarticledirectory := filepath.Join(wiki.Basepath, relativearticledirectory)
	destarticle := filepath.Join(absarticledirectory, md.FileName())
	srcarticle := path

	if srcarticle == destarticle {
		// nothing to do
		return nil
	}

	if fm.dryrun {
		log.Printf("mv %s -> %s\n", srcarticle, destarticle)
		fm.removeddirectories[filepath.Dir(srcarticle)] = struct{}{}
		return nil
	}

	if err := os.MkdirAll(absarticledirectory, 0700); err != nil {
		return fmt.Errorf("can't mkdir %s because: %v", absarticledirectory, err)
	}

	if err := os.Link(srcarticle, destarticle); err != nil {
		return fmt.Errorf("can't link %s to %s because %v", srcarticle, destarticle, err)
	}

	if err := os.Remove(srcarticle); err != nil {
		return fmt.Errorf("can't remove %s because %v", srcarticle, err)
	}

	// Walk does a pre-order traversal. So we might have removed all of the
	// files in a given directory but we don't know if this directory is
	// empty. But we can record the fact that we've removed something from
	// the directory and clean the directories in the Summary
	fm.removeddirectories[filepath.Dir(srcarticle)] = struct{}{}

	return nil
}

func (fm *fileMover) Summary() error {
	dirs := fm.removeddirectories

	if fm.dryrun {
		log.Println("not removing non-empty directories in dryrun mode")
		return nil
	}

	for workremains := true; workremains; {
		parentdirs := make(map[string]struct{})
		for d := range dirs {
			workremains = false
			if err := os.Remove(d); err == nil {
				if pd := filepath.Dir(d); pd != wiki.Basepath {
					parentdirs[pd] = struct{}{}
					workremains = true
				}
			}
		}
		dirs = parentdirs
	}

	return nil
}
