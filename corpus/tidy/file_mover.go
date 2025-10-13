package tidy

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"html/template"
	"strings"

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

func newFilemoverImpl(settings *wiki.Settings) *fileMover {
	return &fileMover{
		moves: make([]FileMove, 0),
		settings:           settings,
	}
}

// NewFilemover creates a Tidying implementation that positions files in
// the right wiki directories
func NewFilemover(settings *wiki.Settings) (corpus.Tidying, error) {
	return newFilemoverImpl(settings), nil
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
	if destarticle != abspath {
		fm.moves = append(fm.moves, FileMove{ From: abspath, To: destarticle})
	}
	return nil
}

// TODO(rjk): Don't forget to update the index data here.
func (fm *fileMover) moveFiles() []string {
	dirs := make(map[string]struct{}, len(fm.moves))
	errors := make([]string,0)
	for _, v := range fm.moves {
		dirs[filepath.Dir(v.From)] = struct{}{}
		if err := wiki.SafeMoveFile(v.From, v.To); err != nil {
			errors = append(errors, fmt.Sprintf("move %q to %q failed: %v", v.From, v.To, err))
		}
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
	return errors
}


const movewikihtml = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>File Motion Report</title>
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
{{if .Dryrun}}
<h1>Article Candidates for Relocation</h1>
 <div class="list-wrapper">
    <ul class="fill-across">
{{range .Actions}}
    <li>{{ filetourl .From }} → {{ trim .To }}</li>
{{end}}
    </ul>
  </div>
{{else}}
<h1>Articles actually Relocated</h1>
 <div class="list-wrapper">
    <ul class="fill-across">
{{range .Actions}}
    <li>{{ trim .From }} → {{ filetourl .To }}</li>
{{end}}
    </ul>
  </div>
{{if .Errors}}
<h1>Errors</h1>
<ul>
{{range .Errors}}
	<li>{{.}}</li>
{{end}}
</ul>
{{end}}
{{end}}
</body>
</html>
`

func (fm *fileMover) _htmlFileMotionReport(w io.Writer, results *Results) error {
	// TODO(rjk): This needs to be cached for reuse.
	// Central state tracking needs to happen.
		tmpl, err := template.New("movewikihtml").Funcs(template.FuncMap{
				"filetourl": func(path string) template.HTML {
					return template.HTML(filetourl(fm.settings.Wikidir, path))
				},
				"trim":  func(path string) string {
					return strings.TrimPrefix(path, fm.settings.Wikidir)
				},
			}).Parse(movewikihtml)
		if  err != nil {
			return fmt.Errorf("can't movewikihtml template%v", err)
		}

	return tmpl.ExecuteTemplate(w, "movewikihtml", results)
}

const movementtmpl = `{{if .Dryrun}}will move{{else}}moved{{end}}
{{range .Actions}}
{{.From}} → {{.To}}
{{end}}
{{if .Errors}}
Errors:
{{range .Errors}}
{{.}}
{{end}}
{{end}}
`


func (fm *fileMover) SummaryWrite(w io.Writer) error {
	// TODO(rjk): Dump the content here.
	errors := []string{}
	if !fm.settings.Dryrun {
		errors = fm.moveFiles()
	}

	results := &Results{
		Dryrun: fm.settings.Dryrun,
		Actions: fm.moves,
		Errors: errors,
	}

	// TODO(rjk): need to plumb this nicely.
	if fm.settings.OutputType == wiki.OutputHTML {
		// TODO(rjk): Wire me up
		return fm._htmlFileMotionReport(w, results)
	}


	// TODO(rjk): Cache this properly.
	t := template.Must(template.New("movementtmpl").Parse(movementtmpl))
	return t.Execute(w, results)
}

type Results struct {
	Dryrun bool
	Actions  []FileMove
	Errors []string
}

func (fm *fileMover) SummaryEncode(e *json.Encoder) error {
	errors := []string{}
	if !fm.settings.Dryrun {
		errors = fm.moveFiles()
	}

	results := &Results{
		Dryrun: fm.settings.Dryrun,
		Actions: fm.moves,
		Errors: errors,
	}

	return e.Encode(results)
}

var _ corpus.Tidying = (*fileMover)(nil)
