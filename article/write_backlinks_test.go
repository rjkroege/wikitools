package article

import (
	"testing"
	"path/filepath"
	"os"

	"github.com/rjkroege/wikitools/corpus"
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

	// TODO(rjk): Decode and compare
		
}

