package corpus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"strings"

	"github.com/rjkroege/wikitools/wiki"
)

// Tidying is the interface implemented by each of the kinds of Tidying
// passes.
type Tidying interface {
	// EachFile is called by the filepath.Walk over each valid wiki file in the wiki tree.
	EachFile(path string, info os.FileInfo, err error) error

	// SummaryWrite provides the final output to the provided io.Writer.
	// TODO(rjk): I should make this more complicated. In a way that
	// permits all the file actions to happen in parallel? The parsing of all
	// the articles is definitely something that can transpire concurrently.
	SummaryWrite(w io.Writer) error

	// SummaryEncode provides the final output to the provided JSON
	// encoder.
	SummaryEncode(e *json.Encoder) error
}

type FileRecord struct {
    Path    string    `json:"path"`
	RelPath string 
    ModTime time.Time `json:"modtime"`
}

// ListAllWikiFiles is a boring implementation of Tidying that lists all files.
type listAllWikiFiles struct {
	Files []FileRecord
	settings *wiki.Settings
}

func NewListAllWikiFilesTidying(settings *wiki.Settings) (Tidying, error) {
	return &listAllWikiFiles{
		settings: settings,
	}, nil
}

func (tidy *listAllWikiFiles) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}

	tidy.Files = append(tidy.Files, FileRecord{
		Path: path,
		ModTime: info.ModTime(),
		RelPath: strings.TrimPrefix(path, tidy.settings.Wikidir),
	})
	return nil
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

func Everyfile(settings *wiki.Settings, tidying Tidying) error {
	// TODO(rjk): I have a Map/Reduce op here. I could make it parallel.

	if err := filepath.Walk(settings.Wikidir, func(path string, info os.FileInfo, err error) error {
		if settings.NotArticle(path, info) {
			return nil
		}
		return tidying.EachFile(path, info, err)
	}); err != nil {
		return fmt.Errorf("Everyfile walking: %v", err)
	}

	return nil
}

// TODO(rjk): This version would search the corpus
// Write me. Use the Spotlight tooling to extract a window.
func Filteredfiles() {
}

// TODO(rjk): It's conceivable that this API could be better?
// I had considered using the link structure but I think no.
type UrlRecorder interface {
	RecordUrl(displaytext, url, filepath string)
	RecordWikilink(displaytext, wikitext, filepath string)
}

var _ Tidying = (*listAllWikiFiles)(nil)
