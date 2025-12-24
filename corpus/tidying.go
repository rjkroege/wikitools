package corpus

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/wiki"
)

// WindowManager is an interface for managing acme windows.
type WindowManager interface {
	Add(win AcmeWindow)
	Remove(id int)
	Snapshot() []AcmeWindow
}

// AcmeWindow represents a window in acme.
type AcmeWindow struct {
	ID   int
	Name string
}

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

	// UpdateFiles updates acme windows based on the WindowManager.
	UpdateFiles(wm WindowManager) error
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
// TODO(rjk): work tracked in [[IncrementalLinkDatabase]]
type LinkRecording interface {
	RecordUrl(displaytext, url string)
	RecordWikilink(displaytext, wikitext string)
	Commit()
}

type LinksRecorder interface {
	StartRecordingForFile(filepath string) LinkRecording
	AppendStringVectorForwardLinks(f func(Wikilink) string, articles map[string][]string)
	AppendStringVectorBackLinks(f func(Wikilink) string, articles map[string][]string)
	AppendStringVectorOutUrls(f func(Urllink) string, articles map[string][]string)
	AppendStringVectorDamagedLinks(f func(Wikilink) string, articles map[string][]string)
}
