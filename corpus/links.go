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

func (wl *Wikilink) Markdown() string {
	if wl.Title != "" {
		return fmt.Sprintf("[[%s | %s]]", wl.Id, wl.Title)
	}
	return fmt.Sprintf("[[%s]]", wl.Id)
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

// TODO(rjk):
// 2 distinct interfaces:
// 1 for taking a wikilink to its corresponding file
// 1 returns an array of (wikilink, path) corresponding to a given string. (to use
// for both auto-completion and to generate the report. (Go refactor the reporting)
// But: DO NOT WRITE CODE THAT IS NOT YET NECESSARY.

// LinkToFile is implemented by objects that can return a unique or all file paths corresponding
// to a given wikilink.
type LinkToFile interface {
	// Returns a single unique path corresponding to the wikitext or error.
	Path(location, lsd, wikitext string) (string, error)

	// Returns all (absolute) paths in the wiki that would match wikitext.
	Allpaths(location, lsd, wikitext string) ([]string, error)
}

// Do I still use this?
type StubLinkToFile struct {
}

func (_ *StubLinkToFile) Path(_, _, _ string) (string, error) {
	return "", fmt.Errorf("StubLinkToFile.Path not implemented")
}
func (_ *StubLinkToFile) Allpaths(_, _, _ string) ([]string, error) {
	return nil, fmt.Errorf("StubLinkToFile.Allpaths not implemented")
}

// Allpaths returns all (absolute) paths of files in the wiki that could
// be referred to by [[wl.Id]] by using a provided index.
// TODO(rjk): this is wrongs.
func (wl *Wikilink) Allpaths(index LinkToFile) ([]string, error) {
	// TODO(rjk): Very wrong
	return index.Allpaths("", "", wl.Id)
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

func (ul *Urllink) Markdown() string {
	return fmt.Sprintf("[%s](%s)", ul.Title, ul.Url)
}

// Markdownable requires the Markdown function to produce a Markdown
// representation of the object.
type Markdownable interface {
	Markdown() string
}

type Links struct {
	// A set of outgoing wikitext links from each fullpath-specified article.
	ForwardLinks map[string]map[Wikilink]Empty

	// A set of incoming (i.e. back) links for each fullpath-specified article.
	BackLinks map[string]map[Wikilink]Empty

	// A set of outgoing URLs from each fullpath-specified article.
	OutUrls map[string]map[Urllink]Empty

	// Forward wikitext links that do not unambiguously refer to a specific target.
	// TODO(rjk): I should track *why* they're damaged.
	DamagedLinks map[string]map[Wikilink]Empty

	// mapper instance takes a wikitext link to its corresponding filename.
	mapper LinkToFile
}

func MakeLinks(mapper LinkToFile) *Links {
	return &Links{
		ForwardLinks: make(map[string]map[Wikilink]Empty),
		BackLinks:    make(map[string]map[Wikilink]Empty),
		OutUrls:      make(map[string]map[Urllink]Empty),
		DamagedLinks: make(map[string]map[Wikilink]Empty),
		mapper: mapper,
	}
}

// AddWikilink adds a URLs leaving the node. There is no node for them
// to point to so the destination URL is nil.
func (links *Links) AddWikilink(displaytext, wikitext, filepath string) {
	urlref := MakeWikilink(wikitext, displaytext)

	// TODO(rjk): Allpaths requires some kind of index whether that's the
	// index provided by Spotlight or something lower-tech.
	paths, err := urlref.Allpaths(&StubLinkToFile{})
	if err != nil || len(paths) == 0 || len(paths) > 1 {
		log.Printf("wikilink %v in %q experienced an error: %v or is missing or ambiguous: %v", urlref, filepath, err, paths)

		perfilemap, ok := links.DamagedLinks[filepath]
		if ok {
			perfilemap[urlref] = Empty{}
		} else {
			perfilemap = make(map[Wikilink]Empty)
			perfilemap[urlref] = Empty{}
			links.DamagedLinks[filepath] = perfilemap
		}
		return
	}

	perfilemap, ok := links.ForwardLinks[filepath]
	if ok {
		perfilemap[urlref] = Empty{}
		log.Println(filepath)
	} else {
		perfilemap = make(map[Wikilink]Empty)
		perfilemap[urlref] = Empty{}
		links.ForwardLinks[filepath] = perfilemap
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
