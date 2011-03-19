/*
  Metadata extraction  
  ; fn gogo  {make test}

*/

package article;

import (
  // "fmt";
  "os";
  "bufio";
  "io";
  "strings";
  "strconv";
  "regexp";
  "time";
)

var metadataMatcher = regexp.MustCompile("^([A-Za-z]*):[ \t]*(.*)$");
var commentDataMatcher = regexp.MustCompile("<!-- *([0-9]*) *-->");

/**
 * Attempts to parse the date sequence that I have used in multiple
 * files that consists of a string of digits.
 */
func parseDateCmdFmt(numericDate string) int64 {
  d0 := int64(0);
  t := time.LocalTime();
  var e os.Error;

  if len(numericDate) != 14 { return d0; }
  
  t.Year, e = strconv.Atoi64(numericDate[0:4]);
  if e != nil { return d0; }

  t.Month, e = strconv.Atoi(numericDate[4:6]);
  if e != nil { return d0; }

  t.Day, e = strconv.Atoi(numericDate[6:8]);
  if e != nil { return d0; }

  t.Hour, e = strconv.Atoi(numericDate[8:10]);
  if e != nil { return d0; }

  t.Minute, e = strconv.Atoi(numericDate[10:12]);
  if e != nil { return d0; }

  t.Second, e = strconv.Atoi(numericDate[12:14]);
  if e != nil { return d0; }
  
  return int64(t.Seconds() * 1e9);
}

/**
 * Attempts to parse the metadata of the file. I will require file
 * metadata to be in UNIX data format (because I can since there
 * is no legacy.)
 *
 * Returns the first time corresponding to the first data match.
 */
func parseDateUnix(numericDate string) int64 {
  dateFormats := [2]string {
      "Mon _2 Jan 2006 15:04:05 MST",
      "2006/01/02 15:04:05" };
  resultDate := parseDateCmdFmt(numericDate);
  if resultDate > int64(0) { return resultDate; }

  for _, df := range(dateFormats) {
    t, e := time.Parse(df, numericDate);
    if e == nil {
      resultDate = int64(t.Seconds() * 1e9);
      break;
    }
  }
  return resultDate;
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
func (md *MetaData) RootThroughFileForMetadata() {
  fd, _ := os.Open(md.Name, os.O_RDONLY, 0)
  rd := bufio.NewReader(io.Reader(fd))
  lc := 0
  inMetaData := false
  md.hadMetaData = false

  var resultLine string;
  resultDate := int64(0);

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
      if strings.ToLower(m1[1]) == "date" { resultDate = parseDateUnix(m1[2]); }
      md.hadMetaData = true
    } else if len(m2) > 0 {
      // fmt.Print("matched for  <" + m2[1] + ">\n");
      resultDate = parseDateCmdFmt(m2[1]);
    }
  
    if resultDate == int64(0) {
      // fmt.Print("zero resultData, trying to parse title");
      resultDate = parseDateUnix(resultLine);
    }
  
    lc++;
  }
  fd.Close();
  md.DateFromMetadata, md.Title = resultDate, resultLine;
}

