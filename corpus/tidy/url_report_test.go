package tidy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rjkroege/wikitools/corpus/links"
	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
)


func Test_onefileimpl(t *testing.T) {
	wikiroot, err := filepath.Abs("../testdata")
	if err != nil {
		t.Fatalf("test can't run: %v", err)
	}
	mapper := search.MakeWikilinkNameIndex(wikiroot)

			settings := &wiki.Settings{
				Wikidir: wikiroot,
			}
			lnks := links.MakeLinks(mapper, wikiroot)

	tests := []struct {
		name         string
		fpath string
		passedErr error
		wantErr bool
		wantOutUrls  []string
		wantForward  []string
		wantDamaged  []string
	}{
		{
			name:    "empty file",
			fpath: "../testdata/wiki/unsorted/Saturday.md",
			wantOutUrls: []string{},
			wantForward: []string{},
			wantDamaged: []string{},
		},
// TODO(rjk): test the pass through error case (should fire an error)
		{
			name: "file with external URL",
			fpath: "../testdata/wiki/2023/05-May/6/Saturday.md",
			wantOutUrls: []string{
"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/05-May/6/Saturday.md",
"[link](https://example.com)",

			},
			wantForward: []string{},
			wantDamaged: []string{},
		},
// 		{
// 			name: "file with multiple external URLs",
// 			content: `# Multiple Links
// 
// - [Example](https://example.com)
// - [Google](https://google.com)
// `,
// 			wantErr: false,
// 			wantOutUrls: map[string]corpus.LinkMap[corpus.Urllink]{
// 				"FILEPATH": {
// 					corpus.MakeUrllink("https://example.com", "Example"): {},
// 					corpus.MakeUrllink("https://google.com", "Google"):   {},
// 				},
// 			},
// 			wantForward: map[string]corpus.LinkMap[corpus.Wikilink]{},
// 			wantDamaged: map[string]corpus.LinkMap[corpus.Wikilink]{},
// 		},
// 		{
// 			name: "file with damaged wikilink",
// 			content: `# Wikilink Test
// 
// Here is a [[nonexistent article]] wikilink.
// `,
// 			wantErr: false,
// 			wantOutUrls: map[string]corpus.LinkMap[corpus.Urllink]{},
// 			wantForward: map[string]corpus.LinkMap[corpus.Wikilink]{},
// 			wantDamaged: map[string]corpus.LinkMap[corpus.Wikilink]{
// 				"FILEPATH": {
// 					corpus.MakeWikilink("nonexistent article.md", ""): {},
// 				},
// 			},
// 		},
// 		{
// 			name: "file with metadata header",
// 			content: `---
// title: Test Article
// tags: test
// ---
// 
// # Test Article
// 
// Content with a [link](https://example.org).
// `,
// 			wantErr: false,
// 			wantOutUrls: map[string]corpus.LinkMap[corpus.Urllink]{
// 				"FILEPATH": {
// 					corpus.MakeUrllink("https://example.org", "link"): {},
// 				},
// 			},
// 			wantForward: map[string]corpus.LinkMap[corpus.Wikilink]{},
// 			wantDamaged: map[string]corpus.LinkMap[corpus.Wikilink]{},
// 		},
// 		{
// 			name: "file with mixed links",
// 			content: `# Mixed Links
// 
// External: [Example](https://example.com)
// Internal: [[some page]]
// `,
// 			wantErr: false,
// 			wantOutUrls: map[string]corpus.LinkMap[corpus.Urllink]{
// 				"FILEPATH": {
// 					corpus.MakeUrllink("https://example.com", "Example"): {},
// 				},
// 			},
// 			wantForward: map[string]corpus.LinkMap[corpus.Wikilink]{},
// 			wantDamaged: map[string]corpus.LinkMap[corpus.Wikilink]{
// 				"FILEPATH": {
// 					corpus.MakeWikilink("some page.md", ""): {},
// 				},
// 			},
// 		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fpath, err := filepath.Abs( tt.fpath)
			if err != nil {
				t.Fatalf("can't abs %q: %v", tt.fpath, err)
			}

			info, err := os.Stat(fpath)
			if err != nil {
				t.Fatalf("couldn't stat %q just made: %v", fpath, err)
			}

			err = onefileimpl(settings, lnks, fpath, info, tt.passedErr)
			if err != nil && !tt.wantErr {
				t.Errorf("onefileimpl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err == nil {
				t.Errorf("onefileimpl() failed to err when expected")
				return
			}

			if diff := cmp.Diff(tt.wantOutUrls, links.StringVector(lnks.OutUrls)); diff != "" {
				t.Errorf("OutUrls mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantForward, links.StringVector(lnks.ForwardLinks)); diff != "" {
				t.Errorf("ForwardLinks mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantDamaged, links.StringVector(lnks.DamagedLinks)); diff != "" {
				t.Errorf("DamagedLinks mismatch (-want +got):\n%s", diff)
			}
		})
	}
}


func Test_onefileimpl_NonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	settings := &wiki.Settings{
		Wikidir: tmpDir,
	}
	lnks := links.MakeLinks(search.MakeWikilinkNameIndex(tmpDir), tmpDir)

	err := onefileimpl(settings, lnks, filepath.Join(tmpDir, "nonexistent.md"), nil, nil)
	if err == nil {
		t.Error("onefileimpl() should return error for nonexistent file")
	}
}


	