package tidy

import (
	"bufio"
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
	"go.abhg.dev/goldmark/wikilink"
)

type urlReport struct {
	settings *wiki.Settings

	// The store of links both forward and backwards.
	links *corpus.Links

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
		links:    corpus.MakeLinks(search.MakeWikilinkNameIndex(settings.Wikidir), settings.Wikidir),
		tmpl:     tmpl,
	}, nil
}

func onefileimpl(settings *wiki.Settings, links *corpus.Links, path string, info os.FileInfo, err error) error {
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

// TODO(rjk): Above, I blithered about how to refactor this to share the
// logic for writing a backing database of URLs with this code. I can
// pull the walking out and just create a different Summary
// implementation.
func (abc *urlReport) Summary() error {
	path, err := abc.settings.MakeGenDir()
	if err != nil {
		return err
	}

	if _, err := abc.tmpl.New("urlreport").Parse(urllistingreport); err != nil {
		return fmt.Errorf("can't cleaningreport template%v", err)
	}

	tpath := filepath.Join(path, "urlreport"+wiki.Extension)
	nfd, err := os.Create(tpath)
	if err != nil {
		return fmt.Errorf("can't urlReport Create %q: %v", tpath, err)
	}
	defer nfd.Close()

	// Zipper over the various outgoing links.
	articles := make(map[string][]string)
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
	for k, v := range abc.links.DamagedLinks {
		for u := range v {
			articles[k] = append(articles[k], "*damaged* "+u.Markdown())
		}
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

	if err := abc.tmpl.ExecuteTemplate(nfd, "urlreport", report); err != nil {
		log.Println("oops, bad template write because", err)
		return fmt.Errorf("can't urlReport Execute template: %v", err)
	}
	return nil
}
