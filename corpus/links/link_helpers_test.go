package links

import (
	"slices"
	"testing"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/google/go-cmp/cmp"
)

func TestBackLinksIterator(t *testing.T) {
	// Create a Links instance with some dummy data
	testLinks := &Links{
		BackLinks: map[string]corpus.LinkMap[corpus.Wikilink]{
			"/tmp/wiki/pageB": map[corpus.Wikilink]empty{
				{ Id: "pageA" }: struct{}{},
				{ Id: "pageC" }: struct{}{},
			},
			"/tmp/wiki/pageC": map[corpus.Wikilink]empty{
				{ Id: "pageA" }: struct{}{},
			},
		},
	}

	// Iterate using BackLinksIterator
	got := make([]string, 0)
	for tuple := range testLinks.BackLinksIterator() {
		got = append(got, tuple.From + "→" + corpus.Stringify(tuple.To))
	}

	slices.Sort(got)

	want := []string{
		"/tmp/wiki/pageB→[[pageA]][[pageC]]",
		"/tmp/wiki/pageC→[[pageA]]",
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("dump mismatch (-want +got):\n%s",  diff)
	}
}