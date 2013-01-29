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
    Title string
    PrettyDate string
    Tags []string
//    Buffy *bytes.Buffer
}


const (
basepath = "/Users/rjkroege/Dropbox/wiki2/"
// basepath = "/Users/rjkroege/"
extension = ".md"
timeformat = "20060102-150405"
)

type SystemImpl int;

func (s SystemImpl) Exists(path string) bool {
    _,  err := os.Stat(path)
    if err == nil {
        return true
    }
    return false
}

func (s SystemImpl) Now() time.Time {
    return time.Now()
}

func Makearticle(args []string, tags []string) *Article {
    s := strings.Join(args, " ")    
    a := Article{ wiki.UniqueValidName(basepath, wiki.ValidBaseName(args), extension, SystemImpl(0)), 
            s, time.Now().Format(time.UnixDate), tags}
    return &a
}

type ExpandedArticle struct {
   *Article
    buffy *bytes.Buffer
}

func (md *ExpandedArticle) Plumb() {
    win, err  := acme.New();
    if err != nil {
        log.Fatal(err)
    }
    
    err =  win.Name(md.filename)
    if err != nil {
        log.Fatal(err)
    }

    _, err = win.Write("body", md.buffy.Bytes())
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

func Expand(a *Article, tpl string ) *ExpandedArticle {
    f :=  template.Must(template.New("footer").Parse(tpl));

    b := new(bytes.Buffer)
    f.Execute(b, a)
    return &ExpandedArticle{a, b}
}

// TODO(rjkroege): add usage output on failure.
func main() {
    args, tags := wiki.Split(os.Args[1:])
    tm, args, tags := wiki.Picktemplate(args, tags)

    Expand(Makearticle(args, tags), tm).Plumb()
}

