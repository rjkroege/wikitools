package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rjkroege/wikitools/wiki"
)

// ArgsAlfredPreprocess reprocesses os.Args iff we are running inside Alfred
// because Alfred app doesn't split apart its args. It just dumps
// everything as a single string and execs the command. This makes a
// (certain kind of) sense but requires running this helper to adjust
// os.Args before running kong. Returns the new os.Args value.
func ArgsAlfredPreprocess(args []string) []string {
	if _, forkedbyalfred := os.LookupEnv("alfred_workflow_uid"); !forkedbyalfred || len(args) < 2 {
		return args
	}

	nwa := make([]string, 0, 8)
	nwa = append(nwa, args[0])
	for i, s := range strings.Split(args[1], " ") {
		// can do something different here
		if i == 0 && s == ActionMarker {
			nwa = append(nwa, "new")
			continue
		}
		if i == 0 && s == PlumbMarker {
			nwa = append(nwa, "plumb")
			continue
		}

		if i == 0 {
			nwa = append(nwa, "newautocomplete")
		}
		if len(s) > 0 {
			nwa = append(nwa, s)
		}
	}

	return nwa
}

func WikinewAutocomplete(settings *wiki.Settings, args []string) {
	//	log.Println("WikinewAutocomplete")
	var result Result

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "	")

	// Alfred doesn't tokenize or divide its argument when they're passed
	// directly. So I have to split them up. I also presume that they might
	// have been already split so that I can execute this code from the
	// shell.
	splitargs, pargs := wiki.Split(args)

	//	log.Println("args:", splitargs, pargs)

	if len(pargs) > 0 { // There are tags. These might auto-completable.
		// TODO(rjk): more of these values should be configurable and in settings?
		tagpath := filepath.Join(settings.Wikidir, wiki.Reportpath, "taglist")

		tagfile, err := os.ReadFile(tagpath)
		if err != nil {
			log.Printf("can't read tags from %q: %v", tagpath, err)
		}
		tags := strings.Split(string(tagfile), "\n")

		th := make(map[string]int)

		// TODO(rjk): This could written to run in a faster way.
		for _, pr := range pargs {
			for _, t := range tags {
				if strings.HasPrefix(t, pr) && pr != t {
					if s, ok := th[t]; !ok || s < len(pr) {
						th[t] = len(pr)
					}
				}
			}
		}

		// log.Println(th)
		for t, v := range th {
			var sb strings.Builder
			for i, a := range splitargs {
				if i > 0 {
					sb.WriteRune(' ')
				}
				sb.WriteString(a)
			}

			for _, a := range pargs {
				sb.WriteString(" ")
				if strings.HasPrefix(t, a) && t != a {
					sb.WriteRune('@')
					sb.WriteString(t)
				} else {
					sb.WriteRune('@')
					sb.WriteString(a)
				}
			}

			// Add the tags that aren't an auto-completion

			// Exclude the Uid field to make sure that the items aren't re-ordered.
			result.Items = append(result.Items, &Item{
				Title:        sb.String(),
				Arg:          sb.String(),
				Autocomplete: sb.String() + " ",
				relevance:    v,
			})
		}
		sort.Sort(result.Items)
	}

	// Alfred requires a non-empty Item to offer it in the list. So we create
	// one that we can pass downstream. The downstream (e.g. action handler)
	// will then take different action based on the presence of the flag.

	var sb strings.Builder
	for i, s := range splitargs {
		if i > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteString(s)
	}
	for _, s := range pargs {
		sb.WriteRune(' ')
		sb.WriteRune('@')
		sb.WriteString(s)
	}

	fs := sb.String()
	if fs != "" {
		result.Items = append(result.Items, &Item{
			Title: fs,
			Arg:   ActionMarker + " " + fs,
		})
	}
	if err := encoder.Encode(result); err != nil {
		log.Fatalf("can't write json %v", err)
	}
}

const ActionMarker = "#actionphase"
const PlumbMarker = "^plumbphase"

type Item struct {
	Uid          string `json:"uid,omitempty"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle,omitempty"`
	Arg          string `json:"arg"`
	Autocomplete string `json:"autocomplete"`
	relevance    int
}

type Result struct {
	Items ItemCollection `json:"items"`
}

type ItemCollection []*Item

func (c ItemCollection) Len() int {
	return len(c)
}

func (c ItemCollection) Less(i, j int) bool {
	return c[i].relevance > c[j].relevance
}

func (c ItemCollection) Swap(i, j int) {
	tmp := c[i]
	c[i] = c[j]
	c[j] = tmp
}

var _ = ItemCollection(nil)
