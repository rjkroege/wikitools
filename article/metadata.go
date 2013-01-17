/*
  Constructs the actual note. Does what the current python buildPage.py
  script does.
*/
package article;

import (
  // "fmt"
  "time"
)

// Clarify the purpose of the struct members.
// Note the use of the named fields for generating
// Timeline JSON.
type MetaData struct {
  Name string					`json:"-"`
  Url string					`json:"link"`
  DateFromStat time.Time		`json:"-"`
  DateFromMetadata time.Time	`json:"-"`
  Title string					`json:"title"`
  FinalDate string				`json:"start"`
  hadMetaData bool 			`json:"-"`
//  PrettyDate string				`json:"-"`
  SourcePath string			`json:"-"`
}

// Converts an article name into its name as a formatted object.s
func (md *MetaData) FormattedName() string {
  oname := md.Name[0:len(md.Name) - len(".md")] + ".html";
  return oname
}

// Constructs a URL path equivalent to the given source file.
func (md *MetaData) UrlForName(path string) string {
  // Prefix file:///<path>/fname.html
  md.Url = "file://" + path + "/" + md.FormattedName();
  return md.Url;
}

func (md *MetaData) SourceForName(path string) string {
  md.SourcePath = path + "/" + md.Name
  return md.SourcePath
}

// TODO(rjkroege): Add a constructor.
// TODO(rjkroege): Make your tests less brittle

func (md *MetaData) PrettyDate() string {
    const df = "Monday, Jan _2, 2006"
    if (!md.DateFromMetadata.IsZero()) {
        return md.DateFromMetadata.Format(df);
    }
    return md.DateFromStat.Format(df);
}

/*
func (md *MetaData) JsonDate(fd io.Writer) {
    const df = "Monday, Jan _2, 2006"
    if (!md.DateFromMetadata.IsZero()) {
        return md.DateFromMetadata.Format(time.RFC3339);
    }
    return md.DateFromMetadata.Format(df);
}
*/
