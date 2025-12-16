package links

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/corpus/search"
)

func fxpth(wikiroot, relpath string) string {
	return filepath.Join(wikiroot, relpath)
}

// getMapper returns a mapper instance for testing.
func getMapper(t *testing.T) corpus.LinkToFile {

	bp, err := filepath.Abs("../testdata")
	if err != nil {
		t.Fatalf("test can't run: %v", err)
	}
	return search.MakeWikilinkNameIndex(bp)
}

// TestAddForwardUrl tests the AddForwardUrl function.
func TestAddForwardUrl(t *testing.T) {
	mapper := getMapper(t)
	links := &Links{
		mapper:  mapper,
		OutUrls: make(map[string]corpus.UrlMap),
	}

	tests := []struct {
		name        string
		displaytext string
		url         string
		fpath       string
	}{
		{
			name:        "Add new forward URL",
			displaytext: "Example",
			url:         "https://example.com",
			fpath:       "test.md",
		},
		{
			name:        "Add another forward URL to the same file",
			displaytext: "Another Example",
			url:         "https://another.example.com",
			fpath:       "test.md",
		},
		{
			name:        "Add forward URL to a different file",
			displaytext: "Different File",
			url:         "https://different.example.com",
			fpath:       "different.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links.addForwardUrl(tt.displaytext, tt.url, tt.fpath)

			// Check if the URL was added to the correct file's OutUrls
			if _, ok := links.OutUrls[tt.fpath]; !ok {
				t.Errorf("Expected OutUrls to contain file %s", tt.fpath)
			}

			// Check if the specific URL link was added
			urlref := corpus.MakeUrllink(tt.url, tt.displaytext)
			if _, ok := links.OutUrls[tt.fpath][urlref]; !ok {
				t.Errorf("Expected OutUrls[%s] to contain URL link %v", tt.fpath, urlref)
			}
		})
	}
}

// TestAddWikilink tests the AddWikilink function.
// TODO(Actually exercise the description part of the wikitext)
func TestAddWikilink(t *testing.T) {
	wikiroot, err := filepath.Abs("../testdata")
	if err != nil {
		t.Fatalf("test can't run: %v", err)
	}
	mapper := getMapper(t)
	links := MakeLinks(mapper, wikiroot)

	tests := []struct {
		name        string
		displaytext string
		wikitext    string
		fpath       string

		forwardlinks []string
		backlinks    []string
		outurls      []string
		damagedlinks []string
	}{
		{
			name:        "Add new wikilink",
			displaytext: "",
			wikitext:    "16/Decisions.md",
			fpath:       "../testdata/wiki/unsorted/Saturday.md",

			forwardlinks: []string{fxpth(wikiroot, "../testdata/wiki/unsorted/Saturday.md"), "[[16/Decisions.md]]"},
			backlinks:    []string{fxpth(wikiroot, "../testdata/wiki/2023/08-Aug/16/Decisions.md"), "[[unsorted/Saturday.md]]"},
			outurls:      []string{},
			damagedlinks: []string{},
		},
		{
			name:        "Add another wikilink",
			displaytext: "",
			wikitext:    "16/Decisions.md",
			fpath:       "../testdata/wiki/2023/05-May/6/Saturday.md",

			forwardlinks: []string{fxpth(wikiroot, "../testdata/wiki/2023/05-May/6/Saturday.md"), "[[16/Decisions.md]]", fxpth(wikiroot, "../testdata/wiki/unsorted/Saturday.md"), "[[16/Decisions.md]]"},
			backlinks:    []string{fxpth(wikiroot, "../testdata/wiki/2023/08-Aug/16/Decisions.md"), "[[6/Saturday.md]][[unsorted/Saturday.md]]"},
			outurls:      []string{},
			damagedlinks: []string{},
		},
		{
			name:        "Add failing wikilink",
			displaytext: "",
			wikitext:    "Decisions.md",
			fpath:       "../testdata/wiki/2023/05-May/6/Saturday.md",

			forwardlinks: []string{fxpth(wikiroot, "../testdata/wiki/2023/05-May/6/Saturday.md"), "[[16/Decisions.md]]", fxpth(wikiroot, "../testdata/wiki/unsorted/Saturday.md"), "[[16/Decisions.md]]"},
			backlinks:    []string{fxpth(wikiroot, "../testdata/wiki/2023/08-Aug/16/Decisions.md"), "[[6/Saturday.md]][[unsorted/Saturday.md]]"},
			outurls:      []string{},
			damagedlinks: []string{fxpth(wikiroot, "../testdata/wiki/2023/05-May/6/Saturday.md"), "[[Decisions.md]]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fpath, err := filepath.Abs(tt.fpath)
			if err != nil {
				t.Fatalf("can't abs %q: %v", tt.fpath, err)
			}
			
			wr := corpus.MakeWikilink(tt.wikitext, tt.displaytext)
			links.addWikilink(wr , fpath)

			// Dump for diagnostics.
			// 			t.Logf("dump it links\nForwardLinks\n%s\nBackLinks\n%s\nDamagedLinks\n%s\n",
			// 				lstring(links.ForwardLinks),
			// 				lstring(links.BackLinks),
			// 				lstring(links.DamagedLinks))

			if diff := cmp.Diff(StringVector(links.ForwardLinks), tt.forwardlinks); diff != "" {
				t.Errorf("ForwardLinks dump mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(StringVector(links.BackLinks), tt.backlinks); diff != "" {
				t.Errorf("BackLinks dump mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(StringVector(links.DamagedLinks), tt.damagedlinks); diff != "" {
				t.Errorf("DamagedLinks dump mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
