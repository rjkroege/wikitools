package article

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rjkroege/wikitools/bibtex"
	"github.com/rjkroege/wikitools/wiki"
)

var pathForContent = ""

func SetPathForContent(p string) {
	pathForContent = p
}

const (
	MdInvalid = iota
	MdLegacy
	MdIaWriter
	MdModern
)

var Metadatanametable = [...]string{
	"MdInvalid ",
	"MdLegacy",
	"MdIaWriter",
	"MdModern",
}

type MetadataType int

// Prefer lower-case fields.
type MetaData struct {
	filename         string
	DateFromStat     time.Time
	DateFromMetadata time.Time
	Title            string
	Dynamicstring    string
	mdtype           MetadataType
	Tags             []string
	extraKeys        map[string]string

	// The path where the article would go in the date-based article
	// categorization.
	datepath string
}

func MakeMetaData(name string, statTime time.Time) *MetaData {
	return &MetaData{
		filename:         name,
		DateFromStat:     statTime,
		DateFromMetadata: time.Time{},
		extraKeys:        map[string]string{},
	}
}

// TODO(rjk): filenamechange NewArticleWithTime -> NewArticle, NewArticle -> NewArticleDefaultTimes
// NewArticleTest makes an article for testing.
func NewArticleWithTime(name string, stat time.Time, meta time.Time, title string, mdtype MetadataType) *MetaData {
	return &MetaData{
		filename:         name,
		DateFromStat:     stat,
		DateFromMetadata: meta,
		Title:            title,
		mdtype:           mdtype,
		extraKeys:        map[string]string{},
	}
}

func NewArticle(name string, title string, tags []string) *MetaData {
	return &MetaData{
		filename:         name,
		DateFromStat:     time.Time{},
		DateFromMetadata: time.Now(),
		Title:            title,
		Tags:             tags,
		extraKeys:        map[string]string{},
	}
}

// TODO(rjk): Remove this.
func (md *MetaData) Name() string {
	// TODO(rjk): filenamechange make sure that every call to this thinks that it's the filename
	// TODO(rjk): conceivably this function can be removed eventually.
	return md.filename
}

// FileName is an accessor for the filename of this journal article.
// A complete journal article path will be filepath.Join(config.Basepath, RelativeDataDirectory, FileName)
func (md *MetaData) FileName() string {
	return md.filename
}

// PreferredFileName is the preferred name for this journal article
// without an extension.
func (md *MetaData) PreferredFileName(settings *wiki.Settings) string {
	ns := md.filename
	if md.Title != "" {
		ns = md.Title
	} else {
		// Strip the extension.
		ns = ns[0 : len(ns)-len(filepath.Ext(ns))]
		// TODO(rjk): Strip the uniquing extension here
	}
	return wiki.ValidName(ns)
}

// Type is an accessor for the type (e.g. version/vintage) of this MetaData
func (md *MetaData) Type() MetadataType {
	return md.mdtype
}

// Tagstring generates the string from the list of tags.
func (md *MetaData) Tagstring() string {
	ta := make([]string, 0, len(md.Tags))
	for _, v := range md.Tags {
		ta = append(ta, "#"+v)
	}
	return strings.Join(ta, " ")
}

// HadMetaData returns true if the file has metadata of some kind.
// TODO(rjk): This exists to support legacy code. Remove when possible.
func (md *MetaData) HadMetaData() bool {
	return md.mdtype != MdInvalid
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

func RelativeDateDirectoryForTime(t time.Time) string {
	return filepath.Join(fmt.Sprintf("%d", t.Year()), fmt.Sprintf("%02d-%s", t.Month(), shortmonths[t.Month()-1]), fmt.Sprintf("%02d", t.Day()))
}

// RelativeDateDirectory generates the name of the file in the structured
// date-based sorting.
func (md *MetaData) RelativeDateDirectory() string {
	return RelativeDateDirectoryForTime(md.PreferredDate())
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
	return a.filename == b.filename && a.DateFromStat == b.DateFromStat && a.DateFromMetadata == b.DateFromMetadata && a.Title == b.Title && a.mdtype == b.mdtype && a.Tagstring() == b.Tagstring() && a.ExtraKeysString() == b.ExtraKeysString()
}

// Converts an article name into its name as a formatted object.
// TODO(rjk): filenamechange
func (md *MetaData) FormattedName() string {
	oname := md.filename[0:len(md.filename)-len(".md")] + ".html"
	return oname
}

// Constructs a URL path equivalent to the given source file.
// TODO(rjk): filenamechange This code implies that FormattedName is not an absolute path
func (md *MetaData) UrlForPath() string {
	return "file://" + pathForContent + "/" + md.FormattedName()
}

// TODO(rjk): filenamechange
func (md *MetaData) SourceForPath() string {
	return pathForContent + "/" + md.filename
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
	return DetailedDateImpl(md.PreferredDate())
}

// Conceivably, this utility function belongs in wiki.
func DetailedDateImpl(d time.Time) string {
	const df = "Mon _2 Jan 2006, 15:04:05 MST"
	return d.Format(df)
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
	_, err := bibtex.ExtractBibTeXEntryType(md.Tags)
	return err == nil
}

// BibTexEntry return a BibTex Entry for this article.
func (md *MetaData) BibTexEntry() string {
	s, err := bibtex.CreateBibTexEntry(md.Tags, md.extraKeys)
	if err != nil {
		// TODO(rjkroege): Do something more rational with errors.
		log.Print("Problem with bibtex entry: " + err.Error())
	}
	return s
}
