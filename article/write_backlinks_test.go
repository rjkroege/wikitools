package article

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rjkroege/wikitools/corpus"
	// "github.com/google/go-cmp/cmp/cmpopts"
)

func TestWriteBacklinks(t *testing.T) {
	td := t.TempDir()

	fpath := filepath.Join(td, "testfile")

	wikilink := corpus.MakeWikilink("jiminy", "cricket")
	backmap := make(map[corpus.Wikilink]corpus.Empty)
	backmap[wikilink] = corpus.Empty{}

	if err := os.WriteFile(fpath, []byte("hi there"), 0600); err != nil {
		t.Fatalf("can't write to %q: %v", fpath, err)
	}

	if err := WriteBacklinks(fpath, backmap); err != nil {
		t.Fatalf("can't write xattr to %q: %v", fpath, err)
	}

	gotmap, err := ReadBacklinks(fpath)
	if err != nil {
		t.Fatalf("can't read xattr from %q: %v", fpath, err)
	}

	if diff := cmp.Diff(backmap, gotmap); diff != "" {
		t.Errorf("[%d] dump mismatch (-want +got):\n%s", 0, diff)
	}

	wikilink2 := corpus.MakeWikilink("hello", "xx/foo/bar")
	backmap[wikilink2] = corpus.Empty{}

	if err := WriteBacklinks(fpath, backmap); err != nil {
		t.Fatalf("can't second write xattr to %q: %v", fpath, err)
	}

	gotmap, err = ReadBacklinks(fpath)
	if err != nil {
		t.Fatalf("can't read xattr from %q: %v", fpath, err)
	}

	if diff := cmp.Diff(backmap, gotmap); diff != "" {
		t.Errorf("[%d] dump mismatch (-want +got):\n%s", 1, diff)
	}

}
