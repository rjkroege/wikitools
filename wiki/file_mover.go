package wiki

import (
	"fmt"
	"os"
)

func SafeReplaceFile(newpath, oldpath string) error {
	backup := oldpath + ".back"
	if err := os.Link(oldpath, backup); err != nil {
		return fmt.Errorf("replaceFile backup: %v", err)
	}

	if err := os.Remove(oldpath); err != nil {
		return fmt.Errorf("replaceFile remove: %v", err)
	}

	if err := os.Link(newpath, oldpath); err != nil {
		return fmt.Errorf("replaceFile emplace: %v", err)
	}

	if err := os.Remove(newpath); err != nil {
		return fmt.Errorf("replaceFile remove: %v", err)
	}
	if err := os.Remove(backup); err != nil {
		return fmt.Errorf("replaceFile remove: %v", err)
	}
	return nil
}


