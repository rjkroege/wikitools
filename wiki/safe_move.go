package wiki

import (
	"fmt"
	"os"
	"path/filepath"
)

// SafeMoveFile safely moves oldpath to newpath where oldpath and newpath
// are both absolute.
func SafeMoveFile(oldpath, newpath string) error {
	if err := os.MkdirAll(filepath.Dir(newpath), 0700); err != nil {
		return fmt.Errorf("can't mkdir %s because: %v", filepath.Dir(newpath), err)
	}

	if err := os.Link(oldpath, newpath); err != nil {
		return fmt.Errorf("can't link %s to %s because %v", oldpath, newpath, err)
	}

	if err := os.Remove(oldpath); err != nil {
		return fmt.Errorf("can't remove %s because %v", oldpath, err)
	}
	return nil
}
