package wiki

// TODO(rjk): Should be in the "config" directory.

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Toplevel settings.
type Settings struct {
	Wikidir        string            `json:"wikidir"`
	TemplateForTag map[string]string `json:"templatefortag"`
}

// Read opens a json format configuration file.
func Read(path string) (*Settings, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no config file %q: %v", path, err)
	}

	settings := &Settings{}
	decoder := json.NewDecoder(fd)
	if err := decoder.Decode(settings); err != nil {
		return nil, fmt.Errorf("error parsing config %q: %v", path, err)
	}

	// TODO(rjk): Validate the configurable settings.
	return settings, nil
}

// Inquiry functions based on configurable values.

// IsWikiArticle returns true for files that we don't want to process
// goes in wiki
// TODO(rjk): This is misnamed. Maybe call it "NotArticle"
// TODO(rjk): Flip its sense too?
func (s *Settings) IsWikiArticle(path string, info os.FileInfo) bool {
	relp, err := filepath.Rel(s.Wikidir, path)
	if err != nil {
		return true // Always skip bad paths
	}

	// TODO(rjk): parametrize it all.
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

// TODO(rjk): Revisit if these are right or not.

// isWikiLink returns true if the provided dest is a link inside of the
// wiki. Links are "inside" the wiki if they are relative or absolute
// with the root of the wiki as prefix.
// TODO(rjk): Should I make sure that there's a file at the end of the
// link? wikipp shoudln't but wikiclean should probably check link
// validity for all of the wiki articles and generate a report if they
// contain invalid links.
func (s *Settings) IsWikiLink(dest []byte) bool {
	pth := path.Clean(string(dest))
	if !path.IsAbs(pth) || strings.HasPrefix(pth, s.Wikidir) {
		return true
	}
	return false
}

func (s *Settings) IsWikiMarkdownLink(dest []byte) bool {
	pth := path.Clean(string(dest))
	// TODO(rjk): Consider making the extension configurable.
	if path.Ext(pth) == Extension && (!path.IsAbs(pth) || strings.HasPrefix(pth, s.Wikidir)) {
		return true
	}
	return false
}
