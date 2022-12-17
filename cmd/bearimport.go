package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rjkroege/wikitools/wiki"
)

// stripextension removes an extension from a filename if it's present and returns it.
func stripextension(fn string) string {
	extension := filepath.Ext(fn)
	return strings.TrimSuffix(fn, extension)
}

// makesafename makes a name from the provided filename argument that would be
// unique inside of the wiki directory but is actually a valid filename for the outputdir.
func makesafename(outputdir, wikidir, fn string) string {
	noextensionname := stripextension(fn)
	pathinwiki := wiki.UniqueAbsolutePath(wikidir,
		wiki.ValidBaseName([]string{noextensionname}), ".md", wiki.SystemImpl(0))

	return filepath.Join(outputdir, pathinwiki)
}

func getalltags(contents []byte) []string {
	matcher := regexp.MustCompile("#[0-9a-zA-Z]+")
	bmatches := matcher.FindAll(contents, -1)

	smatches := make([]string, 0, len(bmatches))
	for _, b := range bmatches {
		smatches = append(smatches, "@"+string(b[1:]))
	}
	return smatches
}

const almostunixlike = "Monday _2 Jan 2006 15:04:05 MST"

func Bearimport(settings *wiki.Settings, outputdir string, filestoprocess []string) {
	if err := os.MkdirAll(outputdir, 0755); err != nil {
		log.Panicf("Can't make output dir %s because: %v\n", outputdir, err)
	}

	for _, fn := range filestoprocess {
		fsi, err := os.Stat(fn)
		if err != nil {
			log.Printf("Can't stat %s because: %v\n", fn, err)
			continue
		}

		fd, err := os.Open(fn)
		if err != nil {
			log.Printf("Can't open %s because: %v\n", fn, err)
			continue
		}

		filecontents, err := ioutil.ReadAll(fd)
		if err != nil {
			log.Printf("Can't read %s because: %v\n", fn, err)
			fd.Close()
			continue
		}
		fd.Close()

		tags := getalltags(filecontents)
		origdate := fsi.ModTime()
		outputpath := makesafename(outputdir, settings.Wikidir, fsi.Name())

		// What we've done so far.
		log.Println("title: ", fsi.Name(), "date:", origdate, "filename", outputpath, "tags:", tags)

		// need to actually write the stuffs

		fd, err = os.Create(outputpath)
		if err != nil {
			log.Printf("can't open output file %s because: %v\n", outputpath, err)
			continue
		}

		if _, err := fmt.Fprintf(fd, "title: %s\ndate: %s\ntags: %s\n\n",
			fsi.Name(), origdate.Format(almostunixlike), strings.Join(tags, " ")); err != nil {
			log.Printf("can't wirte file header %s because: %v\n", outputpath, err)
			fd.Close()
			continue
		}

		if _, err := fd.Write(filecontents); err != nil {
			log.Printf("can't save file %s because: %v\n", outputpath, err)
		}
		fd.Close()
	}
}
