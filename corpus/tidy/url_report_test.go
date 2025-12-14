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
		name        string
		fpath       string
		passedErr   error
		wantErr     bool
		wantOutUrls []string
		wantForward []string
		wantBack    []string
		wantDamaged []string
	}{
		{
			name:        "empty file",
			fpath:       "../testdata/wiki/unsorted/Saturday.md",
			wantOutUrls: []string{},
			wantForward: []string{},
			wantBack:    []string{},
			wantDamaged: []string{},
		},
		// TODO(rjk): test the pass through error case (should fire an error)
		{
			name:  "file with external URL",
			fpath: "../testdata/wiki/2023/05-May/6/Saturday.md",
			wantOutUrls: []string{
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/05-May/6/Saturday.md",
				"[link](https://example.com)",
			},
			wantForward: []string{},
			wantBack:    []string{},
			wantDamaged: []string{},
		},
		{
			name:  "file with multiple external URLs",
			fpath: "../testdata/wiki/unsorted/EveningJournal.md",
			wantOutUrls: []string{
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/05-May/6/Saturday.md",
				"[link](https://example.com)",
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/unsorted/EveningJournal.md",
				"[Example](https://example.com)[Google](https://google.com)",
			},
			wantForward: []string{},
			wantBack:    []string{},
			wantDamaged: []string{},
		},
		{
			name:        "file with damaged wikilink",
			fpath:       "../testdata/wiki/2023/02-Feb/28/Saturday.md",
			wantForward: []string{},
			wantBack:    []string{},
			wantOutUrls: []string{
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/05-May/6/Saturday.md",
				"[link](https://example.com)",
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/unsorted/EveningJournal.md",
				"[Example](https://example.com)[Google](https://google.com)",
			},
			wantDamaged: []string{
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/02-Feb/28/Saturday.md",
				"[[nonexistentArticle]]",
			},
		},
		{
			name:  "file with internal links",
			fpath: "../testdata/wiki/2023/10-Oct/1/PlottingTools.md",
			wantOutUrls: []string{
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/05-May/6/Saturday.md",
				"[link](https://example.com)",
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/unsorted/EveningJournal.md",
				"[Example](https://example.com)[Google](https://google.com)",
			},
			wantForward: []string{
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/10-Oct/1/PlottingTools.md",
				"[[28/Saturday]][[EveningJournal]]",
			},

			wantBack: []string{
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/02-Feb/28/Saturday.md",
				"[[10-Oct/1/PlottingTools.md]]",
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/unsorted/EveningJournal.md",
				"[[10-Oct/1/PlottingTools.md]]",
			},

			wantDamaged: []string{
				"/Users/rjkroege/tools/wikitools/corpus/testdata/wiki/2023/02-Feb/28/Saturday.md",
				"[[nonexistentArticle]]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fpath, err := filepath.Abs(tt.fpath)
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
			if diff := cmp.Diff(tt.wantBack, links.StringVector(lnks.BackLinks)); diff != "" {
				t.Errorf("BackLinks mismatch (-want +got):\n%s", diff)
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
