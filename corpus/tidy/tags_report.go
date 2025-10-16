package tidy

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

type tagsReport struct {
	tags     map[string]int
	settings *wiki.Settings
}

func NewTagsReporter(settings *wiki.Settings) (corpus.Tidying, error) {
	return newTagsReporterImpl(settings)
}

func newTagsReporterImpl(settings *wiki.Settings) (*tagsReport, error) {
	return &tagsReport{
		tags:     make(map[string]int),
		settings: settings,
	}, nil
}

func (abc *tagsReport) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("TagsReporter couldn't read %#v: %v", path, err)
	}

	d, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("can't TagsReporter Stat %s: %#v", path, err)
	}

	ifd, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can't TagsReporter Open %s: %#v", path, err)
	}
	defer ifd.Close()
	fd := bufio.NewReader(ifd)

	// TODO(rjk): RootThroughFileForMetadata needs to return an error when it fails
	md := article.MakeMetaData(filepath.Base(path), d.ModTime())
	md.RootThroughFileForMetadata(fd)

	abc.recordTags(md)
	return nil
}

func (tr *tagsReport) recordTags(md *article.MetaData) {
	for _, t := range md.Tags {
		if _, ok := tr.tags[t]; ok {
			tr.tags[t] += 1
		} else {
			tr.tags[t] = 1
		}
	}
}

type Tagreport struct {
	Tag   string
	Count int
}

func (tr *tagsReport) prepReport() []Tagreport {
	ts := make([]Tagreport, 0)

	for k, v := range tr.tags {
		ts = append(ts, Tagreport{k, v})
	}

	slices.SortFunc(ts, func(a, b Tagreport) int {
		if n := cmp.Compare(a.Count, b.Count); n != 0 {
			return -n
		}
		// If names are equal, order by tag
		return cmp.Compare(a.Tag, b.Tag)
	})

	return ts
}

func (tr *tagsReport) SummaryWrite(w io.Writer) error {
	if tr.settings.OutputType == wiki.OutputHTML {
		return _htmlTagsReport(w, tr.prepReport())
	}

	b := bufio.NewWriter(w)
	defer b.Flush()

	for _, t := range tr.prepReport() {
		if _, err := fmt.Fprintf(b, "%s: %d\n", t.Tag, t.Count); err != nil {
			return err
		}
	}
	return nil
}

func (tr *tagsReport) SummaryEncode(e *json.Encoder) error {
	return e.Encode(tr.prepReport())
}

const tagreporttmpl = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Wiki Tag Summary</title>
  <style>
    .auto-column-list {
      column-width: 20ch;
      column-gap: 4rem;
    }
    .auto-column-list li {
      break-inside: avoid;
    }
  </style>
</head>
<body>
<h1>Wiki Tag Summary</h1>
<ul class="auto-column-list">
{{range .}}
    <li>{{.Tag}}: {{.Count}}</li>
{{end}}
</ul>
  </ul>
</body>
</html>
`

func _htmlTagsReport(w io.Writer, taglist []Tagreport) error {
	t := template.Must(template.New("list").Parse(tagreporttmpl))
	return t.Execute(w, taglist)
}

var _ corpus.Tidying = (*tagsReport)(nil)
