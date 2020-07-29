package article

import (
	"encoding/json"
	"fmt"
	"github.com/rjkroege/wikitools/bibtex"
	"log"
	"path/filepath"
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
	Name             string
	DateFromStat     time.Time
	DateFromMetadata time.Time
	Title            string
	Dynamicstring    string
	HadMetaData      bool
	tags             []string
	extraKeys        map[string]string

	// The path where the article would go in the date-based article
	// categorization.
	datepath string
}

// TODO(rjk): specify if name should be absolute or not.
func MakeMetaData(name string, statTime time.Time) *MetaData {
	return &MetaData{
		Name:             name,
		DateFromStat:     statTime,
		DateFromMetadata: time.Time{},
		Title:            "",
		Dynamicstring:    "",
		HadMetaData:      false,
		tags:             []string{},
		extraKeys:        map[string]string{},
	}
}

func NewArticle(name string, title string, tags []string) *MetaData {
	return &MetaData{
		Name:             name,
		DateFromStat:     time.Time{},
		DateFromMetadata: time.Now(),
		Title:            title,
		Dynamicstring:    "",
		HadMetaData:      false,
		tags:             tags,
		extraKeys:        map[string]string{},
	}
}

// Generate the string from the list of tags.
func (md *MetaData) Tagstring() string {
	return strings.Join(md.tags, " ")
}

var shortmonths = [...]string{
  	"Jan",
  	"Feb",
  	"Mar",
  	"Apr",
  	"May",
  	"Jun",
  	"Jul",
  	"Aug",
  	"Sep",
  	"Oct",
  	"Nov",
  	"Dec",
  }

// RelativeDateDirectory generates the name of the file in the structured
// date-based sorting.
func (md *MetaData) RelativeDateDirectory() string {
	t := md.PreferredDate()

	return filepath.Join(fmt.Sprintf("%d", t.Year()), fmt.Sprintf("%02d-%s", t.Month(), shortmonths[t.Month()-1]), fmt.Sprintf("%d", t.Day()))

}

// Generate the string of the extra keys.
func (md *MetaData) ExtraKeysString() string {
	result := make([]string, 0)
	for k, v := range md.extraKeys {
		result = append(result, k+":"+v)
	}
	sort.Strings(result)
	return strings.Join(result, ", ")
}

func (md *MetaData) ExtraKeys() map[string]string {
	return md.extraKeys
}

func (a *MetaData) equals(b *MetaData) bool {
	return a.Name == b.Name && a.DateFromStat == b.DateFromStat && a.DateFromMetadata == b.DateFromMetadata && a.Title == b.Title && a.HadMetaData == b.HadMetaData && a.Tagstring() == b.Tagstring() && a.ExtraKeysString() == b.ExtraKeysString()
}

// Converts an article name into its name as a formatted object.
func (md *MetaData) FormattedName() string {
	oname := md.Name[0:len(md.Name)-len(".md")] + ".html"
	return oname
}

// Constructs a URL path equivalent to the given source file.
func (md *MetaData) UrlForPath() string {
	return "file://" + pathForContent + "/" + md.FormattedName()
}

func (md *MetaData) SourceForPath() string {
	return pathForContent + "/" + md.Name
}

func (md *MetaData) PreferredDate() time.Time {
	if !md.DateFromMetadata.IsZero() {
		return md.DateFromMetadata
	}
	return md.DateFromStat
}

func (md *MetaData) PrettyDate() string {
	const df = "Monday, Jan _2, 2006"
	return md.PreferredDate().Format(df)
}

func (md *MetaData) DetailedDate() string {
	const df = "Mon _2 Jan 2006, 15:04:05 MST"
	return md.PreferredDate().Format(df)
}

type jsonmetadata struct {
	Link  string `json:"link"`
	Start string `json:"start"`
	Title string `json:"title"`
}

func (md *MetaData) MarshalJSON() ([]byte, error) {
	// const df = "Monday, Jan _2, 2006"
	jmd := jsonmetadata{md.UrlForPath(), md.PreferredDate().Format(time.RFC3339), md.Title}
	return json.Marshal(jmd)
}

// HaveBibTex returns true if this article has a BibTeX entry.
func (md *MetaData) HaveBibTex() bool {
	_, err := bibtex.ExtractBibTeXEntryType(md.tags)
	return err == nil
}

// BibTexEntry return a BibTex Entry for this article.
func (md *MetaData) BibTexEntry() string {
	s, err := bibtex.CreateBibTexEntry(md.tags, md.extraKeys)
	if err != nil {
		// TODO(rjkroege): Do something more rational with errors.
		log.Print("Problem with bibtex entry: " + err.Error())
	}
	return s
}
