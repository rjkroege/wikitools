package tidy

import (
	"fmt"
	"os"

	"github.com/rjkroege/wikitools/wiki"
)

type backlinkWriter struct {
}

// TODO(rjk): Share code with the report generator.


func NewBacklinkwriter(settings *wiki.Settings, dryrun bool) (*backlinkWriter, error) {
	return &backlinkWriter{
	}, nil
}

func (fm *backlinkWriter) EachFile(path string, info os.FileInfo, err error) error {
	return fmt.Errorf("not implemented")
}

func (fm *backlinkWriter) Summary() error {
	return fmt.Errorf("not implemented")
}
