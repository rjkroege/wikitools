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

var pathForContent = ""

func SetPathForContent(p string) {
    pathForContent = p
}

type MetaData struct {
  Name string
  DateFromStat time.Time
  DateFromMetadata time.Time
  Title string
  hadMetaData bool
}

func NewMetaData(name string, statTime time.Time) (*MetaData) {
    return &MetaData{ name, statTime, time.Time{}, "", false }
}

// Converts an article name into its name as a formatted object.
func (md *MetaData) FormattedName() string {
  oname := md.Name[0:len(md.Name) - len(".md")] + ".html";
  return oname
}

// Constructs a URL path equivalent to the given source file.
func (md *MetaData) UrlForPath() string {
  return "file://" + pathForContent + "/" + md.FormattedName();
}

func (md *MetaData) SourceForPath() string {
  return pathForContent + "/" + md.Name
}

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
    jmd := jsonmetadata{ md.UrlForPath(), md.PreferredDate().Format(df), md.Title }
    return json.Marshal(jmd)    
}
