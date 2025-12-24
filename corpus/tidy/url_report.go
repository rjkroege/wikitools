package tidy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	// TODO(rjk): Support parsing math.
	//	mathjax "github.com/litao91/goldmark-mathjax"
	"github.com/rjkroege/wikitools/article/wikiextension"
	"github.com/rjkroege/wikitools/corpus/links"
	"go.abhg.dev/goldmark/wikilink"
)

type urlReport struct {
	settings *wiki.Settings

	// The store of links both forward and backwards.
	links *links.Links

	tmpl *template.Template
}

func newUrlReporterImpl(settings *wiki.Settings) (*urlReport, error) {
	// TODO(rjk): Centralize the report parsing. Do it only once.
	// I am doing this wrongs?
	tmpl, err := template.New("newstylemetadata").Parse(iawritermetadataformat)
	if err != nil {
		return nil, fmt.Errorf("can't NewUrlReporter template %v", err)
	}
	return &urlReport{
		settings: settings,
		links:    links.MakeLinks(settings.Mapper, settings.Wikidir),
		tmpl:     tmpl,
	}, nil
}

// TODO(rjk): Consider how I will refactor this to make it easier to structure
// the dumping of the URLs to some kind of backing storage.
// I should be able to compose the reporting vs logging functionality into
// this in some way.
func NewUrlReporter(settings *wiki.Settings, r *http.Request) (corpus.Tidying, error) {
	return newUrlReporterImpl(settings)
}

func onefileimpl(settings *wiki.Settings, links *links.Links, path string, info os.FileInfo, err error) error {
	log.Println(path)
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}

	d, err := os.Stat(path)
	if err != nil {
		log.Println("urlReport Stat error", err)
		return fmt.Errorf("can't urlReport Stat %s: %v", path, err)
	}

	ifd, err := os.Open(path)
	if err != nil {
		log.Println("urlReport Open error", err)
		return fmt.Errorf("can't urlReport Open %s: %v", path, err)
	}
	defer ifd.Close()
	fd := bufio.NewReader(ifd)

	// skip past the metadata?
	md := article.MakeMetaData(filepath.Base(path), d.ModTime())
	md.RootThroughFileForMetadata(fd)

	// I don't need to do this. Goldmark can do this.
	markdowntext, err := io.ReadAll(fd)
	if err != nil {
		log.Println("urlReport can't read the text", err)
		return fmt.Errorf("urlReport can't read the markdown file %s: %v", path, err)
	}

	// Start recording.
	recorder := links.StartRecordingForFile(path)

	// TODO(rjk): Add an extension that can record all of the links that have been
	// seen.
	// MathJax seems to make a sad
	mkp := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			//			mathjax.MathJax,
			wikiextension.NewLinkminer(settings, recorder),
			// TODO(rjk): Figure out what kind of resolver that I need to write.
			&wikilink.Extender{},
		),
	)

	if err := mkp.Convert(markdowntext, io.Discard); err != nil {
		log.Printf("couldn't process and discard %q: %v", path, err)
	}

	// Commit the link records.
	recorder.Commit()

	return nil
}

func (abc *urlReport) EachFile(path string, info os.FileInfo, err error) error {
	return onefileimpl(abc.settings, abc.links, path, info, err)
}

// TODO(rjk): I might want to make the paths better.
const urllistingreport = `{{template "newstylemetadata" .Metadata}}{{range $index, $element :=  .Articles}}*  {{ $index }}
{{range . }}	* {{.}}
{{end}}
{{end}}
`

type CompleteUrlReportDocument struct {
	Metadata *IaWriterMetadataOutput
	Articles map[string][]string
}

// TODO(rjk): These parameters should use the Pike optional parameter pattern.
func (abc *urlReport) _urlReportGen(damagedonly bool, html bool) map[string][]string {
	articles := make(map[string][]string)

	if html {
		links.AppendStringVector(func(l corpus.Wikilink) string { return l.Html() }, abc.links.DamagedLinks, articles)
	} else {
		links.AppendStringVector(func(l corpus.Wikilink) string { return l.Markdown() }, abc.links.DamagedLinks, articles)
	}


	if damagedonly {
		return articles
	}

	if html {
		links.AppendStringVector(func(l corpus.Urllink) string { return l.Html() }, abc.links.OutUrls, articles)
		links.AppendStringVector(func(l corpus.Wikilink) string { return l.Html() }, abc.links.ForwardLinks, articles)
	} else {
		links.AppendStringVector(func(l corpus.Urllink) string { return l.Markdown() }, abc.links.OutUrls, articles)
		links.AppendStringVector(func(l corpus.Wikilink) string { return l.Markdown() }, abc.links.ForwardLinks, articles)
	}

	return articles
}

// TODO(rjk): Above, I blithered about how to refactor this to share the
// logic for writing a backing database of URLs with this code. I can
// pull the walking out and just create a different Summary
// implementation.
func (abc *urlReport) SummaryWrite(w io.Writer) error {
	if abc.settings.OutputType == wiki.OutputHTML {
		articles := abc._urlReportGen(false, true)
		return abc._htmlUrlsSummaryWrite(w, articles)
	}

	articles := abc._urlReportGen(false, false)
	if _, err := abc.tmpl.New("urlreport").Parse(urllistingreport); err != nil {
		return fmt.Errorf("can't cleaningreport template%v", err)
	}

	nmd := &IaWriterMetadataOutput{
		Title: "Forward URL Report",
		Date:  article.DetailedDateImpl(time.Now()),
		Tags:  "@report @urls",
	}
	report := CompleteUrlReportDocument{
		Metadata: nmd,
		Articles: articles,
	}

	nfd := bufio.NewWriter(w)
	defer nfd.Flush()

	if err := abc.tmpl.ExecuteTemplate(nfd, "urlreport", report); err != nil {
		log.Println("oops, bad template write because", err)
		return fmt.Errorf("can't urlReport Execute template: %v", err)
	}
	return nil
}

const urlhtmlreport = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Wiki Article List</title>
 <style>
    /* --- Container --- */
    .list-wrapper {
      width:  90%;             /* whatever width you need */
      margin: 2rem auto;    /* center the block */
      border: 1px solid #ccc;
      padding: 1rem;
    }

    /* --- List reset --- */
    ul.fill-across {
      list-style: none;
      margin: 0;
      padding: 0;
      display: flex;        /* put items in a row */
      flex-wrap: wrap;        /* allow wrapping to next line */
      gap: 1rem;              /* space between items */
    }

    /* --- List items --- */
    ul.fill-across li {
      flex: 1 1 200px;        /* grow, shrink, base width 200px */
     padding: 1rem;
    }

li.inner {
	padding: 0 !important;
}
  </style>
</head>
<body>
<h1>URL Report</h1>
 <div class="list-wrapper">
	<ul class="fill-across">
	{{range $index, $element :=  .}}
	<li>{{ filetourl $index }}
	<ul>{{range . }}
		<li class="inner">{{.}}</li>{{end}}
	</ul></li>
	{{end}}
	</ul>
 </div>
</body>
</html>
`

func filetourl(prefix, path string) string {
	short := strings.TrimPrefix(path, prefix)
	return fmt.Sprintf("<a href=\"plumb:/%s\">%s</a>", path, short)
}

func (abc *urlReport) _htmlUrlsSummaryWrite(w io.Writer, articles map[string][]string) error {
	if _, err := abc.tmpl.New("urlhtmlreport").Funcs(template.FuncMap{
		"filetourl": func(path string) string {
			return filetourl(abc.settings.Wikidir, path)
		},
	}).Parse(urlhtmlreport); err != nil {
		return fmt.Errorf("can't urlhtmlreport template%v", err)
	}

	return abc.tmpl.ExecuteTemplate(w, "urlhtmlreport", articles)
}

func (abc *urlReport) SummaryEncode(e *json.Encoder) error {
	return e.Encode(abc._urlReportGen(false, false))
}

func (abc *urlReport) UpdateFiles(wm corpus.WindowManager) error {
	return nil
}

var _ corpus.Tidying = (*urlReport)(nil)
