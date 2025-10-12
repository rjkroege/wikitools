package wiki

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSafeMoveFile(t *testing.T) {
	// create a temp directory for test files
	tmp := t.TempDir()

	// create source file
	old := filepath.Join(tmp, "src.txt")
	if err := ioutil.WriteFile(old, []byte("data"), 0600); err != nil {
		t.Fatal(err)
	}

	// destination inside a new sub-directory
	new := filepath.Join(tmp, "subdir", "dst.txt")

	// perform the move
	if err := SafeMoveFile(old, new); err != nil {
		t.Fatalf("SafeMoveFile failed: %v", err)
	}

	// verify old file is gone
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Errorf("old file still exists: %v", err)
	}

	// verify new file exists with correct content
	got, err := ioutil.ReadFile(new)
	if err != nil {
		t.Fatalf("cannot read new file: %v", err)
	}
	if string(got) != "data" {
		t.Errorf("file content mismatch: got %q, want %q", got, "data")
	}
}

func TestSafeMoveFileErrors(t *testing.T) {
	tmp := t.TempDir()

	t.Run("nonexistent source", func(t *testing.T) {
		err := SafeMoveFile(filepath.Join(tmp, "no-such-file"), filepath.Join(tmp, "dst"))
		if err == nil {
			t.Error("expected error for missing source, got nil")
		}
	})

	t.Run("empty paths", func(t *testing.T) {
		err := SafeMoveFile("", "")
		if err == nil {
			t.Error("expected error for empty paths, got nil")
		}
	})
}
