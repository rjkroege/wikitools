package main

import (
    "fmt"
    "os"
    "log"
   "time"
    "strings"
)


/*
    Go is awesome.
*/

// insert constants here for the various templates.


type Handler func([]string)

// TODO(rjkroege): I can refactor this with the other component.
type Article struct {
    Filename string
    Title string
    Date time.Time
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

func makearticle(args []string) *Article {
    s := strings.Join(args, " ");
    a := Article{ strings.Map(filter, s) + ".md", s, time.Now()}
    return &a;
}

func journal(args []string) {
    fmt.Print("setup a new journal article", args, "\n");
    a := makearticle(args)
    fmt.Print(a, "\n");
}

func book(args []string) {
    fmt.Print("setup a new book review", args, "\n");
    a := makearticle(args) 
    fmt.Print(a, "\n");
}

func main() {
    fmt.Print("hello world\n");

    handlers := map[string]Handler{
        "journal": journal,
        "book": book }

    // debugging...
    fmt.Print("length", len(os.Args), "\n");
    for i, _ := range(os.Args) {
        fmt.Print("> ", os.Args[i], "\n");
    }

    if len(os.Args) < 2 {
        log.Fatal("Not enough arguments\n");
    }

     f, ok := handlers[os.Args[1]];
    if !ok {
        log.Fatal("Unsupported sub-command\n");
    }

     f(os.Args[2:]);
  
}

