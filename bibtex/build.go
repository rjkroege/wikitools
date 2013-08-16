/*
	Process extraKeys
*/
package bibtex

import (
//    "github.com/rjkroege/wikitools/article"
	"strings"
)

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
//func CreateBibTexEntry(entrytype string, filteredkeys map[string]string) string, err {
	/*
		Use entrytype to pick a template. Populate the template for that type from
		the filteredkeys. Validate that we have the right kind of keys.
	*/
//}

/*
	We validate the entries.

	Note that there is as yet no reporting structure so there is some work here to
	support that. This work is not currently in scope.

	Each kind of BibTex entry has a required list of keys. Make sure that they're present.
*/
//var validatator_table map[string][]string;

// TODO(rjkroege): Insert the templates into this.

