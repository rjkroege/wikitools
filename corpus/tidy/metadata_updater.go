package tidy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"encoding/json"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

type metadataUpdater struct {
	mdrp *metadataReport
}

// NewMetadataUpdater creates a new Tidying implementation to update
// metadata.
func NewMetadataUpdater(settings *wiki.Settings, r *http.Request) (corpus.Tidying, error) {
	return makeMetadataUpdaterImpl(settings)
}

func makeMetadataUpdaterImpl(settings *wiki.Settings) (*metadataUpdater, error) {
	mdrp, err := newMetadataReporterImpl(settings)
	if err != nil {
		// TODO(rjk): Maybe wrap this?
		return nil, err
	}

	return &metadataUpdater{
		mdrp: mdrp,
	}, nil
}

// TODO(rjk): The current implementation of this code likely leaves the
// database in an invalid state. Test this carefully.
// TODO(rjk): Rename, relocate and purge empty directory hierarchies and
// update the various indexes.
func (mup *metadataUpdater) EachFile(path string, info os.FileInfo, err error) error {
	return mup.mdrp.EachFile(path, info, err)
}

// Concurrency note: there is no need for the mup state to be shared
// between different http requests. However: only one update pass should
// be running at a time.
func (mup *metadataUpdater) updateAllMetadata() []string {
	allerrors := make([]string, 0)

	for _, v := range mup.mdrp.missingmd[article.MdLegacy] {
		npth, err := mup.updateMetadata(v.Path)
		if err != nil {
// TODO(rjk): Be sure to record the faulting file to make sure that this is useful.
			allerrors = append(allerrors, fmt.Sprintf("updateMetadata %q: %v", npth, err))
			continue
		}
		if err := wiki.SafeReplaceFile(npth, npth); err != nil {
			allerrors = append(allerrors, fmt.Sprintf("metadata SafeReplaceFile %q: %v", npth, err))
		}
	}
	return allerrors
}

// updateMetadata updates the format of the metadata from legacy to
// modern. It does not rename the file. It will standardize the date
// format iff the underlying date format can be parsed.
func (abc *metadataUpdater) updateMetadata(path string) (string, error) {
	d, err := os.Stat(path)
	if err != nil {
		log.Println("updateMetadata Stat error", err)
		return "", fmt.Errorf("can't DoMetadataUpdate Stat %s: %v", path, err)
	}

	ifd, err := os.Open(path)
	if err != nil {
		log.Println("updateMetadata Open error", err)
		return "", fmt.Errorf("can't DoMetadataUpdate Open %s: %v", path, err)
	}
	defer ifd.Close()
	fd := bufio.NewReader(ifd)

	// TODO(rjk): RootThroughFileForMetadata needs to return an error when it fails.
	md := article.MakeMetaData(filepath.Base(path), d.ModTime())
	md.RootThroughFileForMetadata(fd)

	tpath := path + "-updating"
	nfd, err := os.Create(tpath)
	if err != nil {
		log.Println("replaceLegacyMetadata Create error", err)
		return "", fmt.Errorf("can't updateMetadata Create %s: %v", tpath, err)
	}
	defer nfd.Close()

	// TODO(rjk): Could do buffered output.
	if err := abc.writeUpdatedMetadata(path, fd, nfd, md); err != nil {
		log.Println("DoMetadataUpdate", err)
		return "", fmt.Errorf("can't updateMetadata: %v", err)
	}
	return tpath, nil
}

// TODO(rjk): There are other transformations that I'll want to
// implement. Refactor when I need to. Assumes that ofd's read point is
// at the end of the metadata in original file.
func (abc *metadataUpdater) writeUpdatedMetadata(path string, ofd io.Reader, nfd io.Writer, md *article.MetaData) error {
	// write new metadata to nfd
	nmd := &IaWriterMetadataOutput{
		Title:     md.Title,
		Date:      md.DetailedDate(),
		Tags:      md.Tagstring(),
		Extrakeys: md.ExtraKeys(),
	}

	//	log.Printf("nmd: %#v\n", nmd)

	if err := abc.mdrp.tmpl.Execute(nfd, nmd); err != nil {
		log.Println("oops, bad template write because", err)
		return fmt.Errorf("can't writeUpdatedMetadata Execute template: %v", err)
	}

	// write existing file minus its metadata to it (first line after the first blank line)
	_, err := io.Copy(nfd, ofd)
	return err
}

type IaWriterMetadataOutput struct {
	Title     string
	Date      string
	Tags      string
	Extrakeys map[string]string
}

const iawritermetadataformat = `---
title: {{.Title}}
date: {{.Date}}{{if ne .Tags "" }}
tags: {{.Tags}}{{end}}{{range $key, $value := .Extrakeys}}
{{$key}}: {{$value}}{{end}}
---

`

func (mup *metadataUpdater) SummaryEncode(e *json.Encoder) error {
	mubu := &MetadataUpdaterOutputBundle{
		Articles: mup.mdrp.missingmd[article.MdLegacy],
	}

	dryrun := mup.mdrp.settings.Dryrun
	if !dryrun {
		 mubu.Errors = mup.updateAllMetadata()
	}
	 return e.Encode(mubu)
}

type MetadataUpdaterOutputBundle struct {
	Errors []string
	Articles []*articleReportEntry
}

const meta_updater_report = `{{range .Articles}}{{.Path}}
{{end}}
{{if .Errors}}Errors:
{{range .Errors}}
{{end}}{{end}}
`

func (mup *metadataUpdater) SummaryWrite(w io.Writer) error {
	mubu := &MetadataUpdaterOutputBundle{
		Articles: mup.mdrp.missingmd[article.MdLegacy],
	}

	dryrun := mup.mdrp.settings.Dryrun
	if !dryrun {
		 mubu.Errors = mup.updateAllMetadata()
	}

	if mup.mdrp.settings.OutputType == wiki.OutputHTML {
		// TODO(rjk): Refactor these together more nicely.
		return mup._htmlMetaUpdaterReport(w, mubu)
	}

	// Blah
	if _, err := mup.mdrp.tmpl.New("meta_updater_report").Parse(meta_updater_report); err != nil {
		return fmt.Errorf("can't meta_updater_report template%v", err)
	}
	return mup.mdrp.tmpl.ExecuteTemplate(w, "meta_updater_report", mubu)
}

var _ corpus.Tidying = (*metadataUpdater)(nil)

func (mup *metadataUpdater) _htmlMetaUpdaterReport(w io.Writer, mubu *MetadataUpdaterOutputBundle) error {
	if _, err := mup.mdrp.tmpl.New("meta_updater_html_report").Parse(meta_updater_html_report); err != nil {
		return fmt.Errorf("can't meta_updater_html_report template%v", err)
	}
	return mup.mdrp.tmpl.ExecuteTemplate(w, "meta_updater_html_report", mubu)
}

// TODO(rjk): treat an empty list more nicely.
const meta_updater_html_report = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Metadata Updater Report</title>
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
  </style>
</head>
<body>
<h1>Metadata Updater Report</h1>
	<h2>Updated Articles</h2>
	 <div class="list-wrapper">
		<ul class="fill-across">
			{{range .Articles}}
				<li><a href="plumb:/{{.Path}}">{{.Title}}</a></li>
			{{end}}
		</ul>
	</div>
{{if .Errors }}
	<h2>Errors</h2>
	<ul>
		{{range .Errors}}
			<li>{{.}}</li>
		{{end}}
	</ul>
{{end}}
</body>
</html>
`
