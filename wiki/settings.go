package wiki

// TODO(rjk): Should be in the "config" directory.

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// Toplevel settings.
type Settings struct {
	Wikidir        string            `json:"wikidir"`
	TemplateForTag map[string]string `json:"templatefortag"`
	// TODO(rjk): Consider making the extension configurable.
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

// NotArticle returns true for files that we don't want to process
// goes in wiki
// TODO(rjk): Flip its sense too?
func (s *Settings) NotArticle(path string, info os.FileInfo) bool {
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

var nowfunc = time.Now

const (
	timeformat = "20060102-150405"
)

// UniquingExtension computes a string that will make the filename unique
// in the desired directory. datepath is the relative date specific path, filename
// is the desired filename without an extension.
// Currently includes the leading -.
func (s *Settings) UniquingExtension(datepath, filename string) string {
	fn := s.ExtensionedFileName(filename)

	if _, err := os.Stat(filepath.Join(s.Wikidir, datepath, fn)); err == nil {
		nfn := "-" + nowfunc().Format(timeformat)
		return nfn
	}
	return ""
}

// TODO(rjk): The extension should become configurable.
func (s *Settings) ExtensionedFileName(filename string) string {
	return filename + ".md"
}

func (s *Settings) Extension() string {
	return ".md"
}

// SplitActualDir returns the relative date dir (possibly empty) for a
// given absolute path.
func (s *Settings) SplitActualDir(abspath string) string {
	d := filepath.Dir(abspath)
	rp, err := filepath.Rel(s.Wikidir, d)
	if err != nil {
		log.Fatalf("can't split the dir %#v: %v", abspath, err)
	}
	return rp
}

// SplitActualName divides fn into the preferred name, its (optional)
// uniquing extension (with its leading -) and the filename extension.
func SplitActualName(fn string) (string, string, string) {
	ext := filepath.Ext(fn)
	nexs := fn[0 : len(fn)-len(ext)]
	ps := strings.Split(nexs, "-")

	// It's possible that the last two pieces make a date.
	uniqueing := ""
	if len(ps) > 2 {
		uniqueing = strings.Join([]string{ps[len(ps)-2], ps[len(ps)-1]}, "-")
		if _, err := time.Parse(timeformat, uniqueing); err == nil {
			nexs = strings.Join(ps[0:len(ps)-2], "-")
			uniqueing = "-" + uniqueing
		} else {
			uniqueing = ""
		}
	}
	return nexs, uniqueing, ext
}

func (s *Settings) MakeGenDir() (string, error) {
	// TODO(rjk): Consider adding Reportpath to the configurable state.
	genpath := filepath.Join(s.Wikidir, Reportpath)
	if err := os.MkdirAll(genpath, 0700); err != nil {
		return "", fmt.Errorf("MakeGenDir can't mkdir %#v: %v", genpath, err)
	}

	return genpath, nil
}
