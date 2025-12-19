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

func fxonepath(wikiroot, relpath string) string {
	return filepath.Join(wikiroot, relpath)
}

func tempReplicate(t testing.TB, srcDir string) string {
	t.Helper()

	tmp := t.TempDir() // automatically removed when test ends

	if err := ReplicateDir(srcDir, tmp); err != nil {
		t.Fatalf("TempReplicate(%q) failed: %v", srcDir, err)
	}
	return tmp
}

func mapby2(s []string, fn func(string) string) []string {
	for i := 0; i < len(s); i += 2 {
		s[i] = fn(s[i])
	}
	return s
}

func fxpth(paths []string, wikiroot string) []string {
	return mapby2(paths, func(s string) string {
		return fxonepath(wikiroot, s)
	})
}

func Test_onefileimpl(t *testing.T) {
	wikiroot := tempReplicate(t, "../testdata")
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
			name:        "empty_file",
			fpath:       "wiki/unsorted/Saturday.md",
			wantOutUrls: []string{},
			wantForward: []string{},
			wantBack:    []string{},
			wantDamaged: []string{},
		},
		// TODO(rjk): test the pass through error case (should fire an error)
		{
			name:  "file_with_external_URL",
			fpath: "wiki/2023/05-May/6/Saturday.md",
			wantOutUrls: []string{
				"wiki/2023/05-May/6/Saturday.md",
				"[link](https://example.com)",
			},
			wantForward: []string{},
			wantBack:    []string{},
			wantDamaged: []string{},
		},
		{
			name:  "file_with_multiple_external_URLs",
			fpath: "wiki/unsorted/EveningJournal.md",
			wantOutUrls: []string{
				"wiki/2023/05-May/6/Saturday.md",
				"[link](https://example.com)",
				"wiki/unsorted/EveningJournal.md",
				"[Example](https://example.com)[Google](https://google.com)",
			},
			wantForward: []string{},
			wantBack:    []string{},
			wantDamaged: []string{},
		},
		{
			name:        "file_with_damaged_wikilink",
			fpath:       "wiki/2023/02-Feb/28/Saturday.md",
			wantForward: []string{},
			wantBack:    []string{},
			wantOutUrls: []string{
				"wiki/2023/05-May/6/Saturday.md",
				"[link](https://example.com)",
				"wiki/unsorted/EveningJournal.md",
				"[Example](https://example.com)[Google](https://google.com)",
			},
			wantDamaged: []string{
				"wiki/2023/02-Feb/28/Saturday.md",
				"[[nonexistentArticle]]",
			},
		},
		{
			name:  "file_with_internal_links",
			fpath: "wiki/2023/10-Oct/1/PlottingTools.md",
			wantOutUrls: []string{
				"wiki/2023/05-May/6/Saturday.md",
				"[link](https://example.com)",
				"wiki/unsorted/EveningJournal.md",
				"[Example](https://example.com)[Google](https://google.com)",
			},
			wantForward: []string{
				"wiki/2023/10-Oct/1/PlottingTools.md",
				"[[28/Saturday]][[EveningJournal]]",
			},

			wantBack: []string{
				"wiki/2023/02-Feb/28/Saturday.md",
				"[[10-Oct/1/PlottingTools.md]]",
				"wiki/unsorted/EveningJournal.md",
				"[[10-Oct/1/PlottingTools.md]]",
			},

			wantDamaged: []string{
				"wiki/2023/02-Feb/28/Saturday.md",
				"[[nonexistentArticle]]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
 
t.Logf("wikiroot: %q", wikiroot)
			fpath := filepath.Join(wikiroot, tt.fpath)
			info, err := os.Stat(fpath)
			if err != nil {
				t.Fatalf("couldn't stat %q just made: %v", fpath, err)
			}
t.Logf("fpath: %q", fpath)

			err = onefileimpl(settings, lnks, fpath, info, tt.passedErr)
			if err != nil && !tt.wantErr {
				t.Errorf("onefileimpl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err == nil {
				t.Errorf("onefileimpl() failed to err when expected")
				return
			}

			if diff := cmp.Diff(fxpth(tt.wantOutUrls,wikiroot), links.StringVector(lnks.OutUrls)); diff != "" {
				t.Errorf("OutUrls mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(fxpth(tt.wantForward,wikiroot), links.StringVector(lnks.ForwardLinks)); diff != "" {
				t.Errorf("ForwardLinks mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(fxpth(tt.wantBack,wikiroot), links.StringVector(lnks.BackLinks)); diff != "" {
				t.Errorf("BackLinks mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(fxpth(tt.wantDamaged,wikiroot), links.StringVector(lnks.DamagedLinks)); diff != "" {
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
