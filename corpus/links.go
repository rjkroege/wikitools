package corpus

import (
	"fmt"
	"log"
)

// TODO(rjk): overview I am working towards a system where I get
// backlinks added to the Edwood tag area. For this to happen on a `we`
// invocation, `we` needs to open the database of backlinks for this
// file. And find the backlinks for the just-opened file and add them to
// title area.
//
// The link as added to the window tag area should be the minimum form
// such that `wikilink` can find a unique wiki article.
//
// There is no need to deal with links that aren't wikilinks. They can be
// added to the report but they don't need need to be persisted. However
//
// The database must also support guidance for moving files.
//
// Full paths are the unambiguous name of items in the wiki. Note that
// this means that the wiki location (prefix) is perhaps redundantly
// encoded. But absolute paths are cannonically absolute and completely
// unambiguous.
//
// Note that a link is not just a `string` path. Additional attributes
// are necessary for reporting and content generation such as the
// "display name".

type Empty struct{}

// Wikilink holds a wikilink's [[id | title ]] where the title is
// optional. Note that per
// https://ia.net/writer/support/library/wikilinks that there can also be
// a [[rootspec: id]] where rootspec is the root specifier of the
// tree in which to search for id.
type Wikilink struct {
	// The contents of the id part of the wikilink.
	Id string

	// The title portion of the wikilink.
	Title string
}

func MakeWikilink(id, title string) Wikilink {
	return Wikilink{
		Id:    id,
		Title: title,
	}
}

// TODO(rjk): These stubs need refinement. They are here with comments to
// capture my API surface thinking.

// MakeWikilinkFromPathPair creates a Wikilink object such that right-clicking on
// the result's string representation in frompath will open topath.
// TODO(rjk): somewhere there will be code that will determine the actual
// paths. Line that up with some tests.
func MakeWikilinkFromPathPair(frompath, topath string) Wikilink {
	return Wikilink{
		Id: "not implemented",
	}
}

// WikilinkNameIndex defines an interface implementing a per-corpus
// index of wikilink text to path names.
type WikilinkNameIndex interface {
	// Returns all (absolute) paths in the wiki that would match wikitext.
	Allpaths(wikitext string) ([]string, error)
}

type StubWikilinkNameIndex struct {
}

func (_ *StubWikilinkNameIndex) Allpaths(_ string) ([]string, error) {
	return nil, fmt.Errorf("StubWikilinkNameIndex not implemented")
}

// Allpaths returns all (absolute) paths of files in the wiki that could
// be referred to by [[wl.Id]] by using a provided index.
func (wl *Wikilink) Allpaths(index WikilinkNameIndex) ([]string, error) {
	return index.Allpaths(wl.Id)
}

// A Urllink holds a Markdown URL where there is a [title](http://foo.foo) structure.
// TODO(rjk): This is large overlap between these structures. Consider refactoring
// them together later.
type Urllink struct {
	Url   string
	Title string
}

func MakeUrllink(url, title string) Urllink {
	return Urllink{
		Url:   url,
		Title: title,
	}
}

type Links struct {
	// A set of outgoing wikitext links from each fullpath-specified article.
	ForwardLinks map[string]map[Wikilink]Empty

	// A set of incoming (i.e. back) links for each fullpath-specified article.
	BackLinks map[string]map[Wikilink]Empty

	// A set of outgoing URLs from each fullpath-specified article.
	OutUrls map[string]map[Urllink]Empty

	// TODO(rjk): Record the d
}

func MakeLinks() *Links {
	return &Links{
		ForwardLinks: make(map[string]map[Wikilink]Empty),
		BackLinks:    make(map[string]map[Wikilink]Empty),
		OutUrls:      make(map[string]map[Urllink]Empty),
	}
}

// AddWikilink adds a URLs leaving the node. There is no node for them
// to point to so the destination URL is nil.
func (links *Links) AddWikilink(displaytext, wikitext, filepath string) {
	urlref := MakeWikilink(wikitext, displaytext)

	perfilemap, ok := links.ForwardLinks[filepath]
	if ok {
		perfilemap[urlref] = Empty{}
		log.Println(filepath)
	} else {
		perfilemap = make(map[Wikilink]Empty)
		perfilemap[urlref] = Empty{}
		links.ForwardLinks[filepath] = perfilemap
	}

	// TODO(rjk): Allpaths requires some kind of index whether that's the
	// index provided by Spotlight or something lower-tech.
	paths, err := urlref.Allpaths(&StubWikilinkNameIndex{})
	if err != nil {
		log.Printf("not back-linking %v in %q because %v", urlref, filepath, err)
		return
	}

	if len(paths) == 0 {
		log.Printf("outgoing wikilink %v in %q points at nothing", urlref, filepath)
		return
	}

	if len(paths) > 1 {
		log.Printf("outgoing wikilink %v in %q is ambiguous and could refer to %v", urlref, filepath, paths)
		return
	}

	destpath := paths[0]
	backref := MakeWikilinkFromPathPair(filepath, destpath)

	// Update the reverse links.
	perfilemap, ok = links.BackLinks[filepath]
	if ok {
		perfilemap[backref] = Empty{}
	} else {
		perfilemap = make(map[Wikilink]Empty)
		perfilemap[backref] = Empty{}
		links.BackLinks[filepath] = perfilemap
	}
}

// AddForwardUrl adds a URLs leaving the node. There is no node for them
// to point to so the destination URL is nil.
func (links *Links) AddForwardUrl(displaytext, url, filepath string) {
	urlref := MakeUrllink(url, displaytext)

	perfilemap, ok := links.OutUrls[filepath]
	if ok {
		perfilemap[urlref] = Empty{}
		log.Println(filepath)
	} else {
		perfilemap = make(map[Urllink]Empty)
		perfilemap[urlref] = Empty{}
		links.OutUrls[filepath] = perfilemap
	}
}
