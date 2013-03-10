package article;

import (
     // "fmt"
    "encoding/json"
    "sort"
    "strings"
    "time"
)

var pathForContent = ""

func SetPathForContent(p string) {
    pathForContent = p
}

// Prefer lower-case fields.
type MetaData struct {
  Name string
  DateFromStat time.Time
  DateFromMetadata time.Time
  Title string
  HadMetaData bool
  tags []string
  extraKeys map[string]string
}

func NewMetaData(name string, statTime time.Time) (*MetaData) {
    return &MetaData{ name, statTime, time.Time{}, "", false, []string{}, map[string]string{}}
}

func NewArticle(name string, title string, tags []string) *MetaData {
    return &MetaData{name, time.Time{}, time.Now(), title, false, tags, map[string]string{}}
}

// Use only from tests. (How could I enforce this?)
func NewArticleTest(name string, stat time.Time, meta time.Time, title string, has bool) *MetaData {
    return &MetaData{name, stat, meta , title, has, []string{}, map[string]string{}}
}

// Generate the string from the list of tags.
func (md *MetaData) Tagstring() string {
    return strings.Join(md.tags,  " ")
}

// Generate the string of the extra keys.
func (md *MetaData) ExtraKeysString() string {
    result := make([]string, 0)
    for k, v := range md.extraKeys {
        result = append(result, k + ":" + v)        
    }
    sort.Strings(result)
    return strings.Join(result, ", ")
}

func (md *MetaData) ExtraKeys() map[string]string {
    return md.extraKeys
}

func (a *MetaData)  equals(b *MetaData) bool {
    return a.Name == b.Name &&  a.DateFromStat == b.DateFromStat &&  a.DateFromMetadata == b.DateFromMetadata && a.Title == b.Title &&  a.HadMetaData == b.HadMetaData &&  a.Tagstring() == b.Tagstring() && a.ExtraKeysString() == b.ExtraKeysString()
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

func (md *MetaData) DetailedDate() string {
    const df = "Mon _2 Jan 2006, 15:04:05 MST"
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
