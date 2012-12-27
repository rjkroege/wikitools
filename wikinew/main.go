package main

import (
    "bytes" 
    // "fmt"			// needed for debugging.
    "github.com/rjkroege/wikitools/wiki"
    "log"
    "os"
    "strings"
    "text/template"
   "code.google.com/p/goplan9/plan9/acme"
   "time"
)

/*
    Go is awesome.
*/

// TODO(rjkroege): I can refactor this with the other tool.
type Article struct {
    filename string
    filepath string
    Title string
    PrettyDate string
    Tags []string
    Buffy *bytes.Buffer
}

func filter(r rune) rune {
    lut := map[rune]rune { 
        ' ':  '-',
        '/':  ',',
        '#':  ',',
        '\t': '-'  }
    nr, ok := lut[r]
    if !ok {
        return r
    }
    return nr
}

const (
basepath = "/Users/rjkroege/Dropbox/wiki2/"
// basepath = "/Users/rjkroege/"
extension = ".md"
timeformat = "20060102-150405"
)

func Makearticle(args []string, tags []string) *Article {
    s := strings.Join(args, " ");
    a := Article{ strings.Map(filter, s), "", s, time.Now().Format(time.UnixDate), tags, nil}
    return &a
}

func (md *Article) Filepath() string {
    if md.filepath != "" {
        return md.filepath
    }

    p :=  basepath + md.filename + extension
    _,  err := os.Stat(p)
    if err != nil {
        md.filepath = p
        return p
    }
    md.filepath = md.filename + "-" + time.Now().Format(timeformat) + extension
    return md.filepath
}

func (md *Article) Plumb() {
    win, err  := acme.New();
    if err != nil {
        log.Fatal(err)
    }
    
    err =  win.Name(md.Filepath())
    if err != nil {
        log.Fatal(err)
    }

    _, err = win.Write("body", md.Buffy.Bytes())
    if err != nil {
        log.Fatal(err)
    }

    err = win.Fprintf("tag", "wikimake")
    if err != nil {
        log.Fatal(err)
    }
}

func (md* Article) Tagstring() string {
    return strings.Join(md.Tags,  " ")
}

func (md *Article) Emit(tpl string) *Article {
    f :=  template.Must(template.New("footer").Parse(tpl));

    md.Buffy = new(bytes.Buffer)
    f.Execute(md.Buffy, md)
    return md
}

// TODO(rjkroege): add usage output on failure.
func main() {
    args, tags := wiki.Split(os.Args[1:])
    tm, args, tags := wiki.Picktemplate(args, tags)

    Makearticle(args, tags).Emit(tm).Plumb()
}

