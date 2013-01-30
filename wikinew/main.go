package main

import (
    "bytes" 
    "github.com/rjkroege/wikitools/article"
    "github.com/rjkroege/wikitools/wiki"
    "log"
    "os"
    "strings"
    "text/template"
    // "fmt"			// needeebugging.
   "code.google.com/p/goplan9/plan9/acme"
   "time"
)

const (
basepath = "/Users/rjkroege/Dropbox/wiki2/"
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

func Makearticle(args []string, tags []string) *article.MetaData {
    s := strings.Join(args, " ")
    filename :=  wiki.UniqueValidName(basepath, wiki.ValidBaseName(args), extension, SystemImpl(0))
    return article.NewArticle(filename, s, tags)
}

type ExpandedArticle struct {
   *article.MetaData
    buffy *bytes.Buffer
}

func (md *ExpandedArticle) Plumb() {
    win, err  := acme.New();
    if err != nil {
        log.Fatal(err)
    }
    
    err =  win.Name(md.Name)
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

/*
func (md* Article) Tagstring() string {
    return strings.Join(md.Tags,  " ")
}
*/

func Expand(a *article.MetaData, tpl string) *ExpandedArticle {
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

