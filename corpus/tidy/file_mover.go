package tidy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
)


type FileMove struct {
	From string
	To string
}

type fileMover struct {
	moves []FileMove
	settings           *wiki.Settings
}

// NewFilemover creates a Tidying implementation that positions files in
// the right wiki directories
func NewFilemover(settings *wiki.Settings) (corpus.Tidying, error) {
	return &fileMover{
		moves: make([]FileMove, 0),
		settings:           settings,
	}, nil
}

// TODO(rjk): Need to move dependent files (i.e. images)
// fixing is not as good as I'd like
// TODO(rjk): rename source links too
// TODO(rjk): Update the index correctly
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
		// nothing to do for this file.
		return nil
	}

	destarticle := filepath.Join(fm.settings.Wikidir, destreldir, destname+destuniquing+destext)
	fm.moves = append(fm.moves, FileMove{ From: abspath, To: destarticle})

	if fm.settings.Dryrun {
		return nil
	}

	wiki.SafeMoveFile(abspath, destarticle)

	return nil
}



func (fm *fileMover) SummaryWrite(_ io.Writer) error {

	if fm.settings.Dryrun {
		return nil
	}

	dirs := make(map[string]struct{}, len(fm.moves))
	for _, v := range fm.moves {
		dirs[filepath.Dir(v.From)] = struct{}{}
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

func (tr *fileMover) SummaryEncode(_ *json.Encoder) error {
	return nil
}

var _ corpus.Tidying = (*fileMover)(nil)
