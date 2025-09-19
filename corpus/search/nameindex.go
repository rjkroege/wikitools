package search

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unique"

	"github.com/rjkroege/wikitools/corpus"
)

type spotlightWikilinkIndexer struct {
	wikiroot string
}

type neo_spotlightWikilinkIndexer struct {
	wikiroot string
	index    map[string][][]unique.Handle[string]
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

// pathsforwikitext returns all the absolute paths for resources in
// directory tree specified by location with leaf path wikitextfile.
// wikitextfile is the complete file path (i.e. includes the extension.)
func (_ *spotlightWikilinkIndexer) pathsforwikitext(location, wikitextfile string) ([]string, error) {
	matches := []string{}
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
	spidx := &spotlightWikilinkIndexer{
		wikiroot: wikiroot,
	}
	return spidx
}

// TODO(rjk): consider making suffix stripping configurable. For example,
// I expect that I'd want svg etc to keep its suffix?
func (spix *spotlightWikilinkIndexer) Wikitext(frompath, topath string) (string, error) {
	// Find allz of the paths
	allpaths, err := spix.pathsforwikitext(filepath.Dir(frompath), filepath.Base(topath))
	if err != nil {
		return "", fmt.Errorf("Wikitext pathsforwikitext %w", err)
	}
	return buildshortestwikitext(spix.wikiroot, topath, allpaths)
}

func splitPathPartsHandle(dir string) []unique.Handle[string] {
	stringparts := splitPathParts(dir)
	handles := make([]unique.Handle[string], 0, len(stringparts))
	for _, p := range stringparts {
		handles = append(handles, unique.Make(p))
	}
	return handles
}

func splitPathParts(dir string) []string {
	cleanpath := filepath.Clean(dir)
	cpp := cleanpath
	if cleanpath[0] == os.PathSeparator {
		cpp = cpp[1:]
		parts := strings.Split(cpp, "/")
		parts[0] = cleanpath[0 : len(parts[0])+1]
		return parts
	}
	return strings.Split(cpp, "/")
}

// TODO(rjk): Factor out the inner code as a separate entry point for
// adding new files.
func NeoMakeWikilinkNameIndex(wikiroot string) *neo_spotlightWikilinkIndexer {
	index := make(map[string][][]unique.Handle[string])

	if err := filepath.WalkDir(wikiroot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		dir := filepath.Dir(path)
		base := filepath.Base(path)

		if !d.IsDir() {
			pl := index[base]
			pl = append(pl, splitPathPartsHandle(dir))
			index[base] = pl
		}
		return nil
	}); err != nil {
		log.Fatalf("can't continue without an index")
	}

	spidx := &neo_spotlightWikilinkIndexer{
		wikiroot: wikiroot,
		index:    index,
	}
	return spidx
}

// pathsforwikitext returns all the absolute paths for resources in
// directory tree specified by location with leaf path wikitextfile.
func (spix *neo_spotlightWikilinkIndexer) neo_pathsforwikitext(location, wikitextfile string) ([]string, error) {
	pls, ok := spix.index[wikitextfile]
	if !ok {
		return []string{}, nil
	}

	// TODO(rjk): strip the spix.wikiroot from location?
	locationparts := splitPathParts(location)
	matches := []string{}

	for _, pl := range pls {
		if len(pl) < len(locationparts) {
			continue
		}

		mc := 0
		for i, p := range locationparts {
			if s := pl[i].Value(); s == p {
				mc++
			}
		}

		if mc == len(locationparts) {
			npps := []string{}
			for _, h := range pl {
				npps = append(npps, h.Value())
			}
			npps = append(npps, wikitextfile)
			matches = append(matches, filepath.Join(npps...))
		}
	}
	return matches, nil
}
