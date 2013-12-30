/*
  Scans for files, extracts metadata, emits
  JSON file list.
  
  ; fn gogo  {6g -I . listnotes.go  metadata.go generatemarkup.go && 6l  listnotes.6  && ./6.out ; echo}

*/

package main

import (
  "fmt";
  "os";
  "strings";
  "time";
)

type FileMetaData struct {
  Name string;
  Url string;
  DateFromStat int64;
  DateFromMetadata int64;
  Title string;
  FinalDate string;
}

/**
 * Constructs a URL path equivalent to the given source
 * file.
 */ 
func makeUrlFromName(f string, path string) string {
  // Prefix file:///<path>/fname.html
  return "file://" + path + "/" + f[0:len(f) - len(".md")] + ".html";
}

/**
 * Turns a time in ns since epoch into a string
 */
func dateToString(ti int64) string {
  t := time.SecondsToLocalTime(int64(ti / 1e9));
  return t.Format(time.RFC3339);
}


func main() {
  // fmt.Printf("hello world\n");
  pwd, _ := os.Getwd();
  
  // get a directory listing
  fd, _ := os.Open(".", os.O_RDONLY, 0);
  dirs, _ := fd.Readdir(-1);	// < 0 means get all of them
  
  e := make([]*FileMetaData, len(dirs));
  i := 0;
  
  for _, d := range dirs {
    if strings.HasSuffix(d.Name, ".md") {
      e[i] = new(FileMetaData);
      e[i].Name = d.Name;
      e[i].Url = makeUrlFromName(e[i].Name, pwd);
      e[i].DateFromStat = d.Mtime_ns;
      e[i].DateFromMetadata, e[i].Title = rootThroughFileForMetadata(d.Name);
      
      i++;
    }
  }
  e = e[0:i];  // fix up the slice

  // Update the metadata objects with the desired date.
  for _, d := range e {
    // TODO(rjkroege): insert computing the Date string here.
    if (d.DateFromMetadata > int64(0)) {
      d.FinalDate = dateToString(d.DateFromMetadata);
    } else {
      d.FinalDate = dateToString(d.DateFromStat);
    }
  }
  
  // Time to generate the stuff...
  ofd, err := os.Open("note_list.js", os.O_WRONLY | os.O_CREATE, 0644);
  if err != nil {
    fmt.Print(err);
  } else {
    writeMarkup(ofd, e);
  }
}
