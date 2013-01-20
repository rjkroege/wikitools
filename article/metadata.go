/*
  Constructs the actual note. Does what the current python buildPage.py
  script does.
*/
package article;

import (
     // "fmt"
    "encoding/json"
    "time"
)

// Clarify the purpose of the struct members.
// Note the use of the named fields for generating
// Timeline JSON.
type MetaData struct {
  Name string
  Url string
  DateFromStat time.Time
  DateFromMetadata time.Time
  Title string
  hadMetaData bool
  SourcePath string
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

func (md *MetaData) PreferredDate() time.Time {
    if (!md.DateFromMetadata.IsZero()) {
        return md.DateFromMetadata
    }
    return md.DateFromStat
}

func (md *MetaData) PrettyDate() string {
    const df = "Monday, Jan _2, 2006"
    return md.PreferredDate().Format(df);
}

type jsonmetadata struct {
    Link string				`json:"link"`
    Start string				`json:"start"`
    Title string				`json:"title"`
}

func (md *MetaData) MarshalJSON() ([]byte, error) {
    const df = "Monday, Jan _2, 2006"
    jmd := jsonmetadata{ md.Url, md.PreferredDate().Format(df), md.Title }
    return json.Marshal(jmd)    
}
