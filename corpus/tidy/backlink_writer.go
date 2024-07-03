package tidy

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
)

type backlinkWriter struct {
	settings *wiki.Settings
	dryrun   bool
}

func NewBacklinkwriter(settings *wiki.Settings, dryrun bool) (*backlinkWriter, error) {
	return &backlinkWriter{
		settings: settings,
		dryrun:   dryrun,
	}, nil
}

func (blw *backlinkWriter) EachFile(path string, info os.FileInfo, err error) error {
	linkies := corpus.MakeLinks(search.MakeWikilinkNameIndex(blw.settings.Wikidir), blw.settings.Wikidir)
	if err := onefileimpl(blw.settings, linkies, path, info, err); err != nil {
		return fmt.Errorf("backlinkWriter.EachFile onefileimpl on %q fail %w", path, err)
	}

	// linkies now contains the desired backlinks right? (Front links too but whatever.)
	nbl := linkies.BackLinks[path]
	obl, err := article.ReadBacklinks(path)

	if err != nil && errors.Is(err, errors.New("attribute not found")) {
		return fmt.Errorf("backlinkWriter.EachFile can't ReadBacklinks on %q fail %w", path, err)
	}

	if err == nil {
		for k := range obl {
			nbl[k] = corpus.Empty{}
		}
	}

	// No point writing nothing.
	if len(nbl) < 1 {
		return nil
	}

	// This is might not be the nicest way to do things.
	log.Println(">>>", path, nbl)
	if blw.dryrun {
		return nil
	}

	if err := article.WriteBacklinks(path, nbl); err != nil {
		return fmt.Errorf("backlinkWriter.EachFile can't WriteBacklinks on %q fail %w", path, err)
	}
	return nil
}

func (fm *backlinkWriter) Summary() error {
	// I don't think that I needed to do something here.
	log.Println("backlinkWriter.Summary")
	return nil
}
