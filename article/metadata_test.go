package article

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/rjkroege/wikitools/testhelpers"
	"github.com/rjkroege/wikitools/wiki"
)

var never = time.Time{}

func Test_Makearticle(t *testing.T) {
	md := NewArticleWithTime("foo.md", never, never, "", MdInvalid)
	if md.FormattedName() != "foo.html" {
		t.Errorf("expected  %s != to actual %s", "foo.html", md.filename)
	}
}

func Test_UrlForName(t *testing.T) {
	md := NewArticleWithTime("foo.md", never, never, "", MdInvalid)
	SetPathForContent("flimmer/blo")
	s := md.UrlForPath()
	if s != "file://flimmer/blo/foo.html" {
		t.Errorf("expected  %s != to actual %s", "file://flimmer/blo/foo.html", s)
	}
}

func Test_ExtraKeysString(t *testing.T) {
	m := MetaData{"", never, never, "", "", MdInvalid, []string{}, map[string]string{"a": "b"}, ""}
	testhelpers.AssertString(t, "a:b", m.ExtraKeysString())

	m = MetaData{"", never, never, "", "", MdInvalid, []string{}, map[string]string{"a": "b", "c": "d"}, ""}
	testhelpers.AssertString(t, "a:b, c:d", m.ExtraKeysString())

	m = MetaData{"", never, never, "", "", MdInvalid, []string{}, map[string]string{"c": "d", "a": "b"}, ""}
	testhelpers.AssertString(t, "a:b, c:d", m.ExtraKeysString())
}

type rtfSR struct {
	testname string
	in       string
	err      error
	ex       MetaData
}

// TODO(rjkroege): Enforce the handling of dates.
// Need to validate that the right thing happens here.
func Test_RootThroughFileForMetadata(t *testing.T) {
	realisticdate, _ := wiki.ParseDateUnix("1999/03/21 17:00:00")
	date, _ := wiki.ParseDateUnix("2012/03/19 06:51:15")
	testfiles := []rtfSR{
		{"Test_header_1", testhelpers.Test_header_1, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdLegacy, []string{}, map[string]string{}, ""}},
		{"Test_header_1_dash", testhelpers.Test_header_1_dash, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdIaWriter, []string{}, map[string]string{}, ""}},
		{"Test_header_2", testhelpers.Test_header_2, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdLegacy, []string{"journal"}, map[string]string{}, ""}},
		{"Test_header_3", testhelpers.Test_header_3, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdLegacy, []string{"journal"}, map[string]string{}, ""}},
		{"Test_header_4", testhelpers.Test_header_4, nil,
			MetaData{"", realisticdate, never, "I need", "", MdInvalid, []string{}, map[string]string{}, ""}},
		{"Test_header_5", testhelpers.Test_header_5, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdLegacy, []string{"journal"}, map[string]string{}, ""}},
		{"Test_header_6", testhelpers.Test_header_6, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdLegacy, []string{"journal"},
				map[string]string{"tag": "empty", "plastic": "yes"}, ""}},
		{"Test_header_6_dash", testhelpers.Test_header_6_dash, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdIaWriter, []string{"journal"},
				map[string]string{"tag": "empty", "plastic": "yes"}, ""}},
		{"Test_header_7", testhelpers.Test_header_7, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdLegacy,
				[]string{"journal", "fiddle"},
				map[string]string{"tag": "empty", "plastic": "yes"}, ""}},
		{"Test_header_8", testhelpers.Test_header_8, nil,
			MetaData{"", realisticdate, date, "What I want", "", MdLegacy,
				[]string{"journal", "hello", "bye"}, map[string]string{"tag": "empty", "plastic": "yes"}, ""}},
		{"Test_header_9", testhelpers.Test_header_9, nil,
			MetaData{"", realisticdate, date, "Business Korea", "", MdLegacy,
				[]string{"book"}, map[string]string{"bib-bibkey": "kenna97", "bib-author": "Peggy Kenna and Sondra Lacy", "bib-title": "Business Korea", "bib-publisher": "Passport Books", "bib-year": "1997"}, ""}},
		{"Test_header_9_dash", testhelpers.Test_header_9_dash, nil,
			MetaData{"", realisticdate, date, "Business Korea", "", MdIaWriter,
				[]string{"book"}, map[string]string{"bib-bibkey": "kenna97", "bib-author": "Peggy Kenna and Sondra Lacy", "bib-title": "Business Korea", "bib-publisher": "Passport Books", "bib-year": "1997"}, ""}},
		{"Test_header_10_dash", testhelpers.Test_header_10_dash, nil,
			MetaData{"", realisticdate, date, "Business Korea", "", MdIaWriter,
				[]string{"book", "business", "korea"}, map[string]string{"bib-bibkey": "kenna97", "bib-author": "Peggy Kenna and Sondra Lacy", "bib-title": "Business Korea", "bib-publisher": "Passport Books", "bib-year": "1997"}, ""}},
	}

	for _, tu := range testfiles {
		if !tu.ex.equals(&tu.ex) {
			t.Errorf("%s: equals has failed for %v", tu.testname, tu)
		}

		md := &MetaData{"", realisticdate, never, "", "", MdInvalid, []string{}, map[string]string{}, ""}
		rd := strings.NewReader(tu.in)
		md.RootThroughFileForMetadata(io.Reader(rd))

		// TODO(rjkroege): Add nicer String() on Metadata?
		if !md.equals(&tu.ex) {
			t.Errorf("%s: got %v, want %v\n", tu.testname, md.Dump(), (&tu.ex).Dump())
		}
	}
}

// TODO(rjk): might belong in wiki.
func Test_PrettyDate(t *testing.T) {
	statdate, _ := wiki.ParseDateUnix("1999/03/21 17:00:00")
	tagdate, _ := wiki.ParseDateUnix("2012/03/19 06:51:15")

	md := MetaData{"", statdate, never, "What I want 0", "", MdInvalid, []string{}, map[string]string{}, ""}
	testhelpers.AssertString(t, "Sunday, Mar 21, 1999", md.PrettyDate())

	md = MetaData{"", statdate, tagdate, "What I want 0", "", MdInvalid, []string{}, map[string]string{}, ""}
	testhelpers.AssertString(t, "Monday, Mar 19, 2012", md.PrettyDate())
}

type tEdMd struct {
	err    error
	result string
	md     MetaData
}

const json1 = `{"link":"file:///url-here/1.html","start":"1999-03-21T17:00:00-05:00","title":"What I want 0"}`
const json2 = `{"link":"file:///url-here/2.html","start":"2012-03-19T06:51:15-04:00","title":"What I want 0"}`

func Test_JsonDate(t *testing.T) {
	statdate, _ := wiki.ParseDateUnix("1999/03/21 17:00:00")
	tagdate, _ := wiki.ParseDateUnix("2012/03/19 06:51:15")
	SetPathForContent("/url-here")

	datas := []tEdMd{
		{nil, json1, MetaData{"1.md", statdate, never, "What I want 0", "", MdInvalid, []string{}, map[string]string{}, ""}},
		{nil, json2, MetaData{"2.md", statdate, tagdate, "What I want 0", "", MdInvalid, []string{}, map[string]string{}, ""}}}

	for _, m := range datas {
		b, e := m.md.MarshalJSON()
		if m.err != e {
			t.Errorf("error value wrong")
		}
		testhelpers.AssertString(t, m.result, string(b))
	}
}

func TestRelativeDateDirectory(t *testing.T) {
	statdate, _ := wiki.ParseDateUnix("1999/03/21 17:00:00")
	tagdate, _ := wiki.ParseDateUnix("2012/03/19 06:51:15")
	otherdate, _ := wiki.ParseDateUnix("2012/03/03 06:51:15")

	for i, tv := range []struct {
		in   *MetaData
		want string
	}{
		{
			in:   &MetaData{"", statdate, never, "What I want 0", "", MdInvalid, []string{}, map[string]string{}, ""},
			want: "1999/03-Mar/21",
		},
		{
			in:   &MetaData{"", statdate, tagdate, "What I want 0", "", MdInvalid, []string{}, map[string]string{}, ""},
			want: "2012/03-Mar/19",
		},
		{
			in:   &MetaData{"", statdate, otherdate, "What I want 0", "", MdInvalid, []string{}, map[string]string{}, ""},
			want: "2012/03-Mar/03",
		},
	} {
		if got, want := tv.in.RelativeDateDirectory(), tv.want; got != want {
			t.Errorf("]%d] got %#v want %#v", i, got, want)
		}
	}
}
