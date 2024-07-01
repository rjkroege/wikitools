package corpus

import (
	"fmt"
	"log"
	"path/filepath"
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

// LinkToFile is implemented by objects that can return a unique or all file paths corresponding
// to a given wikilink.
type LinkToFile interface {
	// Returns a single unique path corresponding to the wikitext with
	// location (i.e. root) and found in file lsd.
	Path(location, lsd, wikitext string) (string, error)

	// Returns all (absolute) paths in the wiki that would match wikitext.
	Allpaths(location, lsd, wikitext string) ([]string, error)

	// Wikitext returns a wikitext such that clicking on it in frompath will
	// open topath or an error if it was impossible to do so.
	Wikitext(frompath, topath string) (string, error)
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

	// location is the root of the wiki tree
	location string
}

// Show that Linkminer is a UrlRecorder
var _ UrlRecorder = (*Links)(nil)

func MakeLinks(mapper LinkToFile, location string) *Links {
	return &Links{
		ForwardLinks: make(map[string]map[Wikilink]Empty),
		BackLinks:    make(map[string]map[Wikilink]Empty),
		OutUrls:      make(map[string]map[Urllink]Empty),
		DamagedLinks: make(map[string]map[Wikilink]Empty),
		mapper:       mapper,
		location:     location,
	}
}

// AddWikilink adds a URLs leaving the node. There is no node for them
// to point to so the destination URL is nil.
func (links *Links) AddWikilink(displaytext, wikitext, fpath string) {
	urlref := MakeWikilink(wikitext, displaytext)

	if filepath.Ext(wikitext) == "" {
		wikitext = wikitext + ".md"
	}

	log.Println("AddWikilink", links.location, fpath, wikitext)
	destpath, err := links.mapper.Path(links.location, filepath.Dir(fpath), wikitext)
	if err != nil {
		log.Printf("links.mapper.Path on %q, [[%s]] error: %v", fpath, wikitext, err)

		perfilemap, ok := links.DamagedLinks[fpath]
		if ok {
			perfilemap[urlref] = Empty{}
		} else {
			perfilemap = make(map[Wikilink]Empty)
			perfilemap[urlref] = Empty{}
			links.DamagedLinks[fpath] = perfilemap
		}
		return
	}

	perfilemap, ok := links.ForwardLinks[fpath]
	if ok {
		perfilemap[urlref] = Empty{}
		log.Println(fpath)
	} else {
		perfilemap = make(map[Wikilink]Empty)
		perfilemap[urlref] = Empty{}
		links.ForwardLinks[fpath] = perfilemap
	}

	backtext, err := links.mapper.Wikitext(fpath, destpath)
	if err != nil {
		log.Printf("links.mapper.Wikitext from %q to %q failed: %v", fpath, destpath, err)
		return
	}
	backref := MakeWikilink(backtext, "")

	// Update the reverse links.
	perfilemap, ok = links.BackLinks[fpath]
	if ok {
		perfilemap[backref] = Empty{}
	} else {
		perfilemap = make(map[Wikilink]Empty)
		perfilemap[backref] = Empty{}
		links.BackLinks[fpath] = perfilemap
	}
}

// AddForwardUrl adds a URLs leaving the node. There is no node for them
// to point to so the destination URL is nil.
func (links *Links) AddForwardUrl(displaytext, url, fpath string) {
	urlref := MakeUrllink(url, displaytext)

	perfilemap, ok := links.OutUrls[fpath]
	if ok {
		perfilemap[urlref] = Empty{}
		log.Println(fpath)
	} else {
		perfilemap = make(map[Urllink]Empty)
		perfilemap[urlref] = Empty{}
		links.OutUrls[fpath] = perfilemap
	}
}

// TODO(rjk): The presence of this forwarder suggests that I might want
// to change the UrlRecorder interface?
func (links *Links) RecordUrl(displaytext, url, filepath string) {
	links.AddForwardUrl(displaytext, url, filepath)
}

func (links *Links) RecordWikilink(displaytext, wikitext, filepath string) {
	links.AddWikilink(displaytext, wikitext, filepath)
}
