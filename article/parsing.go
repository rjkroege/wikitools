package article;

import (
  "bufio"
  "io"
  "regexp"
  "strings"
  "time"
  // "fmt";  // Need for debugging printfs.
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

func parse(layout, value string) (time.Time, error) {
     return time.Parse(layout, value)
}

/**
 * Attempts to parse the metadata of the file. I will require file
 * metadata to be in UNIX date format (because I can since there
 * is no legacy.)
 *
 * Returns the first time corresponding to the first data match.
 */
func ParseDateUnix(ds string) (t time.Time, err error)  {
  timeformats := []string {
    time.UnixDate,
    lsdate,
    lstring,
    slashdate,
    sstring,
    unixlike };

  for _, fs := range(timeformats) {
    t, err = parse(fs, ds);
    if err == nil { return }
  }

 for _, fs := range(timeformats) {
    t, err  = parse(fs, ds + " EDT");
    if err == nil && t.Location().String() == "Local" { return }
    
    t, err = parse(fs, ds + " EST");
    if err == nil && t.Location().String() == "Local" { return }
  }
  return;
}

func trim(line string) string {
    if len(line) > 0 {
        return  line[0:len(line)-1]
    }
  return line
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
  rd := bufio.NewReader(reader)
  lc := 0
  md.HadMetaData = false

  var date time.Time;
  var de error;

  for lc < 5 || md.HadMetaData {
    line, e := rd.ReadString('\n');
    if e != nil || md.HadMetaData  && line == "\n" {
        // Return e.
        break;
    }
    line = trim(line)

    if lc == 0 { 
      md.Title  = line
    }

    // fmt.Print("running regexp matcher...\n")
    m1 := metadataMatcher.FindStringSubmatch(line);
    m2 := commentDataMatcher.FindStringSubmatch(line);
    if len(m1) > 0 {
      s  := strings.ToLower(m1[1])
      if s  == "title" {
        md.Title  = m1[2];
      } else if s == "date" {
        date, de = ParseDateUnix(strings.TrimSpace(m1[2]));
      } else if s == "tags" {
        for _, u := range strings.Split(strings.TrimSpace(m1[2]), " ") {
          if (u != "") {
            md.tags = append(md.tags, u)
          }
        }
      } else {
        md.extraKeys[s] = strings.TrimSpace(m1[2])
      }
      md.HadMetaData = true
    } else if len(m2) > 0 {
      // fmt.Print("matched for  <" + m2[1] + ">\n");
      date, de = ParseDateUnix(m2[1]);
    }

    // I have no test that actually enforces that this is valid.
    // push to a helper  
    if de != nil || date.IsZero() {
      //fmt.Print("date is zero, trying whole resultLine: <" + resultLine + ">\n");
      date, de  = ParseDateUnix(md.Title);
    }
    lc++;
  }
  md.DateFromMetadata = date
}

