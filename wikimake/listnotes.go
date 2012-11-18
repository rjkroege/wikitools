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
  "github.com/rjkroege/wikitools/article"
  "fmt";
  "os";
  "strings";
  "time";
)

/**
 * Turns a time in ns since epoch into a string
 */
func dateToString(ti time.Time) string {
  return ti.Format(time.RFC3339);
}

func dateForPeople(ti time.Time) string {
  return ti.Format("Monday, Jan _2, 2006");
}


// TODO(rjkroege): must read the directory from the command line?
func main() {
  fmt.Printf("hello world\n");
  pwd, _ := os.Getwd();
  
  // get a directory listing
  fd, _ := os.OpenFile(".", os.O_RDONLY, 0);
  dirs, _ := fd.Readdir(-1);	// < 0 means get all of them
  
  e := make([]*article.MetaData, len(dirs));
  i := 0;
  
  // Can finess the "template" here.
  
  for _, d := range dirs {
    if strings.HasSuffix(d.Name(), ".md") {
      // TODO(rjkroege): could be a constructor like object.
      // This code could be much more designed. And less hacky.
      e[i] = new(article.MetaData)
      e[i].Name = d.Name()
      e[i].SourceForName(pwd)
      e[i].UrlForName(pwd)
      e[i].DateFromStat = d.ModTime()
      e[i].RootThroughFileForMetadata()
      i++;
    }
  }
  e = e[0:i];  // fix up the slice

  // Update the metadata objects with the desired date.
  for _, d := range e {
    // TODO(rjkroege): insert computing the Date string here.
    if (!d.DateFromMetadata.IsZero()) {
      d.FinalDate = dateToString(d.DateFromMetadata);
      d.PrettyDate = dateForPeople(d.DateFromMetadata);
    } else {
      d.FinalDate = dateToString(d.DateFromStat);
      d.PrettyDate = dateForPeople(d.DateFromStat);
    }
  }
  
  // Generate the timeline datafile.
  ofd, err := os.OpenFile("note_list.js", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644);
  if err != nil {
    fmt.Print(err);
  } else {
    // fmt.Print("attempt to print out the metadata\n")
    article.WriteTimeline(ofd, e);
  }

  // Generate articles here so that we can inline the JavaScript data if
  // that would prove desirable.
  for _, d := range e {
    // fmt.Print("running WriteHtmlFile\n")
    d.WriteHtmlFile()
  }

}