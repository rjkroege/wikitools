package links

import (
	"testing"
	"path/filepath"
	
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/corpus/search"
)

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
		mapper: mapper,
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
			links.AddForwardUrl(tt.displaytext, tt.url, tt.fpath)

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
		rwikitext    string
		fpath       string
		bpath	string
	}{
		{
			name:        "Add new wikilink",
			displaytext: "",
			wikitext:    "16/Decisions.md",
			rwikitext:	"unsorted/Saturday.md",
			fpath:       "../testdata/wiki/unsorted/Saturday.md",
			bpath:	"../testdata/wiki/2023/08-Aug/16/Decisions.md",
		},
		{
			name:        "Add another wikilink",
			displaytext: "",
			wikitext:    "16/Decisions.md",
			rwikitext:	"6/Saturday.md",
			fpath:       "../testdata/wiki/2023/05-May/6/Saturday.md",
			bpath:	"../testdata/wiki/2023/08-Aug/16/Decisions.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fpath, err := filepath.Abs( tt.fpath)
			if err != nil {
				t.Fatalf("can't abs %q: %v", tt.fpath, err)
			}
			bpath, err := filepath.Abs( tt.bpath)
			if err != nil {
				t.Fatalf("can't abs %q: %v", tt.bpath, err)
			}
		
			links.AddWikilink(tt.displaytext,tt.wikitext, fpath)

			// Check if the wikilink was added to the correct file's ForwardLinks
			if _, ok := links.ForwardLinks[fpath]; !ok {
				t.Errorf("Expected ForwardLinks to contain file %s", tt.fpath)
			}

			// Check if the specific wikilink was added
			furlref := corpus.MakeWikilink(tt.wikitext, tt.displaytext)
			if _, ok := links.ForwardLinks[fpath][furlref]; !ok {
				t.Errorf("Expected ForwardLinks[%s] to contain wikilink %v", tt.fpath, furlref)
			}

			// Check if the specific backlink was added.
			burlref := corpus.MakeWikilink(tt.rwikitext, "")
			if _, ok := links.BackLinks[bpath][burlref]; !ok {
				t.Errorf("Expected BackLinks[%s] to contain wikilink %v", tt.bpath, burlref)
			}

			// Note: Testing DamagedLinks would require a more complex setup
			// with a functional mapper that can resolve paths and generate backreferences.
		})
	}
}
