package search

import (
	"path/filepath"
	"strings"
	"errors"
)


var AmbiguousWikitext = errors.New("ambiguous wikitext")
var NoFileForWikitext = errors.New("no matching wikitext")

// disambiguatewikipaths takes local source directory lsd in which the
// contents containing a wikitext was found, the location of the wiki
// tree and an array of absolute paths and proceeds to find an
// unambiguous match or fail.
//
// wikitext can have a (possibly empty) disambiguating prefix.
// 
// Wiki link resolution is taken from
// https://ia.net/writer/support/library/wikilinks: "When you open a
// wikilink like [[your link]], iA Writer finds the nearest file that
// matches the given name. A file in the same folder be preferred over a
// file in a subfolder, and a file in a subfolder will be preferred over
// a file in a parent folder."
// 
// Hence, resolve ambiguous matches by making three passes over the list
// of file matches to find lsd/wikitext, lsd/.../wikitext,
// location/.../wikitext.
//
// TODO(rjk): Determine how to handle extensions.
func disambiguatewikipaths(location, lsd, wikitext string, allpaths []string) (string, error) {
	if len(allpaths) == 0 {
// TODO(rjk): This is perhaps not the ideal behaviour but it aligns with
// the current shell script behaviour. Since the tools are becoming
// increasingly Edwood-specific, I can imagine writing some kind of
// complaint to $location/+Errors in this context.
		return "", NoFileForWikitext
	}

	// Unambiguous single file name. All is well.
	if len(allpaths) == 1 {
		return allpaths[0], nil
	}

	expanded := filepath.Join(lsd, wikitext)
	matches := 0
	onematch := ""
	for _, p := range(allpaths) {
		if strings.HasSuffix(p, expanded) {
			matches++
			onematch = p
		}
	}

	if matches == 1 {
		return onematch, nil
	}

	expanded = filepath.Join(location, wikitext)
	matches = 0
	onematch = ""
	for _, p := range(allpaths) {
		if strings.HasSuffix(p, expanded) {
			matches++
		}
	}

	if matches == 1 {
		return onematch, nil
	} else {
		// TODO(rjk): make a more useful error message that I can dump somewhere.
		return "", AmbiguousWikitext
	}
}


// The auto-complete functionality (i.e. dismbiguating) string needs to
// find the shortest set of paths needed to unique the name either w.r.t.
// the root of the tree (i.e. ~/Documents/wiki) or the directory
// containing the link origin. (Call this the origin file.)
// 
// NB: Having found the correct prefix, it is sufficient to glue together
// the prefix with the file name and see if there is a file in the match
// list where that string is a suffix.
// 
// What about the generating the prefix + filename? That's a separate
// problem. But: I don't have to be as smart as iaWriter? i.e. My
// auto-completes don't have to be minimal? Correct. The minimal prefix
// is not necessary. So just auto-complete with either no prefix for
// it's a unique name or in directory or the prefix w.r.t. origin or the
// prefix w.r.t. root.

