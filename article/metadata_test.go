package article

import (
    "io"
    "strings"
    "testing"
    "time"
)

func Test_Makearticle(t *testing.T) {
    md := MetaData{ "foo.md", "", time.Now(), time.Now(), "", "", false, "", ""}
    if md.FormattedName() != "foo.html" {
        t.Errorf("expected  %s != to actual %s", "foo.html", md.Name)
    }    
}

func Test_UrlForName(t *testing.T) {
    md := MetaData{ "foo.md", "", time.Now(), time.Now(), "", "", false, "", ""}
    s := md.UrlForName("flimmer/blo")
    if s != "file://flimmer/blo/foo.html" {
        t.Errorf("expected  %s != to actual %s", "file://flimmer/blo/foo.html", s)
    }
}

type pdSR struct {
    ex string
    err error
    in string
} 

func Test_parseDateUnix(t *testing.T) {
    testdates := []pdSR {
        pdSR{ "Tue Sep 11 17:34:00 EDT 2012", nil,  "11 Sep 17:34:00 2012"},
        pdSR{ "Sat Oct 27 11:39:41 PDT 2012", nil,  "Sat Oct 27 11:39:41 PDT 2012"},
        pdSR{ "Wed Jun 15 08:24:39 EDT 2011", nil,  "2011/06/15 08:24:39"},
        pdSR{ "Tue Dec 27 17:46:16 EST 2011", nil,  "2011/12/27 17:46:16"},
        pdSR{ "Sun Mar 14 08:00:00 EST 2004", nil,  "200403140800"},
        pdSR{ "Tue Dec 11 17:34:00 EST 2012", nil,  "11 Dec 17:34:00 2012"},
        pdSR{ "Sat Dec  1 17:34:00 EST 2012", nil,  "1 Dec 17:34:00 2012"}     }

    for _, tu := range(testdates) {
        r, err := parseDateUnix(tu.in)
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
func Test_RootThroughFileForMetadata(t *testing.T) {
    /* General idea: create a constant string. Read from it., validate the resulting output. */

    date, _ := parseDateUnix("2012/03/19 06:51:15")
    testfiles := []rtfSR {
        rtfSR{ test_header_1, nil, MetaData{"", "", time.Time{}, date , "What I want", "", true, "", ""}}, 
        rtfSR{ test_header_2, nil, MetaData{"", "", time.Time{}, date , "What I want", "", true, "", ""}}, 
        rtfSR{ test_header_3, nil, MetaData{"", "", time.Time{}, date , "What I want", "", true, "", ""}}, 
        rtfSR{ test_header_4, nil, MetaData{"", "", time.Time{}, time.Time{} , "I need", "", false, "", ""}},  
        rtfSR{ test_header_5, nil, MetaData{"", "", time.Time{}, date , "What I want", "", true, "", ""}},
        rtfSR{ test_header_6, nil, MetaData{"", "", time.Time{}, date , "What I want", "", true, "", ""}}  }

    for _, tu := range(testfiles) {
        md := MetaData{"", "", time.Time{}, time.Time{}, "", "", false, "", ""};
        rd := strings.NewReader(tu.in)
        md.RootThroughFileForMetadata(io.Reader(rd))

        // TODO(rjkroege): Add nicer String() on Metadata?   
        if md !=  tu.ex {
            t.Errorf("expected %s != actual %s", tu.ex, md)
        }
    }

}
