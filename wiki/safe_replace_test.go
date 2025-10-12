package wiki

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSafeReplaceFile(t *testing.T) {
	// Create a temporary directory for test files.
	tmpDir, err := ioutil.TempDir("", "wiki_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Paths for test files.
	oldFile := filepath.Join(tmpDir, "old.txt")
	newFile := filepath.Join(tmpDir, "new.txt")

	// Create the old file.
	if err := ioutil.WriteFile(oldFile, []byte("old content"), 0644); err != nil {
		t.Fatalf("failed to write old file: %v", err)
	}

	// Create the new file.
	if err := ioutil.WriteFile(newFile, []byte("new content"), 0644); err != nil {
		t.Fatalf("failed to write new file: %v", err)
	}

	// Run the function.
	if err := SafeReplaceFile(newFile, oldFile); err != nil {
		t.Fatalf("SafeReplaceFile failed: %v", err)
	}

	// Verify old file now contains new content.
	got, err := ioutil.ReadFile(oldFile)
	if err != nil {
		t.Fatalf("failed to read old file: %v", err)
	}
	if string(got) != "new content" {
		t.Errorf("old file has wrong content: %s", got)
	}

	// Verify new file and backup are removed.
	if _, err := os.Stat(newFile); !os.IsNotExist(err) {
		t.Errorf("new file still exists")
	}
	backup := oldFile + ".back"
	if _, err := os.Stat(backup); !os.IsNotExist(err) {
		t.Errorf("backup file still exists")
	}
}
