package search

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

/*
/Users/rjkroege/Documents/wiki/2013/07-Jul/14/Sunday.md /Users/rjkroege/Documents/wiki/2013/07-Jul/7/Sunday.md /Users/rjkroege/Documents/wiki/2013/04-Apr/7/Sunday.md /Users/rjkroege/Documents/wiki/2013/04-Apr/21/Warmer-Sunday.md
*/

const location = "/wiki"

type teststim struct {
	location string
	lsd      string
	wikitext string
	allpaths []string
	want     string
	wanterr  error
}

func TestDisambiguateWikiPaths(t *testing.T) {
	testtab := []teststim{
		{
			location: location,
			lsd:      location,
			wikitext: "Sunday.md",
			allpaths: []string{},
			want:     "",
			wanterr:  NoFileForWikitext,
		},
		{
			location: location,
			lsd:      "/wiki/2013/07-Jul/14",
			wikitext: "Sunday.md",
			allpaths: []string{
				"/wiki/Sunday.md",
			},
			want:    "/wiki/Sunday.md",
			wanterr: nil,
		},
		{
			location: location,
			lsd:      "/wiki/2013/07-Jul/14",
			wikitext: "Sunday.md",
			allpaths: []string{
				"/wiki/2013/07-Jul/14/Sunday.md",
				"/wiki/2013/07-Jul/13/Sunday.md",
				"/wiki/2013/06-Jun/13/Sunday.md",
				"/wiki/2013/06-Jun/13/puddle/Sunday.md",
			},
			want:    "/wiki/2013/07-Jul/14/Sunday.md",
			wanterr: nil,
		},
		{
			location: location,
			lsd:      "/wiki/2013/06-Jun/14",
			wikitext: "Sunday.md",
			allpaths: []string{
				"/wiki/2013/07-Jul/13/Sunday.md",
				"/wiki/2013/06-Jun/13/Sunday.md",
				"/wiki/2013/06-Jun/14/puddle/Sunday.md",
			},
			want:    "/wiki/2013/06-Jun/14/puddle/Sunday.md",
			wanterr: nil,
		},
		{
			location: location,
			lsd:      "/wiki/2013/07-Jul/14",
			wikitext: "13/Sunday.md",
			allpaths: []string{
				"/wiki/2013/07-Jul/13/Sunday.md",
				"/wiki/2013/06-Jun/13/Sunday.md",
				"/wiki/2013/06-Jun/13/puddle/Sunday.md",
			},
			want:    "",
			wanterr: AmbiguousWikitext,
		},
		{
			location: location,
			lsd:      "/wiki",
			wikitext: "07-Jul/13/Sunday.md",
			allpaths: []string{
				"/wiki/2013/07-Jul/14/Sunday.md",
				"/wiki/2013/07-Jul/13/Sunday.md",
				"/wiki/2013/06-Jun/13/Sunday.md",
				"/wiki/2013/06-Jun/13/puddle/Sunday.md",
			},
			want:    "/wiki/2013/07-Jul/13/Sunday.md",
			wanterr: nil,
		},
	}

	for i, tv := range testtab {
		got, goterr := disambiguatewikipaths(tv.location, tv.lsd, tv.wikitext, tv.allpaths)
		if diff := cmp.Diff(tv.want, got); diff != "" {
			t.Errorf("[%d] dump mismatch (-want +got):\n%s", i, diff)
		}
		if diff := cmp.Diff(tv.wanterr, goterr, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("[%d] error dump mismatch (-want +got):\n%s", i, diff)
		}
	}

}

func TestBuildshortestwikitext(t *testing.T) {
	// Can use the same teststim.
	testtab := []teststim{
		{
			location: "/wiki",
			lsd:      "/wiki/bar",
			allpaths: []string{},
			want:     "",
			wanterr:  NoValidMatch,
		},
		{
			location: "/wiki/",
			lsd:      "/wiki/2013/06-Jun/13/Sunday.md",
			allpaths: []string{
				"/wiki/2013/07-Jul/14/Sunday.md",
				"/wiki/2013/07-Jul/13/Sunday.md",
				"/wiki/2013/06-Jun/13/Sunday.md",
				"/wiki/2013/06-Jun/13/puddle/Sunday.md",
			},
			want: "06-Jun/13/Sunday.md",
		},
		{
			location: "/wiki/",
			lsd:      "/wiki/2013/06-Jun/13/Bombast.md",
			allpaths: []string{
				"/wiki/2013/06-Jun/13/Bombast.md",
			},
			want: "Bombast.md",
		},
		{
			location: "/wiki/",
			lsd:      "/wiki/2013/06-Jun/13/puddle/Sunday.md",
			allpaths: []string{
				"/wiki/2013/07-Jul/14/Sunday.md",
				"/wiki/2013/07-Jul/13/Sunday.md",
				"/wiki/2013/06-Jun/13/Sunday.md",
				"/wiki/2013/06-Jun/13/puddle/Sunday.md",
			},
			want: "puddle/Sunday.md",
		},
		{
			location: "/wiki/",
			lsd:      "/wiki/2013/07-Jul/14/Sunday.md",
			allpaths: []string{
				"/wiki/2013/07-Jul/14/Sunday.md",
				"/wiki/2013/07-Jul/13/Sunday.md",
				"/wiki/2013/06-Jun/13/Sunday.md",
				"/wiki/2013/06-Jun/13/puddle/Sunday.md",
			},
			want: "14/Sunday.md",
		},
	}

	for i, tv := range testtab {
		got, goterr := buildshortestwikitext(tv.location, tv.lsd, tv.allpaths)
		if diff := cmp.Diff(tv.want, got); diff != "" {
			t.Errorf("[%d] dump mismatch (-want +got):\n%s", i, diff)
		}
		if diff := cmp.Diff(tv.wanterr, goterr, cmpopts.EquateErrors()); diff != "" {
			t.Errorf("[%d] error dump mismatch (-want +got):\n%s", i, diff)
		}
	}

}
