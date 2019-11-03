package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rjkroege/wikitools/wiki"
)

// TODO(rjk): Do I need this program?
// TODO(rjk): Correct the pathing assumptions.

var outputdir = flag.String("odir", "./converted", "Output directory for importable files")
var wikidir = flag.String("wdir", "/Users/rjkroege/gda/wiki2", "Where the wiki files are for naming")

// stripextension removes an extension from a filename if it's present and returns it.
func stripextension(fn string) string {
	extension := filepath.Ext(fn)
	return strings.TrimSuffix(fn, extension)
}

// makesafename makes a name from the provided filename argument that would be
// unique inside of the wiki directory but is actually a valid filename for the outputdir.
func makesafename(fn string) string {
	noextensionname := stripextension(fn)
	pathinwiki := wiki.UniqueValidName(*wikidir,
		wiki.ValidBaseName([]string{noextensionname}), ".md", wiki.SystemImpl(0))

	return filepath.Join(*outputdir, filepath.Base(pathinwiki))
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

func mkdir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

const almostunixlike = "Monday _2 Jan 2006 15:04:05 MST"

func main() {
	flag.Parse()

	log.Println("hello")

	filestoprocess := flag.Args()

	if err := mkdir(*outputdir); err != nil {
		log.Printf("Can't make output dir %s because: %v\n", *outputdir, err)
		os.Exit(-1)
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
		outputpath := makesafename(fsi.Name())

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
