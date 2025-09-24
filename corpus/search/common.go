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
// Historical note: written while on an unhelpful sidequest to fix a bug
// in Wikitext.
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

// buildshortestwikitext generates the shortest wikitext that can
// reference topath (absolute path of the desired destination) given
// allpaths in the wiki that share the same basename with the given wiki
// root.
func buildshortestwikitext(root, topath string, allpaths []string) (string, error) {
	// 2. split the topath into a bundle (たば), chopping the separator.
	束 := strings.Split(strings.Trim(topath, string(filepath.Separator)), string(filepath.Separator))

	// 3. Find the shortest unique prefix
	// TODO(rjk): Could improve with dynamic programming but len(allpaths) typically small.
	b束 := len(
		strings.Split(
			strings.Trim(root, string(filepath.Separator)),
			string(filepath.Separator)))
	c := 0
	for j := len(束) - 1; j > b束; j-- {
		s束 := filepath.Join(束[j:]...)
		c = counter(allpaths, s束)
		if c == 1 {
			return s束, nil
		}
	}

	return "", NoValidMatch
}

func counter(束 []string, s string) int {
	c := 0
	for _, k := range 束 {
		if strings.HasSuffix(k, s) {
			c++
		}
	}
	return c
}
