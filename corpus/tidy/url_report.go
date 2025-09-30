package tidy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/corpus/links"
	"go.abhg.dev/goldmark/wikilink"
)

type urlReport struct {
	settings *wiki.Settings

	// The store of links both forward and backwards.
	links *links.Links

	tmpl *template.Template
}

// TODO(rjk): Consider how I will refactor this to make it easier to structure
// the dumping of the URLs to some kind of backing storage.
// I should be able to compose the reporting vs logging functionality into
// this in some way.
func NewUrlReporter(settings *wiki.Settings) (corpus.Tidying, error) {
	// TODO(rjk): The metadata
	tmpl, err := template.New("newstylemetadata").Parse(iawritermetadataformat)
	if err != nil {
		return nil, fmt.Errorf("can't NewUrlReporter template %v", err)
	}
	return &urlReport{
		settings: settings,
		links:    links.MakeLinks(search.MakeWikilinkNameIndex(settings.Wikidir), settings.Wikidir),
		tmpl:     tmpl,
	}, nil
}

func onefileimpl(settings *wiki.Settings, links *links.Links, path string, info os.FileInfo, err error) error {
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

	// TODO(rjk): Add an extension that can record all of the links that have been
	// seen.
	// MathJax seems to make a sad
	mkp := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			//			mathjax.MathJax,
			wikiextension.NewLinkminer(settings, links, path),
			// TODO(rjk): Figure out what kind of resolver that I need to write.
			&wikilink.Extender{},
		),
	)

	if err := mkp.Convert(markdowntext, io.Discard); err != nil {
		log.Printf("couldn't process and discard %q: %v", path, err)
	}

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

func (abc *urlReport) _urlReportGen() map[string][]string {
	// Zipper over the various outgoing links.
	articles := make(map[string][]string)
	for k, v := range abc.links.DamagedLinks {
		for u := range v {
			articles[k] = append(articles[k], "*damaged* "+u.Markdown())
		}
	}
	for k, v := range abc.links.OutUrls {
		for u := range v {
			articles[k] = append(articles[k], u.Markdown())
		}
	}
	for k, v := range abc.links.ForwardLinks {
		for u := range v {
			articles[k] = append(articles[k], u.Markdown())
		}
	}
	return articles
}

// TODO(rjk): Above, I blithered about how to refactor this to share the
// logic for writing a backing database of URLs with this code. I can
// pull the walking out and just create a different Summary
// implementation.
func (abc *urlReport) SummaryWrite(w io.Writer) error {
	articles := abc._urlReportGen()

	if abc.settings.OutputType == wiki.OutputHTML {
		return abc._htmlUrlsSummaryWrite(w, articles)
	}

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
<h1>URL Report</h1>
<ul class="auto-column-list">
{{range $index, $element :=  .}}<li>{{ $index }}
	<ul>{{range . }}<li>{{.}}</li>{{end}}
</ul></li>
{{end}}
</ul>

</body>
</html>
`

func (abc *urlReport) _htmlUrlsSummaryWrite(w io.Writer, articles map[string][]string) error {
	if _, err := abc.tmpl.New("urlhtmlreport").Parse(urlhtmlreport); err != nil {
		return fmt.Errorf("can't urlhtmlreport template%v", err)
	}

	return abc.tmpl.ExecuteTemplate(w, "urlhtmlreport", articles)
}

func (abc *urlReport) SummaryEncode(e *json.Encoder) error {
	return e.Encode(abc._urlReportGen())
}

var _ corpus.Tidying = (*urlReport)(nil)
