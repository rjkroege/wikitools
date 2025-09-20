package tidy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

type articleReportEntry struct {
	Path     string
	Title    string
	Date     string
	RealDate time.Time
}

type metadataReport struct {
	missingmd [][]*articleReportEntry
	tmpl      *template.Template
	settings  *wiki.Settings
}

func (abc *metadataReport) recordMetadataState(md *article.MetaData, path string) {
	abc.missingmd[md.Type()] = append(abc.missingmd[md.Type()], &articleReportEntry{
		Path:     path,
		Title:    md.Title,
		Date:     md.DetailedDate(),
		RealDate: md.PreferredDate(),
	})
}

func NewMetadataReporter(settings *wiki.Settings) (corpus.Tidying, error) {
	// TODO(rjk): These should be configurable?
	tmpl, err := template.New("newstylemetadata").Parse(iawritermetadataformat)
	if err != nil {
		return nil, fmt.Errorf("can't NewMetadataReporter %v", err)
	}
	return &metadataReport{
		missingmd: make([][]*articleReportEntry, article.MdModern+1),
		tmpl:      tmpl,
		settings:  settings,
	}, nil
}

func (abc *metadataReport) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}

	d, err := os.Stat(path)
	if err != nil {
		log.Println("metadataReport Stat error", err)
		return fmt.Errorf("can't metadataReport Stat %s: %v", path, err)
	}

	ifd, err := os.Open(path)
	if err != nil {
		log.Println("metadataReport Open error", err)
		return fmt.Errorf("can't metadataReport Open %s: %v", path, err)
	}
	defer ifd.Close()
	fd := bufio.NewReader(ifd)

	// TODO(rjk): RootThroughFileForMetadata needs to return an error when it fails
	md := article.MakeMetaData(filepath.Base(path), d.ModTime())
	md.RootThroughFileForMetadata(fd)

	abc.recordMetadataState(md, path)
	return nil
}

const cleaningreportformat = `{{template "newstylemetadata" .Metadata}}{{range .Sections}}# {{ .Name }}

{{range .Articles}}* [{{.Title}}]({{.Path}}), {{.Date}}
{{end}}
{{end}}
`

// I want sections for each type
// A list of liniks (how do I do nested templates?) Time to learn

type CompleteDocument struct {
	Metadata *IaWriterMetadataOutput
	Sections []MetadataSection
}

type MetadataSection struct {
	Name     string
	Articles []*articleReportEntry
}

func (abc *metadataReport) _genMetadataSections() []MetadataSection {
	sections := make([]MetadataSection, len(abc.missingmd))
	for i := range abc.missingmd {
		v := ByDate(abc.missingmd[i])
		sort.Sort(v)
		m := &sections[i]
		m.Name = article.Metadatanametable[i]
		m.Articles = abc.missingmd[i]
	}
	return sections
}

func (abc *metadataReport) SummaryWrite(w io.Writer) error {
	if abc.settings.OutputType == wiki.OutputHTML {
		return abc._htmlMetaReport(w, abc._genMetadataSections())
	}

	// Build up report structure
	nmd := &IaWriterMetadataOutput{
		Title: "Metadata Report",
		Date:  article.DetailedDateImpl(time.Now()),
		Tags:  "@report",
	}

	report := CompleteDocument{
		Metadata: nmd,
		Sections: abc._genMetadataSections(),
	}

	b := bufio.NewWriter(w)
	defer b.Flush()

	if _, err := abc.tmpl.New("cleaningreport").Parse(cleaningreportformat); err != nil {
		return fmt.Errorf("can't cleaningreport template%v", err)
	}

	if err := abc.tmpl.ExecuteTemplate(b, "cleaningreport", report); err != nil {
		log.Println("oops, bad template write because", err)
		return fmt.Errorf("can't writeUpdatedMetadata Execute template: %v", err)
	}
	return nil
}

func (tr *metadataReport) SummaryEncode(e *json.Encoder) error {
	return e.Encode(tr._genMetadataSections())
}

func (abc *metadataReport) _htmlMetaReport(w io.Writer, sections []MetadataSection) error {
	if _, err := abc.tmpl.New("meta_html_report").Parse(meta_html_report); err != nil {
		return fmt.Errorf("can't meta_html_report template%v", err)
	}
	return abc.tmpl.ExecuteTemplate(w, "meta_html_report", sections)
}

const meta_html_report = `
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
<h1>Metadata Report</h1>

{{range .}}<h2> {{ .Name }} </h2>
<ul>{{range .Articles}}
	<li><a href="plumb:/{{.Path}}">{{.Title}}</a>, {{.Date}}</li>
{{end}}</ul>
{{end}}

</body>
</html>
`

var _ corpus.Tidying = (*metadataReport)(nil)

type ByDate []*articleReportEntry

func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].RealDate.Before(a[j].RealDate) }

func (a ByDate) Len() int { return len(a) }
