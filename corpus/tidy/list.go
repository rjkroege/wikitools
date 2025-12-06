package tidy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"html/template"
	"net/http"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/rjkroege/wikitools/corpus"
)

type FileRecord struct {
	Path    string `json:"path"`
	RelPath string
	ModTime time.Time `json:"modtime"`
}

// ListAllWikiFiles is a boring implementation of Tidying that lists all files.
type listAllWikiFiles struct {
	Files    []FileRecord
	settings *wiki.Settings
	tmpl *template.Template
	tags []string
}

// extractPathValueTags pulls the named wildcard value from the request,
// splits it on commas, and returns the resulting tokens.
// If the value is missing or empty the function returns an empty slice.
func extractPathValueTags(r *http.Request, tagName string) []string {
	raw := r.PathValue(tagName)
	if raw == "" {
		return []string{}
	}
	return strings.Split(raw, ",")
}

func NewListAllWikiFilesTidying(settings *wiki.Settings, r *http.Request) (corpus.Tidying, error) {
	tags := extractPathValueTags(r, "tags")
	return &listAllWikiFiles{
		settings: settings,
		tags: tags,
	}, nil
}

func (tidy *listAllWikiFiles) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}

	// If tags filter is specified, parse metadata and check tags
	if len(tidy.tags) > 0 {
		ifd, err := os.Open(path)
		if err != nil {
			log.Printf("couldn't open %s for metadata parsing: %v", path, err)
			return nil // Skip files that can't be opened
		}
		defer ifd.Close()

		md := article.MakeMetaData(filepath.Base(path), info.ModTime())
		md.RootThroughFileForMetadata(ifd)

		// Check if article has all requested tags (AND logic)
		if !hasAllTags(md.Tags, tidy.tags) {
			return nil // Skip files that don't match tag filter
		}
	}

	tidy.Files = append(tidy.Files, FileRecord{
		Path:    path,
		ModTime: info.ModTime(),
		RelPath: strings.TrimPrefix(path, tidy.settings.Wikidir),
	})
	return nil
}

// hasAllTags checks if articleTags contains all of the requiredTags
func hasAllTags(articleTags, requiredTags []string) bool {
	tagSet := make(map[string]bool)
	for _, tag := range articleTags {
		tagSet[tag] = true
	}

	for _, required := range requiredTags {
		if !tagSet[required] {
			return false
		}
	}
	return true
}

func (tidy *listAllWikiFiles) SummaryWrite(w io.Writer) error {
	log.Println("listAllWikiFiles", "SummaryWrite")
	if tidy.settings.OutputType == wiki.OutputHTML {
		return tidy._htmlSummaryWrite(w)
	}
	b := bufio.NewWriter(w)
	defer b.Flush()
	for _, s := range tidy.Files {
		if _, err := fmt.Fprintf(b, "%s: %s\n", s.Path, s.ModTime.Format(time.RFC822)); err != nil {
			return err
		}
	}
	return nil
}

func (tidy *listAllWikiFiles) SummaryEncode(e *json.Encoder) error {
	log.Println("listAllWikiFiles", "SummaryEncode", tidy.Files)
	return e.Encode(tidy.Files)
}

func (tidy *listAllWikiFiles) UpdateFiles(wm corpus.WindowManager) error {
	return nil
}

var _ corpus.Tidying = (*listAllWikiFiles)(nil)
