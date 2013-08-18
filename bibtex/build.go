/*
	Process extraKeys
*/
package bibtex

import (
//    "github.com/rjkroege/wikitools/article"
	"strings"
)

/*
	Reports error conditions.
*/
type BibTeXError struct {
    what string
}

func (e *BibTeXError) Error() string {
	return e.what
}

/*
	Removes all non 'bib-.*' keys from extrakeys, removes the bibliography
	prefix and returns the resulting list.
*/
func FilterExtrakeys(extrakeys []string) (filtered []string) {
	filtered = make([]string, 0, len(extrakeys))	
	for _, s := range(extrakeys) {
		if strings.HasPrefix(s, "bib-") {
			filtered = append(filtered, s[len("bib-"):])
		}
	}
	return
}

/*
	Finds the entry type from an array of article tags. The bibtex entry type is specified
	either implicitly (we have a @book tag) or we have a @book tag and a supplementary
	@bibtex-(.*) tag where the matched substring is the entry types. Entries are as per
	the documentation such as 
	http://newton.ex.ac.uk/tex/pack/bibtex/btxdoc/node6.html#SECTION00031000000000000000
	Entry types do not include the leading '@'
*/
func ExtractBibTeXEntryType(tags []string) (entry string, biberror error) {
	entry_set := 0
	book_tag_present := false
	for _, s := range(tags) {
		switch {
		case  s == "@book":
			book_tag_present = true
		case strings.HasPrefix(s, "@bibtex-"):
			entry = s[len("@bibtex-"):]
			entry_set ++
		}
	}

	switch {
	case entry_set == 0:
		entry = "book"
	case entry_set > 1:
		biberror = &BibTeXError{"More than one supplementary @bibtex-(.*) tag."}
	case !book_tag_present:
		biberror = &BibTeXError{ "No book tag present." }
	default:
		biberror = nil
	}
	return
}

/*
	What's going on here. An article can have exta keys. We want to build bibtex
	entries from these.

	Presumption: extra keys matching bib-.* are valid.
	Presumption: we have an argument that specifies which kind of entry it is. (It needs to be supplementary
	to the article type.) 

	TODO(rjkroege): Need an entry type.
	TODO(rjkroege): Need a stable mechanism for generating a book key. Does it have to be stable.

	Why: because I can't remember the kind of entries that BibTeX would expect.

	Presumption: we pick 

	entrytype is the type of entry
	filteredkeys are all extra keys matching bib-.*

*/
//func CreateBibTexEntry(tags []string, extrakeys map[string]string) string, err {
	/*
		Use entrytype to pick a template. Populate the template for that type from
		the filteredkeys. Validate that we have the right kind of keys.
	*/
//	filtered := FilterExtrakeys(extrakeys)

	/*
		TODO(rjkroege): Are here. Decide how to get the article type.	
		I could make it part of the toplevel tags. That sort of pollutes the
		top-level namespace? In particular: journal, book are semantically
		overlapped with the bibtex entry types. 

		I could have bibtex/book bibtex/article etc.?

		I could have @book @article

		I could have @book @bibtex-article. Which defaults to @book @bibtex-book. Hm. I like it.
	*/
//	entrytype, _ := FindEntryType(tags)
//	entrytype, _ := validator_table[
//}

/*
	We validate the entries.

	Note that there is as yet no reporting structure so there is some work here to
	support that. This work is not currently in scope.

	Each kind of BibTex entry has a required list of keys. Make sure that they're present.
*/
//var validatator_table map[string][]string;

// TODO(rjkroege): Insert the templates into this.

