package search

import (
	"errors"
	"path/filepath"
	"strings"
)

var AmbiguousWikitext = errors.New("ambiguous wikitext")
var NoFileForWikitext = errors.New("no matching wikitext")
var EmptyWikitextFile = errors.New("wikitext file portion is empty")
var NoValidMatch = errors.New("buildshortestwikitext no valid match")

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

	matches := 0
	onematch := ""
	for _, p := range allpaths {
		if strings.HasSuffix(p, wikitext) && strings.HasPrefix(p, lsd) {
			matches++
			onematch = p
		}
	}

	if matches == 1 {
		return onematch, nil
	}

	matches = 0
	onematch = ""
	for _, p := range allpaths {
		if strings.HasSuffix(p, wikitext) && strings.HasPrefix(p, location) {
			matches++
			onematch = p
		}
	}

	if matches == 1 {
		return onematch, nil
	} else {
		return "", AmbiguousWikitext
	}
}

// commonPrefixSplit returns the longest common directory prefix of two
// absolute paths together with the two differing suffixes.
// TODO(rjk): I had wanted: Path(root, a, commonPrefixSplit(a,b).suffixB) == b
// but it is not clear if this is the case.
func commonPrefixSplit(a, b string) (prefix, suffixA, suffixB string) {
	a = filepath.Clean(a)
	b = filepath.Clean(b)

	aParts := strings.Split(a, string(filepath.Separator))
	bParts := strings.Split(b, string(filepath.Separator))

	// Find the split point.
	split := 0
	for ; split < len(aParts) && split < len(bParts); split++ {
		if aParts[split] != bParts[split] {
			break
		}
	}

	common := aParts[:split]
	suffixAParts := aParts[split:]
	suffixBParts := bParts[split:]

	prefix = filepath.Join(common...)
	if !filepath.IsAbs(prefix) {
		prefix = string(filepath.Separator) + prefix
	}
	suffixA = filepath.Join(suffixAParts...)
	suffixB = filepath.Join(suffixBParts...)

	return prefix, suffixA, suffixB
}
