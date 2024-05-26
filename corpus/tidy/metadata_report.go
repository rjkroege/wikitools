package tidy

import (
	"bufio"
	"fmt"
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

func (abc *metadataReport) Summary() error {
	path, err := abc.settings.MakeGenDir()
	if err != nil {
		return err
	}

	if _, err := abc.tmpl.New("cleaningreport").Parse(cleaningreportformat); err != nil {
		return fmt.Errorf("can't cleaningreport template%v", err)
	}

	// Sort the arrays by Date.

	tpath := filepath.Join(path, "metadatareport"+wiki.Extension)
	nfd, err := os.Create(tpath)
	if err != nil {
		return fmt.Errorf("can't writeMetadataUpdateReport Create %s: %v", tpath, err)
	}
	defer nfd.Close()

	sections := make([]MetadataSection, len(abc.missingmd))
	for i := range abc.missingmd {
		v := ByDate(abc.missingmd[i])
		sort.Sort(v)
		m := &sections[i]
		m.Name = article.Metadatanametable[i]
		m.Articles = abc.missingmd[i]
	}

	// Build up giant structure here...
	nmd := &IaWriterMetadataOutput{
		Title: "Metadata Report",
		Date:  article.DetailedDateImpl(time.Now()),
		Tags:  "@report",
	}

	report := CompleteDocument{
		Metadata: nmd,
		Sections: sections,
	}

	if err := abc.tmpl.ExecuteTemplate(nfd, "cleaningreport", report); err != nil {
		log.Println("oops, bad template write because", err)
		return fmt.Errorf("can't writeUpdatedMetadata Execute template: %v", err)
	}
	return nil
}

type ByDate []*articleReportEntry

func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].RealDate.Before(a[j].RealDate) }

func (a ByDate) Len() int { return len(a) }
