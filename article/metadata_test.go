package article

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/rjkroege/wikitools/testhelpers"
)

var never = time.Time{}

// NewArticleTest makes an article for testing.
func NewArticleTest(name string, stat time.Time, meta time.Time, title string, has bool) *MetaData {
	return &MetaData{
		Name:             name,
		DateFromStat:     stat,
		DateFromMetadata: meta,
		Title:            title,
		HadMetaData:      has,
	}
}

func Test_Makearticle(t *testing.T) {
	md := NewArticleTest("foo.md", never, never, "", false)
	if md.FormattedName() != "foo.html" {
		t.Errorf("expected  %s != to actual %s", "foo.html", md.Name)
	}
}

func Test_UrlForName(t *testing.T) {
	md := NewArticleTest("foo.md", never, never, "", false)
	SetPathForContent("flimmer/blo")
	s := md.UrlForPath()
	if s != "file://flimmer/blo/foo.html" {
		t.Errorf("expected  %s != to actual %s", "file://flimmer/blo/foo.html", s)
	}
}

func Test_ExtraKeysString(t *testing.T) {
	m := MetaData{"", never, never, "", "", false, []string{}, map[string]string{"a": "b"}, ""}
	testhelpers.AssertString(t, "a:b", m.ExtraKeysString())

	m = MetaData{"", never, never, "", "", false, []string{}, map[string]string{"a": "b", "c": "d"}, ""}
	testhelpers.AssertString(t, "a:b, c:d", m.ExtraKeysString())

	m = MetaData{"", never, never, "", "", false, []string{}, map[string]string{"c": "d", "a": "b"}, ""}
	testhelpers.AssertString(t, "a:b, c:d", m.ExtraKeysString())
}

type pdSR struct {
	ex  string
	err error
	in  string
}

func Test_ParseDateUnix(t *testing.T) {
	testdates := []pdSR{
		{"Mon Mar 11 00:00:00 EDT 2013", nil, "Monday, Mar 11, 2013"},
		{"Tue Sep 11 17:34:00 EDT 2012", nil, "11 Sep 17:34:00 2012"},
		{"Sat Oct 27 11:39:41 PDT 2012", nil, "Sat Oct 27 11:39:41 PDT 2012"},
		{"Wed Jun 15 08:24:39 EDT 2011", nil, "2011/06/15 08:24:39"},
		{"Tue Dec 27 17:46:16 EST 2011", nil, "2011/12/27 17:46:16"},
		{"Sun Mar 14 08:00:00 EST 2004", nil, "200403140800"},
		{"Tue Dec 11 17:34:00 EST 2012", nil, "11 Dec 17:34:00 2012"},
		{"Fri Jun 14 07:25:48 EDT 2013", nil, "Fri 14 Jun 2013, 07:25:48 EDT"},
		{"Sat Dec  1 17:34:00 EST 2012", nil, "1 Dec 17:34:00 2012"},
		{"Tue Sep  5 11:14:03 PDT 2006", nil, "Tue Sep  5 11:14:03 PDT 2006"},
		{"Tue Feb  5 08:52:22 -0700 2019", nil, "2019-02-05 08:52:22.000000000 -0700"},
	}

	for _, tu := range testdates {
		r, err := ParseDateUnix(tu.in)
		if tu.err != err {
			t.Errorf("invalid error value in test date %s", tu.in)
		}
		if tu.err == nil && tu.ex != r.Format(time.UnixDate) {
			t.Errorf("bad date: expected %s, received %s", tu.ex, r.Format(time.UnixDate))
		}
	}
}

const test_header_1 = `title: What I want
date: 2012/03/19 06:51:15

I need to figure out what I want. 
`

const test_header_1_dash = `---
title: What I want
date: 2012/03/19 06:51:15
---

I need to figure out what I want. 
`

const test_header_2 = `title: What I want
date: 2012/03/19 06:51:15
tags: @journal

I need to figure out what I want. 
`
const test_header_3 = `date: 2012/03/19 06:51:15
title: What I want
tags: @journal

I need to figure out what I want. 
`
const test_header_4 = `I need
to figure out what I want. And code it.
`

const test_header_5 = `Date: 2012/03/19 06:51:15
Title: What I want
tags: @journal
`

const test_header_6 = `plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal

I need to figure out what to code
`

const test_header_6_dash = `---
plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal
---

I need to figure out what to code
`

const test_header_7 = `plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal   @fiddle

I need to figure out what to code
`
const test_header_8 = `plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal  @hello     @bye

I need to figure out what to code
`

const test_header_9 = `title: Business Korea
date: 2012/03/19 06:51:15
tags: @book
bib-bibkey: kenna97
bib-author: Peggy Kenna and Sondra Lacy
bib-title: Business Korea
bib-publisher: Passport Books
bib-year:  1997

Business book.
`

const test_header_9_dash = `---
title: Business Korea
date: 2012/03/19 06:51:15
tags: @book
bib-bibkey: kenna97
bib-author: Peggy Kenna and Sondra Lacy
bib-title: Business Korea
bib-publisher: Passport Books
bib-year:  1997
---

Business book.
`

type rtfSR struct {
	testname string
	in       string
	err      error
	ex       MetaData
}

// TODO(rjkroege): Enforce the handling of dates.
// Need to validate that the right thing happens here.
func Test_RootThroughFileForMetadata(t *testing.T) {
	/* General idea: create a constant string. Read from it., validate the resulting output. */

	realisticdate, _ := ParseDateUnix("1999/03/21 17:00:00")
	date, _ := ParseDateUnix("2012/03/19 06:51:15")
	testfiles := []rtfSR{
		{"test_header_1", test_header_1, nil,
			MetaData{"", realisticdate, date, "What I want", "", true, []string{}, map[string]string{}, ""}},
		{"test_header_1_dash", test_header_1_dash, nil,
			MetaData{"", realisticdate, date, "What I want", "", true, []string{}, map[string]string{}, ""}},
		{"test_header_2", test_header_2, nil,
			MetaData{"", realisticdate, date, "What I want", "", true, []string{"@journal"}, map[string]string{}, ""}},
		{"test_header_3", test_header_3, nil,
			MetaData{"", realisticdate, date, "What I want", "", true, []string{"@journal"}, map[string]string{}, ""}},
		{"test_header_4", test_header_4, nil,
			MetaData{"", realisticdate, never, "I need", "", false, []string{}, map[string]string{}, ""}},
		{"test_header_5", test_header_5, nil,
			MetaData{"", realisticdate, date, "What I want", "", true, []string{"@journal"}, map[string]string{}, ""}},
		{"test_header_6", test_header_6, nil,
			MetaData{"", realisticdate, date, "What I want", "", true, []string{"@journal"},
				map[string]string{"tag": "empty", "plastic": "yes"}, ""}},
		{"test_header_6_dash", test_header_6_dash, nil,
			MetaData{"", realisticdate, date, "What I want", "", true, []string{"@journal"},
				map[string]string{"tag": "empty", "plastic": "yes"}, ""}},
		{"test_header_7", test_header_7, nil,
			MetaData{"", realisticdate, date, "What I want", "", true,
				[]string{"@journal", "@fiddle"},
				map[string]string{"tag": "empty", "plastic": "yes"}, ""}},
		{"test_header_8", test_header_8, nil,
			MetaData{"", realisticdate, date, "What I want", "", true,
				[]string{"@journal", "@hello", "@bye"}, map[string]string{"tag": "empty", "plastic": "yes"}, ""}},
		{"test_header_9", test_header_9, nil,
			MetaData{"", realisticdate, date, "Business Korea", "", true,
				[]string{"@book"}, map[string]string{"bib-bibkey": "kenna97", "bib-author": "Peggy Kenna and Sondra Lacy", "bib-title": "Business Korea", "bib-publisher": "Passport Books", "bib-year": "1997"}, ""}},
		{"test_header_9_dash", test_header_9_dash, nil,
			MetaData{"", realisticdate, date, "Business Korea", "", true,
				[]string{"@book"}, map[string]string{"bib-bibkey": "kenna97", "bib-author": "Peggy Kenna and Sondra Lacy", "bib-title": "Business Korea", "bib-publisher": "Passport Books", "bib-year": "1997"}, ""}},
	}

	for _, tu := range testfiles {
		if !tu.ex.equals(&tu.ex) {
			t.Errorf("%s: equals has failed for %v", tu.testname, tu)
		}

		md := &MetaData{"", realisticdate, never, "", "", false, []string{}, map[string]string{}, ""}
		rd := strings.NewReader(tu.in)
		md.RootThroughFileForMetadata(io.Reader(rd))

		// TODO(rjkroege): Add nicer String() on Metadata?
		if !md.equals(&tu.ex) {
			t.Errorf("%s: got %v, want %v\n", tu.testname, md.Dump(), (&tu.ex).Dump())
		}
	}
}

func Test_PrettyDate(t *testing.T) {
	statdate, _ := ParseDateUnix("1999/03/21 17:00:00")
	tagdate, _ := ParseDateUnix("2012/03/19 06:51:15")

	md := MetaData{"", statdate, never, "What I want 0", "", false, []string{}, map[string]string{}, ""}
	testhelpers.AssertString(t, "Sunday, Mar 21, 1999", md.PrettyDate())

	md = MetaData{"", statdate, tagdate, "What I want 0", "", true, []string{}, map[string]string{}, ""}
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
	statdate, _ := ParseDateUnix("1999/03/21 17:00:00")
	tagdate, _ := ParseDateUnix("2012/03/19 06:51:15")
	SetPathForContent("/url-here")

	datas := []tEdMd{
		{nil, json1, MetaData{"1.md", statdate, never, "What I want 0", "", false, []string{}, map[string]string{}, ""}},
		{nil, json2, MetaData{"2.md", statdate, tagdate, "What I want 0", "", true, []string{}, map[string]string{}, ""}}}

	for _, m := range datas {
		b, e := m.md.MarshalJSON()
		if m.err != e {
			t.Errorf("error value wrong")
		}
		testhelpers.AssertString(t, m.result, string(b))
	}
}

func TestRelativeDateDirectory(t *testing.T) {
	statdate, _ := ParseDateUnix("1999/03/21 17:00:00")
	tagdate, _ := ParseDateUnix("2012/03/19 06:51:15")
	otherdate, _ := ParseDateUnix("2012/03/03 06:51:15")

	for i, tv := range []struct {
		in   *MetaData
		want string
	}{
		{
			in:   &MetaData{"", statdate, never, "What I want 0", "", false, []string{}, map[string]string{}, ""},
			want: "1999/03-Mar/21",
		},
		{
			in:   &MetaData{"", statdate, tagdate, "What I want 0", "", true, []string{}, map[string]string{}, ""},
			want: "2012/03-Mar/19",
		},
		{
			in:   &MetaData{"", statdate, otherdate, "What I want 0", "", true, []string{}, map[string]string{}, ""},
			want: "2012/03-Mar/3",
		},
	} {
		if got, want := tv.in.RelativeDateDirectory(), tv.want; got != want {
			t.Errorf("]%d] got %#v want %#v", i, got, want)
		}
	}
}
