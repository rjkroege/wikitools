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

type fileMover struct {
	removeddirectories map[string]struct{}
	dryrun             bool
	settings           *wiki.Settings
}

// NewFilemover creates a Tidying implementation that positions files in
// the right wiki directories
func NewFilemover(settings *wiki.Settings, dryrun bool) (corpus.Tidying, error) {
	return &fileMover{
		removeddirectories: make(map[string]struct{}),
		dryrun:             dryrun,
		settings:           settings,
	}, nil
}

// TODO(rjk): Need to move dependent files (i.e. images)
// fixing is not as good as I'd like
// TODO(rjk): rename source links too
func (fm *fileMover) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}

	// Get a filedescriptor to read the metadata.
	ifd, err := os.Open(path)
	if err != nil {
		log.Println("updateMetadata Open error", err)
		return fmt.Errorf("can't FileMover Open %s: %v", path, err)
	}
	defer ifd.Close()
	fd := bufio.NewReader(ifd)

	md := article.MakeMetaData(info.Name(), info.ModTime())
	// TODO(rjk): Move the filedescriptor logic etc. to the reading of the metadata?
	md.RootThroughFileForMetadata(fd)

	abspath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("can't find abs for %s: %v", path, err)
	}

	srcname, srcuniquing, srcext := wiki.SplitActualName(info.Name())
	srcreldir := fm.settings.SplitActualDir(abspath)

	destname := md.PreferredFileName(fm.settings)
	destreldir := md.RelativeDateDirectory()
	destuniquing := fm.settings.UniquingExtension(destreldir, destname)
	destext := fm.settings.Extension()

	// TODO(rjk): Handle dependent files. They also need some logic to have right extensions.
	mustrename := func() bool {
		if srcname != destname ||
			srcreldir != destreldir ||
			srcext != destext {
			return true
		}

		if destuniquing == "" && srcuniquing != "" {
			return true
		}

		if destuniquing != "" && srcuniquing != "" || destuniquing != "" && srcuniquing == "" {
			return false
		}

		// I think that this covers all of the cases?
		return false
	}

	if !mustrename() {
		// nothing to do
		return nil
	}

	destarticle := filepath.Join(fm.settings.Wikidir, destreldir, destname+destuniquing+destext)

	if fm.dryrun {
		log.Printf("mv %s -> %s\n", abspath, destarticle)
		fm.removeddirectories[filepath.Dir(abspath)] = struct{}{}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(destarticle), 0700); err != nil {
		return fmt.Errorf("can't mkdir %s because: %v", filepath.Dir(destarticle), err)
	}

	if err := os.Link(abspath, destarticle); err != nil {
		return fmt.Errorf("can't link %s to %s because %v", abspath, destarticle, err)
	}

	if err := os.Remove(abspath); err != nil {
		return fmt.Errorf("can't remove %s because %v", abspath, err)
	}

	// Walk does a pre-order traversal. So we might have removed all of the
	// files in a given directory but we don't know if this directory is
	// empty. But we can record the fact that we've removed something from
	// the directory and clean the directories in the Summary
	fm.removeddirectories[filepath.Dir(abspath)] = struct{}{}

	return nil
}

func (fm *fileMover) Summary() error {
	dirs := fm.removeddirectories

	if fm.dryrun {
		log.Println("not removing non-empty directories in dryrun mode")
		return nil
	}

	for workremains := true; workremains; {
		parentdirs := make(map[string]struct{})
		for d := range dirs {
			workremains = false
			if err := os.Remove(d); err == nil {
				if pd := filepath.Dir(d); pd != fm.settings.Wikidir {
					parentdirs[pd] = struct{}{}
					workremains = true
				}
			}
		}
		dirs = parentdirs
	}

	return nil
}
