package wiki

import (
	"os"
	"path/filepath"
	"strings"
)


// IsWikiArticle returns true for files that we don't want to process
// goes in wiki
func IsWikiArticle(path string, info os.FileInfo) bool {
	relp, err := filepath.Rel(Basepath, path)
	if err != nil {
		return true // Always skip bad paths
	}

	switch {
	case info.IsDir():
		return true
	case filepath.Ext(info.Name()) != ".md":
		return true
	case strings.HasPrefix(relp, "templates"):
		return true
	case strings.HasPrefix(relp, "generated"):
		return true
	case info.Name() == "README.md":
		return true
	}
	return false
}
