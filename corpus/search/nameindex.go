package search

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unique"

	"github.com/rjkroege/wikitools/wiki"
)

type wikilinkIndexerimpl struct {
	wikiroot string
	index    map[string][][]unique.Handle[string]
}

var _ wiki.LinkToFile = (*wikilinkIndexerimpl)(nil)

func (spix *wikilinkIndexerimpl) Close() {}

// Returns a single unique path corresponding to the wikitext found in a
// file in directory lsd, limiting search to files found recursively in location or
// error when impossible.
func (spix *wikilinkIndexerimpl) Path(location, lsd, wikitext string) (string, error) {
	// TODO(rjk): Adjust in the future as needed to support different
	// extensions.
	if filepath.Ext(wikitext) == "" {
		wikitext = wikitext + ".md"
	}

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

func (_ *wikilinkIndexerimpl) Allpaths(location, lsd, wikitext string) ([]string, error) {
	return nil, fmt.Errorf("StubLinkToFile not implemented")
}

// Wikitext returns a wikitext such that clicking on it in file frompath
// will open file topath or an error if it was impossible to do so.
// Wikitext and its dependencies assume that both frompath and topath
// exist.
// TODO(rjk): consider making suffix stripping configurable. For example,
// I expect that I'd want svg etc to keep its suffix?
func (spix *wikilinkIndexerimpl) Wikitext(frompath, topath string) (string, error) {
	base := filepath.Base(topath)
	allpaths, err := spix.pathsforwikitext(spix.wikiroot, base)
	if err != nil {
		return "", fmt.Errorf("Wikitext pathsforwikitext %w", err)
	}

	if len(allpaths) == 1 {
		return base, nil
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

// Only have one index.
 var (
instance *wikilinkIndexerimpl
   once     sync.Once
 )

 func MakeWikilinkNameIndex(wikiroot string) *wikilinkIndexerimpl {
     once.Do(func() {
        instance = implMakeWikilinkNameIndex(wikiroot)
   })
  return instance
 }

// indexOneFile adds a single file to the index.
func indexOneFile(index map[string][][]unique.Handle[string], path string, d os.DirEntry, err error) error {
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
}

// TODO(rjk): Factor out the inner code as a separate entry point for
// adding new files. In particular, there will be a larger design document
// for updating the index cache.
func implMakeWikilinkNameIndex(wikiroot string) *wikilinkIndexerimpl {
	index := make(map[string][][]unique.Handle[string])

	if err := filepath.WalkDir(wikiroot, func(path string, d os.DirEntry, err error) error {
		return indexOneFile(index, path, d, err)
	}); err != nil {
		log.Fatalf("can't continue without an index")
	}

	spidx := &wikilinkIndexerimpl{
		wikiroot: wikiroot,
		index:    index,
	}
	return spidx
}

// pathsforwikitext returns all the absolute paths for resources in
// directory tree specified by location with leaf path wikitextfile.
// wikitextfile is the complete file path (i.e. includes the extension.)
func (spix *wikilinkIndexerimpl) pathsforwikitext(location, wikitextfile string) ([]string, error) {
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
