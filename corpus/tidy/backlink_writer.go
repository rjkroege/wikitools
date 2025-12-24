package tidy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	// "github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
)

type empty = struct{}

type backlinkWriter struct {
	urlrp *urlReport
}

func NewBacklinkwriter(settings *wiki.Settings, r *http.Request) (corpus.Tidying, error) {
	urlrp, err := newUrlReporterImpl(settings)
	if err != nil {
		return nil, err
	}

	return &backlinkWriter{
		urlrp: urlrp,
	}, nil
}

func (blw *backlinkWriter) EachFile(path string, info os.FileInfo, err error) error {
	if err := onefileimpl(blw.urlrp.settings, blw.urlrp.links, path, info, err); err != nil {
		return fmt.Errorf("backlinkWriter.EachFile onefileimpl on %q fail %w", path, err)
	}
	return nil
}

// linkUpdate writes the backlinks to the effected files. Run this during
// the Summary phase to have only O(n) disk writes.
func (blw *backlinkWriter) linkUpdate() []string {
	allerrors := make([]string, 0)

	for path, nbl := range blw.urlrp.links.BackLinks {
		obl, err := article.ReadBacklinks(path)

		if err != nil && errors.Is(err, errors.New("attribute not found")) {
			allerrors = append(allerrors, fmt.Sprintf("backlinkWriter.EachFile can't ReadBacklinks on %q fail %v", path, err))
		}

		if err == nil {
			for k := range obl {
				nbl[k] = empty{}
			}
		}

		if err := article.WriteBacklinks(path, nbl); err != nil {
			allerrors = append(allerrors, fmt.Sprintf("backlinkWriter.EachFile can't WriteBacklinks on %q fail %v", path, err))
		}
	}

	return allerrors
}

func (blw *backlinkWriter) SummaryWrite(w io.Writer) error {
	dryrun := blw.urlrp.settings.Dryrun
	serrors := []string{}
	if !dryrun {
		serrors = blw.linkUpdate()
	}

	bws := BacklinkWritingStatus{
		Dryrun:                dryrun,
		FilesystemErrors:      serrors,
		FilesModified:         make(map[string][]string),
		FilesWithDamagedLinks: make(map[string][]string),
	}

	if blw.urlrp.settings.OutputType == wiki.OutputHTML {
		// Perhaps a bit of a fib. This returns the files *that will be written* which
		// is larger than the files that *need to be written*.
		blw.urlrp.links.AppendStringVectorBackLinks(func(l corpus.Wikilink) string { return l.Html() }, bws.FilesModified)
		bws.FilesWithDamagedLinks = blw.urlrp._urlReportGen(true, true)
		return blw._backlinkWriterHtmlReportWrite(w, &bws)
	}

	bws.FilesWithDamagedLinks = blw.urlrp._urlReportGen(true, false)
	return blw._backlinkWriterReportWrite(w, &bws)
}

type BacklinkWritingStatus struct {
	FilesystemErrors      []string
	FilesWithDamagedLinks map[string][]string
	FilesModified         map[string][]string
	Dryrun                bool
}

const backlinkwriterhtmlreport = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Wiki Tag Summary</title>
  <style>
    .auto-column-list {
      column-width: 30ch;
      column-gap: 4rem;
    }
    .auto-column-list li {
      break-inside: avoid;
    }
  </style>
</head>
<body>
<h1>Backlink Update Status</h1>

{{ if gt (len .FilesystemErrors) 0 }}
<h2>Filesystem Errors</h2>
<ul>
{{ range .FilesystemErrors }}
  <li>{{ . }}</li>
{{ end }}
</ul>
{{ end }}

{{ if gt (len .FilesWithDamagedLinks) 0 }}
<h2>Files with damaged links</h2>
<ul class="auto-column-list">
{{range $index, $element :=  .FilesWithDamagedLinks}}
<li>{{ filetourl $index }}
	<ul>{{range . }}
		<li>{{.}}</li>{{end}}
	</ul>
</li>
{{end}}
</ul>
{{ end }}

{{ if .Dryrun }}
<h2>Dry run, files to be modified</h2>
{{ else }}
<h2>Files modified</h2>
{{ end }}
<ul class="auto-column-list">
{{range $index, $element :=  .FilesWithDamagedLinks}}
	<li>{{ filetourl $index }}</li>
{{end}}
</ul>
</body>
</html>
`

func (blw *backlinkWriter) _backlinkWriterHtmlReportWrite(w io.Writer, bws *BacklinkWritingStatus) error {
	if _, err := blw.urlrp.tmpl.New("backlinkwriterhtmlreport").Funcs(template.FuncMap{
		"filetourl": func(path string) string {
			return filetourl(blw.urlrp.settings.Wikidir, path)
		},
	}).Parse(backlinkwriterhtmlreport); err != nil {
		return fmt.Errorf("can't prepare backlinkwriterhtmlreport template%v", err)
	}

	return blw.urlrp.tmpl.ExecuteTemplate(w, "backlinkwriterhtmlreport", bws)

}

// TODO(rjk): Expand this more.
// TODO(rjk): Remove the <a> tags
const backlinkwriterconsolereport = `
{{ if gt (len .FilesystemErrors) 0 }}
{{ range .FilesystemErrors }}
{{ end }}
{{ end }}

{{ if gt (len .FilesWithDamagedLinks) 0 }}
Files with damaged links
{{range $index, $element :=  .FilesWithDamagedLinks}}
{{ filetourl $index }}{{range . }}
	{{.}}{{end}}
{{end}}
{{ end }}
`

func (blw *backlinkWriter) _backlinkWriterReportWrite(w io.Writer, bws *BacklinkWritingStatus) error {
	if _, err := blw.urlrp.tmpl.New("backlinkwriterconsolereport").Funcs(template.FuncMap{
		"filetourl": func(path string) string {
			return filetourl(blw.urlrp.settings.Wikidir, path)
		},
	}).Parse(backlinkwriterconsolereport); err != nil {
		return fmt.Errorf("can't prepare backlinkwriterconsolereport template%v", err)
	}
	return blw.urlrp.tmpl.ExecuteTemplate(w, "backlinkwriterconsolereport", bws)
}

func (blw *backlinkWriter) SummaryEncode(e *json.Encoder) error {
	dryrun := blw.urlrp.settings.Dryrun
	serrors := []string{}
	if !dryrun {
		serrors = blw.linkUpdate()
	}

	bws := BacklinkWritingStatus{
		Dryrun:                dryrun,
		FilesystemErrors:      serrors,
		FilesModified:         make(map[string][]string),
		FilesWithDamagedLinks: blw.urlrp._urlReportGen(true, false),
	}
	return e.Encode(&bws)
}

func (blw *backlinkWriter) UpdateFiles(wm corpus.WindowManager) error {
	return nil
}

var _ corpus.Tidying = (*backlinkWriter)(nil)
