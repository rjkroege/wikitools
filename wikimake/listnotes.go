/*
  Scans for files, extracts metadata, emits
  JSON file list.
  
  ; fn gogo  {6g -I . listnotes.go  metadata.go generatemarkup.go && 6l  listnotes.6  && ./6.out ; echo}

  Observe: entrylist enumerates over the .md files and opens each one.
  It could also construct the actual HTML files while doing so. Each 
  html file is constructed through (what is currently) a python script.

*/


package main

import (
    "fmt"
    "github.com/rjkroege/wikitools/article"
    "github.com/rjkroege/wikitools/generate"
    "io"
    "os"
    "strings"
    "time"
)

type SystemImpl int;

func (s SystemImpl) OpenFileForReading(name string) (rd io.ReadCloser, err error) {
    rd, err = os.OpenFile(name, os.O_RDONLY, 0)
    return
}

func (s SystemImpl) ModTime(name string) (modtime time.Time, err error) {
    statinfo, err := os.Stat(name)
    if err == nil {
         modtime = statinfo.ModTime()
    }
    return
}

func (s SystemImpl) OpenFileForWriting(name string) (wr io.WriteCloser, err error) {
    wr, err = os.OpenFile(name, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
    return
}

// TODO(rjkroege): must read the directory from the command line?
func main() {
    pwd, _ := os.Getwd();
    article.SetPathForContent(pwd)
  
  // get a directory listing
  fd, _ := os.OpenFile(".", os.O_RDONLY, 0);
  dirs, _ := fd.Readdir(-1);	// < 0 means get all of them
  
  e := make([]*article.MetaData, len(dirs));
  i := 0;
  
  // Can finess the "template" here.
  
  for _, d := range dirs {
    if strings.HasSuffix(d.Name(), ".md") {
        e[i] = article.NewMetaData(d.Name(), d.ModTime())
        fd, _ :=  os.OpenFile(e[i].Name, os.O_RDONLY, 0)
        e[i].RootThroughFileForMetadata(io.Reader(fd))
        fd.Close()
        i++;
    }
  }
  e = e[0:i];  // fix up the slice

  // Generate the timeline datafile.
  ofd, err := os.OpenFile("note_list.js", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644);
  if err != nil {
    fmt.Print(err);
  } else {
    // fmt.Print("attempt to print out the metadata\n")
    generate.WriteTimeline(ofd, e);
  }

  // Generate articles here so that we can inline the JavaScript data if
  // that would prove desirable.
  for _, d := range e {
    // fmt.Print("running WriteHtmlFile\n")
    generate.WriteHtmlFile(SystemImpl(0), d)
  }

}
