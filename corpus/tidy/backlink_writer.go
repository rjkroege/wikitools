package tidy

import (
	"errors"
	"fmt"
	"os"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/corpus/search"
	"github.com/rjkroege/wikitools/wiki"
)

type backlinkWriter struct {
	settings *wiki.Settings
	dryrun   bool

	// The store of links both forward and backwards.
	linkies *corpus.Links
}

func NewBacklinkwriter(settings *wiki.Settings, dryrun bool) (*backlinkWriter, error) {
	return &backlinkWriter{
		settings: settings,
		dryrun:   dryrun,
		linkies:    corpus.MakeLinks(search.MakeWikilinkNameIndex(settings.Wikidir), settings.Wikidir),
	}, nil
}

func (blw *backlinkWriter) EachFile(path string, info os.FileInfo, err error) error {
	if err := onefileimpl(blw.settings, blw.linkies, path, info, err); err != nil {
		return fmt.Errorf("backlinkWriter.EachFile onefileimpl on %q fail %w", path, err)
	}
	return nil
}

// Summary writes the backlinks to the effected files. Run this at the end
// so that I don't have an O(n^2) fan out of disk writes for each merge of the
// backlink attributes.
// TODO(rjk): I will need to pull out the merging process into a helper function
// TODO(rjk): The merging process needs to happen for all the paths that have
// backlink additions.
func (blw *backlinkWriter) Summary() error {
	allerrors := make([]error, 0)
	for path, nbl := range blw.linkies.BackLinks {
		obl, err := article.ReadBacklinks(path)

		if err != nil && errors.Is(err, errors.New("attribute not found")) {
			allerrors = append(allerrors, fmt.Errorf("backlinkWriter.EachFile can't ReadBacklinks on %q fail %w", path, err))
		}

		if err == nil {
			for k := range obl {
				nbl[k] = corpus.Empty{}
			}
		}

		// No point writing nothing.
		if len(nbl) < 1 {
			continue
		}

		if blw.dryrun {
			continue
		}

		if err := article.WriteBacklinks(path, nbl); err != nil {
			allerrors = append(allerrors, fmt.Errorf("backlinkWriter.EachFile can't WriteBacklinks on %q fail %w", path, err))
		}
	}
	return errors.Join(allerrors...)
}
