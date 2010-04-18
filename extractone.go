/*
  Rips out the meta data from a single Journal file.
  {make &&  ./extractmeta *.md}
*/

package main

import (
  "fmt";
  "os";
  "strings";
  "time";
  "flag";
)

// TODO(rjkroege): refactor into something different
// TODO(rjkroege): write some tests.
type FileMetaData struct {
  Name string;
  Url string;
  DateFromStat uint64;
  DateFromMetadata uint64;
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
func dateToString(ti uint64) string {
  t := time.SecondsToLocalTime(int64(ti / 1e9));
  return t.Format(time.ISO8601);
}


func main() {
  flag.Parse();
  pwd, _ := os.Getwd();
  
  for i := 0; i < flag.NArg(); i++ {
    fname := flag.Arg(i);
  
    // Skip files of the wrong metadata
    fi, err := os.Stat(fname);
    if strings.HasSuffix(fname, ".md") && err == nil {
      e := new(FileMetaData);
      e.Name = fname;
      e.Url = makeUrlFromName(e.Name, pwd);
      
      e.DateFromStat = fi.Mtime_ns;
      e.DateFromMetadata, e.Title = rootThroughFileForMetadata(fname);

      if (e.DateFromMetadata > uint64(0)) {
        e.FinalDate = dateToString(e.DateFromMetadata);
      } else {
        e.FinalDate = dateToString(e.DateFromStat);
      }

      fmt.Println(e.Name);
      fmt.Println(e.Url);
      fmt.Println(e.Title);
      fmt.Println(e.FinalDate);

    } else if err != nil {
      fmt.Print(err);
    } else {
      fmt.Print("given file: " + fname + " is not a .md file\n");
    }
  }
}