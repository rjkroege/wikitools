package links

import (
	"log"
	"path/filepath"
	"sort"
	"sync"

	"github.com/rjkroege/wikitools/corpus"
	"golang.org/x/exp/maps"
)

type empty = struct{}

type Links struct {
	// A set of outgoing wikitext links from each fullpath-specified article.
	ForwardLinks map[string]corpus.LinkMap[corpus.Wikilink]

	// A set of incoming (i.e. back) links for each fullpath-specified article.
	BackLinks map[string]corpus.LinkMap[corpus.Wikilink]

	// A set of outgoing URLs from each fullpath-specified article.
	OutUrls map[string]corpus.LinkMap[corpus.Urllink]

	// Forward wikitext links that do not unambiguously refer to a specific target.
	// TODO(rjk): I should track *why* they're damaged.
	DamagedLinks map[string]corpus.LinkMap[corpus.Wikilink]

	// mapper instance takes a wikitext link to its corresponding filename.
	mapper corpus.LinkToFile

	// location is the root of the wiki tree
	location string
}

// Show that Linkminer is a LinksRecorder
var _ corpus.LinksRecorder = (*Links)(nil)

// Only have one index.
var (
	instance *Links
	once     sync.Once
)

// ImplMakeLinks is exposed publically only to make it easier to run tests.
// Real code should only use the future concurrent version.
// TODO(rjk): Clean this up carefully when I convert this code to be concurrent safe.
func implMakeLinks(mapper corpus.LinkToFile, location string) *Links {
	return &Links{
		ForwardLinks: make(map[string]corpus.WikilinkMap),
		BackLinks:    make(map[string]corpus.WikilinkMap),
		OutUrls:      make(map[string]corpus.LinkMap[corpus.Urllink]),
		DamagedLinks: make(map[string]corpus.WikilinkMap),
		mapper:       mapper,
		location:     location,
	}
}

func MakeLinks(mapper corpus.LinkToFile, location string) *Links {
	once.Do(func() {
		instance = implMakeLinks(mapper, location)
	})
	return instance
}

// AddWikilink updates the two-way linking data for a wikitext found in
// fpath with (optional) displaytext so that fpath points to
// Path(wikitext) and Path(wikitext) points to fpath in the link table.
// This function is called for each link found in fpath by the Markdown
// parser.
func (links *Links) addWikilink(wlref corpus.Wikilink, fpath string) {
	// urlref := corpus.MakeWikilink(wikitext, displaytext)

	// Here I fix the links to have the correct extension.
	// TODO(rjk): Adjust in the future as needed to support different
	// extensions.
	wikitext := wlref.Id
	if filepath.Ext(wikitext) == "" {
		wikitext = wikitext + ".md"
	}

	destpath, err := links.mapper.Path(links.location, filepath.Dir(fpath), wikitext)
	if err != nil {
		perfilemap, ok := links.DamagedLinks[fpath]
		if ok {
			perfilemap[wlref] = empty{}
		} else {
			perfilemap = make(corpus.WikilinkMap)
			perfilemap[wlref] = empty{}
			links.DamagedLinks[fpath] = perfilemap
		}
		return
	}

	perfilemap, ok := links.ForwardLinks[fpath]
	if ok {
		perfilemap[wlref] = empty{}
	} else {
		perfilemap = make(corpus.WikilinkMap)
		perfilemap[wlref] = empty{}
		links.ForwardLinks[fpath] = perfilemap
	}

	backtext, err := links.mapper.Wikitext(destpath, fpath)
	// Why this error?
	if err != nil {
		log.Printf("links.mapper.Wikitext from %q to %q failed: %v", fpath, destpath, err)
		return
	}
	backref := corpus.MakeWikilink(backtext, "")

	// Update the reverse links. NB: the destpath needs the update with a
	// synthesized wikilink back to fpath.
	perfilemap, ok = links.BackLinks[destpath]
	if ok {
		perfilemap[backref] = empty{}
	} else {
		perfilemap = make(corpus.WikilinkMap)
		perfilemap[backref] = empty{}
		links.BackLinks[destpath] = perfilemap
	}
}

// addForwardUrl adds a URLs leaving the node. There is no node for them
// to point to so the destination URL is nil.
func (links *Links) addForwardUrl(displaytext, url, fpath string) {
	urlref := corpus.MakeUrllink(url, displaytext)

	perfilemap, ok := links.OutUrls[fpath]
	if ok {
		perfilemap[urlref] = empty{}
	} else {
		perfilemap = make(corpus.UrlMap)
		perfilemap[urlref] = empty{}
		links.OutUrls[fpath] = perfilemap
	}
}

// TODO(rjk): The API surface will change to concurrent access.
func (lr *linkRecording) RecordUrl(displaytext, url string) {
	lr.links.addForwardUrl(displaytext, url, lr.filepath)
}

// TODO(rjk): The API will change for concurrent access.
func (lr *linkRecording) RecordWikilink(displaytext, wikitext string) {
	log.Println("RecordWikilink", displaytext, wikitext, lr.filepath)
	wikiref := corpus.MakeWikilink(wikitext, displaytext)
	lr.links.addWikilink(wikiref, lr.filepath)
}

// TODO(rjk): Currently a nop. This will change.
func (lr *linkRecording) Commit() {
	log.Println("Commit")
}

func StringVector[T corpus.Link](linkmap map[string]corpus.LinkMap[T]) []string {
	keys := maps.Keys(linkmap)
	sort.Strings(keys)

	result := make([]string, 0, len(keys))
	for _, ks := range keys {
		result = append(result, ks)
		result = append(result, corpus.Stringify(linkmap[ks]))
	}

	return result
}

type linkRecording struct {
	filepath string
	links    *Links
	urls     []corpus.Urllink
	wikis    []corpus.Wikilink
}

func (links *Links) StartRecordingForFile(filepath string) corpus.LinkRecording {
	// TODO(rjk): Insert "remove filepath" here so that updates work.
	return &linkRecording{
		filepath: filepath,
		links:    links,
		urls:     []corpus.Urllink{},
		wikis:    []corpus.Wikilink{},
	}
}

// Show that Linkminer is a LinksRecorder
var _ corpus.LinkRecording = (*linkRecording)(nil)
