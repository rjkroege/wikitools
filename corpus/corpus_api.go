package corpus

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
	"github.com/rjkroege/wikitools/wiki"
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

type Link interface {
	Wikilink | Urllink
	Sortid() string
	Markdown() string
	Html() string
}

type LinkMap[T Link] map[T]empty

func MakeWikilink(id, title string) Wikilink {
	return Wikilink{
		Id:    id,
		Title: title,
	}
}

func foo[E Link](t E) {
	log.Println(t.Sortid())
}

// SortWikilinks sorts a slice of Wikilink structures by their Id field.
func SortWikilinks[S ~[]E, E Link](links S) {
	slices.SortFunc(links, func(a, b E) int {
		sa, sb := a.Sortid(), b.Sortid()
		if sa < sb {
			return -1
		} else if sa > sb {
			return 1
		}
		return 0
	})
}

// Helpful definitions to not need to change code
type WikilinkMap = LinkMap[Wikilink]
type UrlMap = LinkMap[Urllink]

func Stringify[T Link](links LinkMap[T]) string {
	keys := maps.Keys(links)
	SortWikilinks(keys)

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k.Markdown())
	}
	return b.String()
}

func (wl Wikilink) Markdown() string {
	if wl.Title != "" {
		return fmt.Sprintf("[[%s | %s]]", wl.Id, wl.Title)
	}
	return fmt.Sprintf("[[%s]]", wl.Id)
}

func (wl Wikilink) Html() string {
	if wl.Title != "" {
		return fmt.Sprintf("<a href=\"plumb://w/%s\">%s</a>", wl.Id, wl.Title)
	}
	return fmt.Sprintf("<a href=\"plumb://w/%s\">%s</a>", wl.Id, wl.Id)
}

// Allpaths returns all (absolute) paths of files in the wiki that could
// be referred to by [[wl.Id]] by using a provided index.
// TODO(rjk): Check if this is working.
func (wl Wikilink) Allpaths(index wiki.LinkToFile) ([]string, error) {
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

func (ul Urllink) Markdown() string {
	return fmt.Sprintf("[%s](%s)", ul.Title, ul.Url)
}

func (ul Urllink) Html() string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", ul.Url, ul.Title)
}

func (wl Wikilink) Sortid() string { return wl.Id }
func (ul Urllink) Sortid() string  { return ul.Url }
