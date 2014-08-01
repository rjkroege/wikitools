package wiki

import (
	"github.com/rjkroege/wikitools/bibtex"
	"log"
)

func Picktemplate(args []string, tags []string) (tm string, oargs []string, otags []string) {
	templatemap := map[string]string{
		"journal": journaltmpl,
		"book":    booktmpl,
		"article": articletmpl,
	}

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
			return tm, args, tags
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
		otags = append(tags, s)
		oargs = args[1:]
		return
	}
	log.Fatal("No tag or first argument selecting a journal type")
	return
}

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
