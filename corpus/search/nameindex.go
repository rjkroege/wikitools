package search

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/corpus"
)

// spotlightWikilinkIndexer hides all of the darwin-specific code needed
// to support indexing using spotlight.
// Changes to this structure need to be synchronized correctly.
type spotlightWikilinkIndexer struct {
	wikiroot string
}

var _ corpus.LinkToFile = (*spotlightWikilinkIndexer)(nil)

func (spix *spotlightWikilinkIndexer) Path(location, lsd, wikitext string) (string, error) {
	basepart := filepath.Base(wikitext)
	if basepart == "" {
		return "", EmptyWikitextFile
	}

	allpaths, err := spix.pathsforwikitext(location, basepart)
	if err != nil {
		return "", err
	}

	return disambiguatewikipaths(location, lsd, wikitext, allpaths)
}

func (_ *spotlightWikilinkIndexer) Allpaths(location, lsd, wikitext string) ([]string, error) {
	return nil, fmt.Errorf("StubLinkToFile not implemented")
}

// pathsforwikitext returns all the absolute paths for resources in directory tree specified by
// location with leaf path wikitextfile.
// TODO(rjk): the input text might or might not have a file name extension. I'm currently
// not clear about that.
// TODO(rjk): does it need an object?
func (_ *spotlightWikilinkIndexer) pathsforwikitext(location, wikitextfile string) ([]string, error) {
    matches :=  []string{}
    err := filepath.WalkDir(location, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return nil
        }
        if filepath.Base(path) == wikitextfile {
            abs, err := filepath.Abs(path)
            if err != nil {
                return err
            }
            matches = append(matches, abs)
        }
        return nil
    })
    return matches, err
}

func MakeWikilinkNameIndex(wikiroot string) *spotlightWikilinkIndexer {
	// The Apple docs imply (very strongly) that there can only be a single
	// query running at a time. Remember this if I should convert the tidy
	// code to run concurrently.
	spidx := &spotlightWikilinkIndexer{
		wikiroot: wikiroot,
	}
	return spidx
}

// make suffix stripping configurable. I expect that I'd want svg etc to keep its
// suffix?
func (spix *spotlightWikilinkIndexer) Wikitext(frompath, topath string) (string, error) {
	// Find allz of the paths
	allpaths, err := spix.pathsforwikitext(filepath.Dir(frompath), filepath.Base(topath))
	if err != nil {
		return "", fmt.Errorf("Wikitext pathsforwikitext %w", err)
	}
	return buildshortestwikitext(spix.wikiroot, topath, allpaths)
}
