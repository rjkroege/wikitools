package main

import (
    "bytes" 
    // "fmt"			// needed for debugging.
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

// insert constants here for the various templates.


type Handler func([]string)

// TODO(rjkroege): I can refactor this with the other tool.
type Article struct {
    filename string
    filepath string
    Title string
    PrettyDate string
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
// basepath = "/Users/rjkroege/Dropbox/wiki2/"
basepath = "/Users/rjkroege/"
extension = ".md"
timeformat = "20060102-150405"
)

func Makearticle(args []string) *Article {
    s := strings.Join(args, " ");
    a := Article{ strings.Map(filter, s), "", s, time.Now().Format(time.UnixDate), nil}
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
    // Would a better time format be nicer?
    md.filepath = md.filename + "-" + time.Now().Format(timeformat) + extension
    return md.filepath
}

const (
journaltmpl = 
`title: {{.Title}}
date: {{.PrettyDate}}
tags: @journal

Yo dawg! Write stuff here.
`

booktmpl =
`title: {{.Title}}
date: {{.PrettyDate}}
tags: @bib

Yo dawg! Put the bookreview here.
`

)

// Connect up to Acme.
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

func (md *Article) Emit(tpl string) *Article {
    f :=  template.Must(template.New("footer").Parse(tpl));

    md.Buffy = new(bytes.Buffer)
    f.Execute(md.Buffy, md)
    return md
}

func journal(args []string) {
    // fmt.Print("setup a new journal article", args, "\n");
   Makearticle(args).Emit(journaltmpl).Plumb()
}

func book(args []string) {
    // fmt.Print("setup a new book review", args, "\n");
   Makearticle(args).Emit(booktmpl).Plumb()
}

// TODO(rjkroege): use @foo as a tag that goes in the tags entry (to create trails or the like)
// TODO(rjkroege): add usage output on failure.
func main() {
    handlers := map[string]Handler{
        "journal": journal,
        "book": book }

    if len(os.Args) < 2 {
        log.Fatal("Not enough arguments\n");
    }

    f, ok := handlers[os.Args[1]];
    if !ok {
        log.Fatal("Unsupported sub-command\n");
    }
     f(os.Args[2:]);
}

