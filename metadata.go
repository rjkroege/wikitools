/*
  Metadata extraction  
  ; fn gogo  {6g -I . listnotes.go  metadata.go && 6l  listnotes.6  && ./6.out ; echo}

*/

package main

import (
  "fmt";
  "os";
  "bufio";
  "io";
  "strings";
  "regexp";
  "time";
)

var metadataMatcher = regexp.MustCompile("^([A-Za-z]*):[ \t]*(.*)$");
var commentDataMatcher = regexp.MustCompile("<!-- *([0-9]*) *-->");

/**
 * Attempts to parse the date sequence that I have used in multiple
 * files that consists of a string of digits.
 */
func parseDateCmdFmt(numericDate string) uint64 {
  fmt.Println(numericDate);
  
  // Moral: I don't understand how date and time work.
  t, e1 := time.Parse("20060102150405", numericDate);
  
  retval := uint64(0);
  
  if e1 != nil {
    fmt.Print(e1);
  } else {
    fmt.Println("y:", t.Year, "m:", t.Month, "d: ", t.Day, t.Hour, t.Minute,
        t.Second, t.ZoneOffset, t.Seconds(), "\n");        
    fmt.Println(time.LocalTime().Seconds());        
    retval = uint64(t.Seconds() * 1e9);
  }
  
  return retval;
}


/**
 * Opens a specified file and attempts to extract meta data.
 * There are two possibilities for metadata. Without either,
 * dates fallback to the modification date of the file and the 
 * the first line as the fallback.
 *
 * 1. The date is in a metadata segment at the top of the file as
 * defined for MetaMarkdown. This format consists of key: value with
 * a following blank line.
 *
 * 2. The data is contained in a comment as a sequence of numbers.
 * To keep this from being too inefficient, it must be found in the top
 * 5 lines.
 */
func rootThroughFileForMetadata(name string) (uint64, string) {
  fd, _ := os.Open(name, os.O_RDONLY, 0);
  // Collect the metadata in a struct?  
  
  // read a line of the file
  rd := bufio.NewReader(io.Reader(fd));

  lc := 0;
  inMetaData := false;

  var resultLine string;
  // var uint64 resultDate uint64;
  resultDate := uint64(0);
  
  for !inMetaData && lc < 5 {
    line, _ := rd.ReadString('\n');
    line = line[0:len(line)-1];
    
    if lc == 0 { resultLine = line; }
    // fmt.Print(line);
    // fmt.Print("\n");
    
    m1 := metadataMatcher.MatchStrings(line);
    m2 := commentDataMatcher.MatchStrings(line);
    if len(m1) > 0 {

      fmt.Print("matched for " + m1[1] + " <" + m1[2] + ">\n");
      if strings.ToLower(m1[1]) == "title" { resultLine = m1[2]; }
      // if strings.ToLower(m1[1]) == "date" { resultDate = parseDateFlexible(m1[2]); }
    } else if len(m2) > 0 {
      // fmt.Print("matched for  <" + m2[1] + ">\n");
      resultDate = parseDateCmdFmt(m2[1]);
    }
  
    lc++;
  }

  return resultDate, resultLine;
}

