package tidy

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)

type tagsReport struct {
	tags     map[string]int
	settings *wiki.Settings
}

func NewTagsReporter(settings *wiki.Settings) (corpus.Tidying, error) {
	return &tagsReport{
		tags:     make(map[string]int),
		settings: settings,
	}, nil
	// TODO(rjk): Goo!
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

func (tr *tagsReport) Summary() error {
	// TODO(rjk): Could sort.
	for k, v := range tr.tags {
		log.Println(k, v)
	}
	return nil
}
