/*
  Metadata extraction  
  ; fn gogo  {make test}

*/

package article;

import (
  "bufio"
  "io"
  "regexp"
  "strings"
  "time"
//  "fmt";  // Need for debugging printfs.
)

var metadataMatcher = regexp.MustCompile("^([A-Za-z]*):[ \t]*(.*)$");
var commentDataMatcher = regexp.MustCompile("<!-- *([0-9]*) *-->");

const (
lsdate = "_2 Jan 15:04:05 2006 MST"
lstring = "20060102150405 MST";
slashdate = "2006/01/02 15:04:05 MST"
sstring = "200601021504 MST";
unixlike = "Mon _2 Jan 2006 15:04:05 MST"
)

/**
 * Attempts to parse the metadata of the file. I will require file
 * metadata to be in UNIX date format (because I can since there
 * is no legacy.)
 *
 * Returns the first time corresponding to the first data match.
 */
func parseDateUnix(ds string) (t time.Time, err error)  {
  // fmt.Print("time string <" + ds + ">\n");

  timeformats := []string {
    time.UnixDate,
    lsdate,
    lstring,
    slashdate,
    sstring,
    unixlike };

  for _, fs := range(timeformats) {
    t, err = time.Parse(fs, ds);
    if err == nil { return }
  }

 for _, fs := range(timeformats) {
    t, err = time.Parse(fs, ds + " EDT");
    // fmt.Print("<" + fs + "> for <" + ds + " EDT> gives location: " + t.Location().String());
    if err == nil && t.Location().String() == "Local" { return }
    
    t, err = time.Parse(fs, ds + " EST");
    // fmt.Print("<" + fs + "> for <" + ds + " EST> gives location: " + t.Location().String());
    if err == nil && t.Location().String() == "Local" { return }
  }
  return;
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
 * 2. The date is contained in a comment as a sequence of numbers.
 * To keep this from being too inefficient, it must be found in the top
 * 5 lines.
 */
func (md *MetaData) RootThroughFileForMetadata(reader io.Reader) {
  // fmt.Print("\nfile: " + md.Name + "\n");
  rd := bufio.NewReader(reader)
  lc := 0
  inMetaData := false
  md.hadMetaData = false

  var resultLine string;
  var date time.Time;
  var de error;

  for !inMetaData && lc < 5 {
    line, _ := rd.ReadString('\n');
    if len(line) > 0 { line = line[0:len(line)-1]; }
    
    if lc == 0 { resultLine = line; }
     // fmt.Print(line);
     // fmt.Print("\n");

    // fmt.Print("running regexp matcher...\n")
    m1 := metadataMatcher.FindStringSubmatch(line);
    m2 := commentDataMatcher.FindStringSubmatch(line);
    if len(m1) > 0 {
      // fmt.Print("matched for " + m1[1] + " <" + m1[2] + ">\n");
      if strings.ToLower(m1[1]) == "title" { resultLine = m1[2]; }
      if strings.ToLower(m1[1]) == "date" {
        date, de = parseDateUnix(strings.TrimSpace(m1[2]));
      }
      md.hadMetaData = true
    } else if len(m2) > 0 {
      // fmt.Print("matched for  <" + m2[1] + ">\n");
      date, de = parseDateUnix(m2[1]);
    }
  
    if de != nil || date.IsZero() {
      //fmt.Print("date is zero, trying whole resultLine: <" + resultLine + ">\n");
      date, de  = parseDateUnix(resultLine);
    }
    lc++;
  }
  md.DateFromMetadata, md.Title = date, resultLine;
}

