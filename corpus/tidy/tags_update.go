package tidy

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"io"
	"encoding/json"
	"net/http"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

type tagsDump struct {
	tagrp *tagsReport
}

func NewTagsDumper(settings *wiki.Settings, r *http.Request) (corpus.Tidying, error) {
	tagrp , err := newTagsReporterImpl(settings)
	if err != nil {
		return nil, err
	}
	return &tagsDump{
		tagrp: tagrp,
	}, nil
}

// TODO(rjk): This functionality needs to done correctly.
func (tr *tagsDump) writeTagList() error {
	genpath, err := tr.tagrp.settings.MakeGenDir()
	if err != nil {
		return err
	}

	path := filepath.Join(genpath, "taglist")
	nfd, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("can't Create %#v: %v", path, err)
	}
	defer nfd.Close()

	fd := bufio.NewWriter(nfd)
	defer fd.Flush()

	for k := range tr.tagrp.tags {
		if _, err := fd.WriteString(k); err != nil {
			return fmt.Errorf("writing %#v failed: %v", path, err)
		}
		if _, err := fd.WriteRune('\n'); err != nil {
			return fmt.Errorf("writing %#v failed: %v", path, err)
		}
	}
	return nil
}

func (tagu *tagsDump) EachFile(path string, info os.FileInfo, err error) error {
	return tagu.tagrp.EachFile(path , info , err )
}
func (tagu *tagsDump) SummaryWrite(w io.Writer) error {
	dryrun := tagu.tagrp.settings.Dryrun
	if !dryrun {
		if err := tagu.writeTagList(); err != nil {
			return err
		}
	}
	// TODO(rjk): Handle errors more nicely.
	return tagu.tagrp.SummaryWrite(w)
}
func (tagu *tagsDump) SummaryEncode(e *json.Encoder) error {
	dryrun := tagu.tagrp.settings.Dryrun
	if !dryrun {
		if err := tagu.writeTagList(); err != nil {
			return err
		}
	}
	// TODO(rjk): Handle errors more nicely.
	return tagu.tagrp.SummaryEncode(e)
}

var _ corpus.Tidying = (*tagsDump)(nil)
