package main

import (
    "fmt"
    "os"
    "log"
   "time"
    "strings"
    "text/template"
    "bufio"
    "os/exec"
)


/*
    Go is awesome.
*/

// insert constants here for the various templates.


type Handler func([]string)

// TODO(rjkroege): I can refactor this with the other tool.
type Article struct {
    Filename string
    Title string
    PrettyDate string
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

func Makearticle(args []string) *Article {
    s := strings.Join(args, " ");
    a := Article{ strings.Map(filter, s) + ".md", s, time.Now().Format(time.UnixDate)}
    return &a;
}

const (
path = "/Users/rjkroege/Dropbox/wiki2/"

// for debugging
// path = "./"

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

// Might want to read where Plan9 lives from the environment?
// Add wikimake to the bar for this edit pain.
// Might want to sdd support for connecting to plumber
func (md *Article) Plumb() {
    // TODO(rjkroege): make pathifying a method on md
    err := exec.Command("/usr/local/plan9/bin/plumb", path + md.Filename).Run()
    if err != nil {
        log.Fatal(err)
    }
}

func (md *Article) Emit(tpl string) *Article {
    pth := path + md.Filename
    ofd, werr := os.OpenFile(pth, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644);
    f :=  template.Must(template.New("footer").Parse(tpl));
    defer ofd.Close()
    if werr != nil {
        log.Fatal("couldn't open", pth)
    }
    w := bufio.NewWriter(ofd)
    defer w.Flush()
    f.Execute(w, md)
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
// TODO(rjkroege): add usage output
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

