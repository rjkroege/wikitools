package links

import (
	"testing"
	"path/filepath"
	
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/corpus/search"
)

// getMapper returns a mapper instance for testing.
func getMapper(t *testing.T) corpus.LinkToFile {
	
	bp, err := filepath.Abs("testdata")
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

/*
// TestAddWikilink tests the AddWikilink function.
func TestAddWikilink(t *testing.T) {
	mapper := getMapper()
	links := &Links{
		mapper: mapper,
		ForwardLinks: make(map[string]map[Wikilink]Empty),
		BackLinks:    make(map[string]map[Wikilink]Empty),
		DamagedLinks: make(map[string]map[Wikilink]Empty),
		location:     "/test/location",
	}

	tests := []struct {
		name        string
		displaytext string
		wikitext    string
		fpath       string
	}{
		{
			name:        "Add new wikilink",
			displaytext: "Page1",
			wikitext:    "page1.md",
			fpath:       "test.md",
		},
		{
			name:        "Add another wikilink to the same file",
			displaytext: "Page2",
			wikitext:    "page2.md",
			fpath:       "test.md",
		},
		{
			name:        "Add wikilink to a different file",
			displaytext: "Page3",
			wikitext:    "page3.md",
			fpath:       "different.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links.AddWikilink(tt.displaytext, tt.wikitext, tt.fpath)

			// Check if the wikilink was added to the correct file's ForwardLinks
			if _, ok := links.ForwardLinks[tt.fpath]; !ok {
				t.Errorf("Expected ForwardLinks to contain file %s", tt.fpath)
			}

			// Check if the specific wikilink was added
			urlref := MakeWikilink(tt.wikitext, tt.displaytext)
			if _, ok := links.ForwardLinks[tt.fpath][urlref]; !ok {
				t.Errorf("Expected ForwardLinks[%s] to contain wikilink %v", tt.fpath, urlref)
			}

			// Note: Testing BackLinks and DamagedLinks would require a more complex setup
			// with a functional mapper that can resolve paths and generate backreferences.
			// This basic test focuses on the ForwardLinks addition.
		})
	}
}
*/