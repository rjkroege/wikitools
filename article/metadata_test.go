package article

import (
    "io"
    "strings"
    "github.com/rjkroege/wikitools/testhelpers"
    "testing"
    "time"
)

var now = time.Now()
var never = time.Time{}

func Test_Makearticle(t *testing.T) {
    md := MetaData{ "foo.md", never, never, "", false, []string{}}
    if md.FormattedName() != "foo.html" {
        t.Errorf("expected  %s != to actual %s", "foo.html", md.Name)
    }    
}

func Test_UrlForName(t *testing.T) {
    md := MetaData{ "foo.md", never, never, "", false, []string{}}
    SetPathForContent("flimmer/blo");
    s := md.UrlForPath()
    if s != "file://flimmer/blo/foo.html" {
        t.Errorf("expected  %s != to actual %s", "file://flimmer/blo/foo.html", s)
    }
}

type pdSR struct {
    ex string
    err error
    in string
} 

func Test_ParseDateUnix(t *testing.T) {
    testdates := []pdSR {
        pdSR{ "Tue Sep 11 17:34:00 EDT 2012", nil,  "11 Sep 17:34:00 2012"},
        pdSR{ "Sat Oct 27 11:39:41 PDT 2012", nil,  "Sat Oct 27 11:39:41 PDT 2012"},
        pdSR{ "Wed Jun 15 08:24:39 EDT 2011", nil,  "2011/06/15 08:24:39"},
        pdSR{ "Tue Dec 27 17:46:16 EST 2011", nil,  "2011/12/27 17:46:16"},
        pdSR{ "Sun Mar 14 08:00:00 EST 2004", nil,  "200403140800"},
        pdSR{ "Tue Dec 11 17:34:00 EST 2012", nil,  "11 Dec 17:34:00 2012"},
        pdSR{ "Sat Dec  1 17:34:00 EST 2012", nil,  "1 Dec 17:34:00 2012"}     }

    for _, tu := range(testdates) {
        r, err := ParseDateUnix(tu.in)
        if tu.err != err {
            t.Errorf("invalid error value in test date %s", tu.in)
        }
        if tu.err == nil && tu.ex != r.Format(time.UnixDate) {
            t.Errorf("bad date: expected %s, received %s", tu.ex, r.Format(time.UnixDate))
        }
    }
}

const test_header_1 = 
`title: What I want
date: 2012/03/19 06:51:15

I need to figure out what I want. 
`

const test_header_2 = 
`title: What I want
date: 2012/03/19 06:51:15
tags: @journal

I need to figure out what I want. 
`
const test_header_3 = 
`date: 2012/03/19 06:51:15
title: What I want
tags: @journal

I need to figure out what I want. 
`
const test_header_4 = 
`I need
to figure out what I want. And code it.
`

const test_header_5 = 
`Date: 2012/03/19 06:51:15
Title: What I want
tags: @journal
`

const test_header_6 = 
`plastic: yes
Date: 2012/03/19 06:51:15
Tag: empty
Title: What I want
tags: @journal

I need to figure out what to code
`

type rtfSR struct {
    in string
    err error
    ex MetaData
}

// TODO(rjkroege): Enforce the handling of dates.
// Need to validate that the right thing happens here.
func Test_RootThroughFileForMetadata(t *testing.T) {
    /* General idea: create a constant string. Read from it., validate the resulting output. */

    realisticdate, _ := ParseDateUnix("1999/03/21 17:00:00")
    date, _ := ParseDateUnix("2012/03/19 06:51:15")
    testfiles := []rtfSR {
        rtfSR{ test_header_1, nil, MetaData{"", realisticdate, date , "What I want", true, []string{}}}, 
        rtfSR{ test_header_2, nil, MetaData{"", realisticdate, date , "What I want",  true, []string{}}}, 
        rtfSR{ test_header_3, nil, MetaData{"", realisticdate, date , "What I want", true, []string{}}}, 
        rtfSR{ test_header_4, nil, MetaData{"", realisticdate, never , "I need",  false, []string{}}},  
        rtfSR{ test_header_5, nil, MetaData{"", realisticdate, date , "What I want", true, []string{}}},
        rtfSR{ test_header_6, nil, MetaData{"", realisticdate, date , "What I want", true, []string{}}}  }

    for _, tu := range(testfiles) {
        md := MetaData{"", realisticdate, never, "", false, []string{}};
        rd := strings.NewReader(tu.in)
        md.RootThroughFileForMetadata(io.Reader(rd))

        // TODO(rjkroege): Add nicer String() on Metadata?   
        if !md.equals(&tu.ex) {
            t.Errorf("expected %s != actual %s", tu.ex, md)
        }
    }
}

func Test_PrettyDate(t *testing.T) {
    statdate, _ := ParseDateUnix("1999/03/21 17:00:00")
    tagdate, _ := ParseDateUnix("2012/03/19 06:51:15")

    md := MetaData{"", statdate, never , "What I want 0", false, []string{}}
    testhelpers.AssertString(t, "Sunday, Mar 21, 1999", md.PrettyDate())

    md = MetaData{"", statdate, tagdate , "What I want 0", true, []string{}}
    testhelpers.AssertString(t, "Monday, Mar 19, 2012", md.PrettyDate())
}

type tEdMd struct {
    err error
    result string
    md MetaData
}


const json1 = `{"link":"file:///url-here/1.html","start":"Sunday, Mar 21, 1999","title":"What I want 0"}`
const json2  = `{"link":"file:///url-here/2.html","start":"Monday, Mar 19, 2012","title":"What I want 0"}`

func Test_JsonDate(t *testing.T) {
    statdate, _ := ParseDateUnix("1999/03/21 17:00:00")
    tagdate, _ := ParseDateUnix("2012/03/19 06:51:15")
    SetPathForContent("/url-here")

    datas := []tEdMd {
        { nil, json1, MetaData{"1.md", statdate, never , "What I want 0", false, []string{}}},
        { nil, json2, MetaData{"2.md", statdate, tagdate , "What I want 0", true, []string{}}}}

    for _, m := range(datas) {
        b, e := m.md.MarshalJSON()
        if m.err != e {
            t.Errorf("error value wrong");
        }
        testhelpers.AssertString(t, m.result, string(b))
    }
}
