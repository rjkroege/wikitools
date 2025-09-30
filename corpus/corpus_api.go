package corpus

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
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
// added to the report but they don't need need to be persisted.
//
// The database must also support guidance for moving files.
//
// Full paths are the unambiguous name of items in the wiki. Note that
// this means that the wiki location (prefix) is perhaps redundantly
// encoded. But absolute paths are unambiguous.
//
// Note that a link is not just a `string` path. Additional attributes
// are necessary for reporti

type empty = struct{}

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

type WikilinkMap map[Wikilink]empty

type UrlMap map[Urllink]empty

func MakeWikilink(id, title string) Wikilink {
	return Wikilink{
		Id:    id,
		Title: title,
	}
}

type ById []Wikilink

func (a ById) Len() int           { return len(a) }
func (a ById) Swap(i, j int)    { a[i], a[j] = a[j], a[i] }
func (a ById) Less(i, j int) bool { return a[i].Id < a[j].Id }

// SortWikilinks sorts a slice of Wikilink structures by their Id field.
func SortWikilinks(links []Wikilink) {
	sort.Sort(ById(links))
}

func (wm WikilinkMap) String() string {
	keys := maps.Keys(wm)
	SortWikilinks(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k.Markdown())
	}
	return b.String()
}

func (wl *Wikilink) Markdown() string {
	if wl.Title != "" {
		return fmt.Sprintf("[[%s | %s]]", wl.Id, wl.Title)
	}
	return fmt.Sprintf("[[%s]]", wl.Id)
}

func (wl *Wikilink) Html() string {
	if wl.Title != "" {
		return fmt.Sprintf("<a href=\"plumb://w/%s\">%s</a>",  wl.Id, wl.Title)
	}
	return fmt.Sprintf("<a href=\"plumb://w/%s\">%s</a>",  wl.Id, wl.Id)
}


// LinkToFile is implemented by objects that can return a unique or all file paths corresponding
// to a given wikilink.
type LinkToFile interface {
// Returns a single unique path corresponding to the wikitext found in
// file lsd limiting the search for target paths to files in location or
// error if this is impossible.
	Path(location, lsd, wikitext string) (string, error)

	// Returns all (absolute) paths in the wiki that would match wikitext.
	// TODO(rjk): Why is lsd here?
	Allpaths(location, lsd, wikitext string) ([]string, error)

// Wikitext returns a wikitext such that clicking on it in file frompath
// will open file topath or an error if it was impossible to do so. In
// particular: Path(wikiroot, frompath, Wikitext(frompath, topath)) ==
// topath
	Wikitext(frompath, topath string) (string, error)
}

// Allpaths returns all (absolute) paths of files in the wiki that could
// be referred to by [[wl.Id]] by using a provided index.
// TODO(rjk): Check if this is working.
func (wl *Wikilink) Allpaths(index LinkToFile) ([]string, error) {
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

func (ul *Urllink) Html() string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>",  ul.Url, ul.Title)
}

// Markdownable requires the Markdown function to produce a Markdown
// representation of the object.
type Markdownable interface {
	Markdown() string
}

type Htmlable interface {
	Html() string
}

var _ Markdownable = (*Urllink)(nil)
var _ Markdownable = (*Wikilink)(nil)
var _ Htmlable = (*Urllink)(nil)
var _ Htmlable = (*Wikilink)(nil)
