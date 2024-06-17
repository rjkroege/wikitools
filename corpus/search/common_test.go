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
	lsd string
	wikitext string
	allpaths []string
	want string
	wanterr error
}

func TestDisambiguateWikiPaths(t *testing.T) {
	testtab  := []teststim {
		{
			location: location,
			lsd: "2013/07-Jul/14",
			wikitext: "Sunday.md",
			allpaths: []string{
			},
			want: "",
			wanterr: NoFileForWikitext,
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