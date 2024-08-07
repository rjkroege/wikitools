package wiki

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/rjkroege/wikitools/bibtex"
)

type Template struct {
	Template   string
	Custombody string
}

type TemplatePalette map[string]Template

func NewTemplatePalette() TemplatePalette {
	return map[string]Template{
		"journalam": {basetmpl, journalamtmpl},
		"journalpm": {basetmpl, journalpmtmpl},
		"entry":     {basetmpl, entrytmpl},
		"book":      {booktmpl, ""},
		"article":   {articletmpl, ""},
		"code":      {basetmpl, codetmpl},
	}
}

// AddDynamicTemplate aguments the provided TemplatePalette
// with templates read from disk. Errors are inlined into the template
// body to help the user figure out why there was a problem.
func (tm TemplatePalette) AddDynamcTemplates(config map[string]string) {
	for k, v := range config {
		if _, ok := tm[k]; !ok {
			tm[k] = Template{
				Template:   basetmpl,
				Custombody: "",
			}
		}

		fd, err := os.Open(v)
		if err != nil {
			tm[k] = Template{
				Template:   tm[k].Template,
				Custombody: fmt.Sprintf("File %s for key %s had error: %v", v, k, err),
			}
			continue
		}
		byteslice, err := ioutil.ReadAll(fd)
		if err != nil {
			tm[k] = Template{
				Template:   tm[k].Template,
				Custombody: fmt.Sprintf("File %s for key %s had error: %v", v, k, err),
			}
			continue
		}
		tm[k] = Template{
			Template:   tm[k].Template,
			Custombody: string(byteslice),
		}
	}
}

// Picktemplate chooses a template for a new wiki entry bases on the
// provided arguments and tags. Tags are bare here with the prefix '[#@]'
// removed.
func (templatemap TemplatePalette) Picktemplate(args []string, tags []string) (Template, []string, []string) {
	// Handle book/article entries
	// TODO(rjk): Are book formats right?
	booktype, err := bibtex.ExtractBibTeXEntryType(tags)
	if err == nil {
		tm, ok := templatemap[booktype]
		if ok {
			return tm, args, tags
		} else {
			// TODO(rjk): Return errors
			log.Fatal("Selected a BibTex entry type without a matching template")
		}
	}

	// TODO(rjk): Perhaps I should error check? The highroad scheme is to
	// create a filesystem *inside* the wikitools that contains the default
	// template.
	tm := templatemap["entry"]

	modifiedtags := make([]string, 0, len(tags))
	for _, t := range tags {
		if t == "journal" {
			tg := journalfortime(t)
			modifiedtags = append(modifiedtags, tg)
			if ntm, ok := templatemap[tg]; ok {
				tm = ntm
			}
			continue
		}
		modifiedtags = append(modifiedtags, t)
	}

	// Otherwise, it's an entry
	return tm, args, modifiedtags
}

// Split divides the provided arguments into those that wil serve as tags
// on the journal entry and those that are conventional arguments.
func Split(all []string) (args []string, tags []string) {
	for _, v := range all {
		if v[0] == '@' || v[0] == '#' {
			if len(v) > 1 {
				// Bare # or @ characters are discarded.
				tags = append(tags, v[1:])
			}
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
			return "journalam"
		} else {
			return "journalpm"
		}
	}
	return tm
}
