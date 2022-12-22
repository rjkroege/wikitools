package tidy

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

type tagsDump struct {
	tagsReport
}

func NewTagsDumper(settings *wiki.Settings) (corpus.Tidying, error) {
	return &tagsDump{
		tagsReport{
			tags:     make(map[string]int),
			settings: settings,
		},
	}, nil
}

func (tr *tagsDump) Summary() error {
	genpath, err := tr.tagsReport.settings.MakeGenDir()
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

	for k := range tr.tags {
		if _, err := fd.WriteString(k); err != nil {
			return fmt.Errorf("writing %#v failed: %v", path, err)
		}
		if _, err := fd.WriteRune('\n'); err != nil {
			return fmt.Errorf("writing %#v failed: %v", path, err)
		}
	}
	return nil
}
