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
	"go.abhg.dev/goldmark/wikilink"
)

type urlReport struct {
	iawriterlinks []string
	settings      *wiki.Settings

	// TODO(rjk): Are these thread safe? If I want to parse all files in
	// parallel, I probably want to have a different parser per goroutine.
	// Aside: if I watch for writes to wiki articles from Edwood (like slurp)
	// then it's O(number of modified articles) and will be fast.
	markdownparser goldmark.Markdown

	// This is a forward database.
	forwardlinks map[string]map[urlrecord]struct{}
	backlinks    map[string]map[string]struct{}

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
		iawriterlinks: []string{},
		settings:      settings,
		forwardlinks:  make(map[string]map[urlrecord]struct{}),
		backlinks:     make(map[string]map[string]struct{}),
		tmpl:          tmpl,
	}, nil
}

func (abc *urlReport) EachFile(path string, info os.FileInfo, err error) error {
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
			wikiextension.NewLinkminer(abc.settings, abc, path),
			// TODO(rjk): Figure out what kind of resolver that I need to write.
			&wikilink.Extender{},
		),
	)

	if err := mkp.Convert(markdowntext, io.Discard); err != nil {
		log.Println("couldn't process and discard %q: %v", path, err)
	}

	return nil
}

// TODO(rjk): I might want to make the paths better.
const urllistingreport = `{{template "newstylemetadata" .Metadata}}{{range $index, $element :=  .Articles}}*  {{ $index }}
{{range $k, $v :=  . }}	* [{{$k.DisplayText}}]({{$k.Url}})
{{end}}
{{end}}
`

type CompleteUrlReportDocument struct {
	Metadata *IaWriterMetadataOutput
	Articles map[string]map[urlrecord]struct{}
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

	nmd := &IaWriterMetadataOutput{
		Title: "Forward URL Report",
		Date:  article.DetailedDateImpl(time.Now()),
		Tags:  "@report @urls",
	}
	report := CompleteUrlReportDocument{
		Metadata: nmd,
		Articles: abc.forwardlinks,
	}

	if err := abc.tmpl.ExecuteTemplate(nfd, "urlreport", report); err != nil {
		log.Println("oops, bad template write because", err)
		return fmt.Errorf("can't urlReport Execute template: %v", err)
	}
	return nil
}

// Show that Linkminer is a goldmark.Extender
var _ wikiextension.UrlRecorder = (*urlReport)(nil)

// A Markdown URL is [blah](http://blah.blah). Or we have wikitext
// and it's [[blah]].
// TODO(rjk): differentiate [[blah]] links in some way?
// TODO(rjk): Need to actually be able to parse the links.
type urlrecord struct {
	// The [blah] part.
	DisplayText string

	// The (http://blah.blah) part or the [[blah]].
	Url string
}

func (abc *urlReport) Record(display, url, file string) {
	log.Println("Record", display, url, file)
	// Update the forward links.
	f, ok := abc.forwardlinks[file]
	if ok {
		f[urlrecord{
			DisplayText: display,
			Url:         url,
		}] = struct{}{}
		log.Println(f)
	} else {
		f = make(map[urlrecord]struct{})
		f[urlrecord{
			DisplayText: display,
			Url:         url,
		}] = struct{}{}
		abc.forwardlinks[file] = f
	}

	// If url has already been uniqued (this is a particular challenge with the
	// wiki text links) then this will all work.
	// TODO(rjk): worry about the URL uniquing.

	// Update the reverse links.
	b, ok := abc.backlinks[url]
	if ok {
		b[file] = struct{}{}
	} else {
		b = make(map[string]struct{})
		b[file] = struct{}{}
		abc.backlinks[url] = b
	}
}
