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

	destpath, err := links.mapper.Path(links.location, filepath.Dir(fpath), wlref.Id)
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
func (links *Links) addForwardUrl(urlref corpus.Urllink, fpath string) {
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
// This code runs on an arbitrary thread.
func (lr *linkRecording) RecordUrl(displaytext, url string) {
	urlref := corpus.MakeUrllink(url, displaytext)
	lr.urls = append(lr.urls, urlref)
}

// TODO(rjk): The API will change for concurrent access.
// This code runs on an arbitrary thread.
func (lr *linkRecording) RecordWikilink(displaytext, wikitext string) {
	log.Println("RecordWikilink", displaytext, wikitext, lr.filepath)
	wikiref := corpus.MakeWikilink(wikitext, displaytext)
	lr.wikis = append(lr.wikis, wikiref)
}

// This code must run on the thread that owns the links database.
// Serializes the updates.
func (lr *linkRecording) commitOnLinksOwner() {
	lr.links.remove(lr.filepath)
	for _, url := range lr.urls {
		lr.links.addForwardUrl(url, lr.filepath)
	}
	for _, wikiref := range lr.wikis {
		lr.links.addWikilink(wikiref, lr.filepath)
	}
}

// Run this on an arbitrary go routine.
func (lr *linkRecording) Commit() {
	log.Println("Commit")
	lr.commitOnLinksOwner()
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
	return &linkRecording{
		filepath: filepath,
		links:    links,
		urls:     []corpus.Urllink{},
		wikis:    []corpus.Wikilink{},
	}
}

// Show that Linkminer is a LinksRecorder
var _ corpus.LinkRecording = (*linkRecording)(nil)

func (links *Links) remove(fpath string) {
	fwlinks := links.ForwardLinks[fpath]
	for k, _ := range fwlinks {
		destpath, err := links.mapper.Path(links.location, filepath.Dir(fpath), k.Id)
		if err != nil {
			// In theory, this should never happen (perhaps if a file was removed?)
			log.Printf("the destination file of wikilink %v would seem to be missing: %v", k, err)
			continue
		}

		backtext, err := links.mapper.Wikitext(destpath, fpath)
		if err != nil {
			log.Printf("links.mapper.Wikitext from %q to %q failed: %v", fpath, destpath, err)
			continue
		}
		backref := corpus.MakeWikilink(backtext, "")
		delete(links.BackLinks[destpath], backref)
	}

	delete(links.DamagedLinks,fpath)
	delete(links.ForwardLinks,fpath)
}
