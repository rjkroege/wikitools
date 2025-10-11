package tidy

import (
	"bufio"
	"fmt"
	"io"
	"log"
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
func NewMetadataUpdater(settings *wiki.Settings) (corpus.Tidying, error) {
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
	// I should apply all the changes once for greater re-use with the report.
	dryrun := mup.mdrp.settings.Dryrun
	if dryrun {
		return mup.mdrp.EachFile(path, info, err)
	}

	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}

	updatedpth, err := mup.updateMetadata(path)
	if err != nil {
		return err
	}

	if updatedpth != "" {
		if err := wiki.SafeReplaceFile(updatedpth, path); err != nil {
			return fmt.Errorf("swapFile can't: %v", err)
		}
	}
	return nil
}

// TODO(rjk): This needs to have some unit tests yes?
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

	// TODO(rjk): RootThroughFileForMetadata needs to return an error when it fails
	md := article.MakeMetaData(filepath.Base(path), d.ModTime())
	md.RootThroughFileForMetadata(fd)

	abc.mdrp.recordMetadataState(md, path)

	if md.Type() != article.MdLegacy {
		// Nothing to do.
		return "", nil
	}

	tpath := path + "-updating"
	nfd, err := os.Create(tpath)
	if err != nil {
		log.Println("replaceLegacyMetadata Create error", err)
		return "", fmt.Errorf("can't updateMetadata Create %s: %v", tpath, err)
	}
	defer nfd.Close()

	if err := abc.writeUpdatedMetadata(path, fd, nfd, md); err != nil {
		log.Println("DoMetadataUpdate", err)
		return "", fmt.Errorf("can't updateMetadata: %v", err)
	}

	return tpath, nil
}

// There are other transformations that I'll want to implement. Refactor when
// I need to. Assumes that ofd's read point is at the end of the metadata in original file.
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

func (abc *metadataUpdater) SummaryEncode(e *json.Encoder) error { return abc.mdrp.SummaryEncode(e) }
func (abc *metadataUpdater) SummaryWrite(w io.Writer) error { return abc.mdrp.SummaryWrite(w) }

var _ corpus.Tidying = (*metadataUpdater)(nil)
