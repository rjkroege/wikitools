/*
	Process extraKeys
*/
package bibtex

import (
//    "github.com/rjkroege/wikitools/article"
	"strings"
	"sort"
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
	Creates a new map with the keys filtered and with the leading 'bib-' removed.
*/
func FilterExtrakeys(extrakeys map[string]string) (filtered map[string]string, keys []string) {
	filtered = make(map[string]string)	
	keys = make([]string, 0, len(extrakeys))
	for s, _ := range(extrakeys) {
		if strings.HasPrefix(s, "bib-") {
			nk := s[len("bib-"):]
			keys = append(keys, nk)
			filtered[nk] = extrakeys[s]
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

// TODO(rjkroege): Insert the code to generate the necessary goo
var required_fields map[string][]string



// TODO(rjkroege): add the bibkey
func init() {
	required_fields = make(map[string][]string)
	required_fields["article"] = []string{"author", "title", "journal", "year"}
	required_fields["book"] = []string{"author", "title", "publisher", "year"}
	required_fields["booklet"] = []string{"title"}
	required_fields["inbook"] = []string{"author", "title", "chapter", "publisher", "year"}
	required_fields["incollection"] = []string{"author", "title", "booktitle", "publisher", "year"}
	required_fields["inproceedings"] = []string{"author", "title", "booktitle", "year"}
	required_fields["manual"] = []string{"title"}
	required_fields["mastersthesis"] = []string{"author", "title", "school", "year"}
	required_fields["misc"] = []string{}
	required_fields["phdthesis"] = []string{"author", "title", "school", "year"}
	required_fields["proceedings"] = []string{"title", "year"}
	required_fields["techreport"] = []string{"author", "title", "institution", "year"}
	required_fields["unpublished"] = []string{"author", "title", "note"}

	for _, s := range(required_fields) {
		// TODO(rjkroege): you are here...
		// required_fields = append(required_fields
		sort.Strings(s)
	}
}

/*
	Linear search: true if t is in fields.
*/
func isinlist(fields []string, t string) bool {
	for _, s := range(fields) {
		if s == t {
			return true
		}
	}
	return false
}

/*
	Intersects fields and expected returning members from expected
	that are not present in fields. fields and rf must both be
	sorted.
*/
func intersectsorted(fields []string, rf []string) []string {
	missing := []string{}
	i := 0
	for _, r := range(rf) {
		for ; i < len(fields) && fields[i] < r; i++ { }
		if i >= len(fields) || fields[i] > r {
			missing = append(missing, r)
		}
	}
	return missing
}

/*
	Generates a BibTexError instance for entrytype for all the missing fields.
*/
func createerror(entrytype string, missing []string) error {
		return &BibTeXError{ "Missing required fields " + strings.Join(missing, " ") +
			 " for entry type " + entrytype };
}

/*
	Determines if the list of BibTeX fields from the article has all the required
	fields for the associated entry type.
*/
func VerifyRequiredFields(entrytype string, fields []string) error {
	_, err := required_fields[entrytype]
	if !err {
		return &BibTeXError{"Invalid entry type"}
	}

	sort.Strings(fields)

	switch entrytype {
	case "book":
		// handle the or cases.
		missing := intersectsorted(required_fields["book"], fields)
		missing_editor :=  intersectsorted(required_fields["book-editor"], fields)
		switch {
		case len(missing) == 0 || len(missing_editor) == 0:
			return nil
		case len(missing_editor) < len(missing):
			return createerror(entrytype, missing_editor)
		default:
			return createerror(entrytype, missing)
		}			
	case "inbook":
		missing := intersectsorted(required_fields["inbook"], fields)
		missing_editor :=  intersectsorted(required_fields["inbook-editor"], fields)
		switch {
		case len(missing) == 0 || len(missing_editor) == 0:
			return nil
		case len(missing) == 1 && missing[0] == "chapter":
			return nil
		case len(missing) == 1 && missing[0] == "pages":
			return nil
		case len(missing_editor) == 1 && missing_editor[0] == "chapter":
			return nil
		case len(missing_editor) == 1 && missing_editor[0] == "pages":
			return nil
		case len(missing_editor) < len(missing):
			return createerror(entrytype, missing_editor)
		default:
			return createerror(entrytype, missing)
		}
	default:
		missing := intersectsorted(required_fields[entrytype], fields)
		if len(missing) == 0 {
			return nil
		} else {
			return createerror(entrytype, missing)
		}
	}
	return &BibTeXError{"internal error: missed a case"}
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
	filteredkeys are all extra fields matching bib-.*

*/
//func CreateBibTexEntry(tags []string, extrakeys map[string]string) string, err {
	/*
		This function is ill-considered. We get a map. We need to produce two things
		A map from field->value without the bib- and a list of field names so that verifier
		works.
	*/
	// TODO(rjkroege): fix this per above.
//	filtered_kv := FilterExtrakeys(extrakeys)

	/*
		TODO(rjkroege): Are here. Decide how to get the article type.	
		I could make it part of the toplevel tags. That sort of pollutes the
		top-level namespace? In particular: journal, book are semantically
		overlapped with the bibtex entry types. 

		I could have bibtex/book bibtex/article etc.?

		I could have @book @article

		I could have @book @bibtex-article. Which defaults to @book @bibtex-book. Hm. I like it.
	*/
//	entrytype, err := ExtractBibTeXEntryType(tags)

/*
	look up the validate table and prove that every required tag for type is in the
	the filtered list. The lists are sort of short. Because you've typed them. But...
*/
//	err = VerifyRequiredFields(entrytype, filtered)

/*
	Use the entry type to choose a template. Generate 

	// TODO(rjkroege): Create templates.
	f := template.Must(template.New("bibtex").Parse(bibtex_template[entrytype])
	b := new(bytes.Buffer)

	// TODO(rjkroege): Learn how to use a map in a template. Must be possible...
	f.Execute(b, filtered)
	// TODO(rjkroege): make right
	return b.toString();

*/



//	
//}

