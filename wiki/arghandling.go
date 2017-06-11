package wiki

import (
	"log"
	"time"

	"github.com/rjkroege/wikitools/bibtex"
)

// Picktemplate chooses a template for a new wiki entry bases on the provided
// arguments and tags.
func Picktemplate(args []string, tags []string) (tm string, oargs []string, otags []string) {
	templatemap := map[string]string{
		"journal": "journal",
		"entry":   entrytmpl,
		"book":    booktmpl,
		"article": articletmpl,
		"code":    codetmpl,
	}

	// Handle book/article entries
	booktype, err := bibtex.ExtractBibTeXEntryType(tags)
	if err == nil {
		tm, ok := templatemap[booktype]
		if ok {
			return tm, args, tags
		} else {
			log.Fatal("Selected a BibTex entry type without a matching template")
		}
	}

	for _, v := range tags {
		tm, ok := templatemap[v[1:]]
		if ok {
			return journalfortime(tm), args, tags
		}
	}

	// If we do not have a @tag that is choosing a journal format, we use the first
	// argument and it becomes a tag.
	if len(args) < 1 {
		log.Fatal("No candidate argument to specify a t enough arguments\n")
	}

	tm, ok := templatemap[args[0]]
	if ok {
		s := "@" + args[0]
		return journalfortime(tm), args[1:], append(tags, s)
	}
	log.Fatal("No tag or first argument selecting a journal type")
	return
}

// Split divides the provided arguments into those that wil serve as tags
// on the journal entry and those that are conventional arguments.
func Split(all []string) (args []string, tags []string) {
	for _, v := range all {
		if v[0] == '@' {
			tags = append(tags, v)
		} else {
			args = append(args, v)
		}
	}
	return
}

// For mockability.
var journaltimepicker = BeforeNoon

// BeforeNoon is a convenience function that determines if
// the current time is before or after noon.
func BeforeNoon() bool {
	now := time.Now()
	return now.Hour() < 12
}

// journalfortime adjusts the journal template based on time of day.
func journalfortime(tm string) string {
	if tm == "journal" {
		if journaltimepicker() {
			tm = journalamtmpl
		} else {
			tm = journalpmtmpl
		}
	}
	return tm
}
