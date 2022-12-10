package corpus


import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/wiki"

)

// Tidying is the interface implemented by each of the kinds of Tidying
// passes.
type Tidying interface {
	// EachFile is called by the filepath.Walk over each file in the wiki tree.
	EachFile(path string, info os.FileInfo, err error) error

	// Summary provides the final output.
	// TODO(rjk): I should make this more complicated. In a way that
	// permits all the file actions to happen in parallel? The parsing of all
	// the articles is definitely something that can transpire concurrently.
	Summary() error
}


// ListAllWikiFiles is a boring implementation of Tidying that lists all files.
type listAllWikiFiles struct{}

func NewListAllWikiFilesTidying() Tidying {
	return &listAllWikiFiles{}
}

func (_ *listAllWikiFiles) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}
	log.Printf("%s: %v\n", path, info)
	return nil
}

func (_ *listAllWikiFiles) Summary() error {
	return nil
}


func Everyfile( settings *wiki.Settings ,  tidying Tidying) error {
	// TODO(rjk): I have a Map/Reduce op here. I could make it parallel.

	if err := filepath.Walk(settings.Wikidir, func(path string, info os.FileInfo, err error) error {
		// TODO(rjk): here is where I should filter out invalid files.
		return tidying.EachFile(path, info, err)
	}); err != nil {
		return fmt.Errorf("Everyfile walking: %v", err)
	}

	return nil
}

// TODO(rjk): This version would search the corpus
// Write me. Use the Spotlight tooling to extract a window.
func Filteredfiles() {
}