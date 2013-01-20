package article

import (
    "bytes"
    "io"
    "strings"
    "testing"
    "time"
)

var now = Date{time.Now()}
var never = Date{time.Time{}}

func Test_Makearticle(t *testing.T) {
    md := MetaData{ "foo.md", "", never, never, "", false, ""}
    if md.FormattedName() != "foo.html" {
        t.Errorf("expected  %s != to actual %s", "foo.html", md.Name)
    }    
}

func Test_UrlForName(t *testing.T) {
    md := MetaData{ "foo.md", "", never, never, "", false,""}
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
// Need to validate that the right thing happens here.
func Test_RootThroughFileForMetadata(t *testing.T) {
    /* General idea: create a constant string. Read from it., validate the resulting output. */

    realisticdate, _ := parseDateUnix("1999/03/21 17:00:00")
    date, _ := parseDateUnix("2012/03/19 06:51:15")
    testfiles := []rtfSR {
        rtfSR{ test_header_1, nil, MetaData{"", "", realisticdate, date , "What I want", true, ""}}, 
        rtfSR{ test_header_2, nil, MetaData{"", "", realisticdate, date , "What I want",  true, ""}}, 
        rtfSR{ test_header_3, nil, MetaData{"", "", realisticdate, date , "What I want", true, ""}}, 
        rtfSR{ test_header_4, nil, MetaData{"", "", realisticdate, never , "I need",  false, ""}},  
        rtfSR{ test_header_5, nil, MetaData{"", "", realisticdate, date , "What I want", true, ""}},
        rtfSR{ test_header_6, nil, MetaData{"", "", realisticdate, date , "What I want", true, ""}}  }

    for _, tu := range(testfiles) {
        md := MetaData{"", "", realisticdate, never, "", false, ""};
        rd := strings.NewReader(tu.in)
        md.RootThroughFileForMetadata(io.Reader(rd))

        // TODO(rjkroege): Add nicer String() on Metadata?   
        if md !=  tu.ex {
            t.Errorf("expected %s != actual %s", tu.ex, md)
        }
    }
}

type closedCount int 

func (cc *closedCount) Close()  error {
    *cc += 1
    return nil
}

type mockReadCloser struct {
    closedCount
    *strings.Reader
    name string
}

type mockWriteCloser struct {
    closedCount
    *bytes.Buffer
    name string
}

type mockSystem struct {
    input string
    readfiles  []*mockReadCloser
    writefiles []*mockWriteCloser
    modtime time.Time
    timedfiles []string      
}

func (ms* mockSystem) OpenFileForReading(name string) (rd io.ReadCloser, err error) {
    p := &mockReadCloser{ 0, strings.NewReader(ms.input), name }
    ms.readfiles = append(ms.readfiles, p)
    return p, nil
}

func (ms* mockSystem) ModTime(name string) (modtime time.Time, err error) {
    ms.timedfiles = append(ms.timedfiles, name)
    return ms.modtime, nil
}

func (ms* mockSystem) OpenFileForWriting(name string) (wc io.WriteCloser, err error) {
    buffy := make([]byte, 0, 5000)
    mwc := &mockWriteCloser{0, bytes.NewBuffer(buffy), name}
    ms.writefiles = append(ms.writefiles, mwc)
    return mwc, nil
}

// Here is a classic example of something that could be
// polymorphic. Is there a way to do this?
func AssertInt(t *testing.T, expected int, actual int) {
    if (expected != actual) {
        t.Errorf("expected %s != actual %s", expected, actual)
    }
}

// TODO(rjkroege): Move testing utilities into a library of they're own.
func AssertString(t *testing.T, expected string, actual string) {
    if (expected != actual) {
        t.Errorf("expected %s != actual %s", expected, actual)
    }
}

const generated_output_2 = 
`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
   "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
<head>
   <title>What I want</title>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

  <!-- date argument for centering -->
  <script language="JavaScript" type="text/javascript">
    var external_titledate = 'Monday, Mar 19, 2012';
  </script>

  <!-- timeline CSS -->
  <link rel="stylesheet" href="timeline/resetfonts.css" type="text/css">
  <link rel="stylesheet" type="text/css" href="timeline/base.css">
  <link rel="stylesheet" href="timeline/timeline-bundle.css" type="text/css">

  <!-- timeline JavaScript -->
  <script src="timeline/simile-ajax-api.js" type="text/javascript"></script>
  <script src="timeline/timeline-bundle.js" type="text/javascript"></script>
  <script src="timeline/timeline.js" type="text/javascript"></script>

  <!-- timeline data -->
  <script src="note_list.js" type="text/javascript"></script>

  <!-- notes CSS -->   
  <link rel="stylesheet" type="text/css" media="all" href="notes.css" />
  <link rel="stylesheet" type="text/css" media="print" href="notes-print.css" />

  <!-- I don't use -->
  <script type="text/javascript" src="styleLineNumbers.js"></script>

</head>
<body onload="onLoad();" onresize="onResize();">
   <div id="container">
      <div id="title">
        <!-- Add editing functionality? -->
        <h1 class="left">What I want</h1>
        <h1 class="right">Monday, Mar 19, 2012</h1>
      </div> <!-- title -->
      <div id="note">
<p>I need to figure out what I want. </p>

<hr />
<p class="info">
   Source: <a href="plumb:hello">one.md</a><br />
   <a href="plumb:/Users/rjkroege/Dropbox/wiki2/New">New Article</a><br />
</p>
</div> <!-- note -->

      <!-- Timeline -->
      <div id="doc3" class="yui-t7">
        <div id="bd" role="main">
          <div class="yui-g">
            <div id='tl'></div>
          </div>
        </div>
      </div>

</div> <!-- container -->
</body>

</html>
`



func Test_WriteHtmlFile(t *testing.T) {
    realisticdate1999, _ := parseDateUnix("1999/03/21 17:00:00")
    realisticdate2012, _ := parseDateUnix("2012/03/19 06:51:15")
   
    // Produce output case.
    ms := &mockSystem { test_header_2,
            make([]*mockReadCloser, 0, 4),
            make([]*mockWriteCloser, 0, 4),
            time.Time{},
            make([]string, 0, 4)}

    md := MetaData{"one.md", "http://two", realisticdate1999, realisticdate2012, "What I want", true, "hello"};
    md.WriteHtmlFile(ms)

    AssertInt(t, 1, len(ms.writefiles))
    AssertInt(t, 1, len(ms.readfiles))
    AssertInt(t, 1, int(ms.writefiles[0].closedCount))
    AssertInt(t, 1, int(ms.readfiles[0].closedCount))

    AssertString(t, "one.md", ms.readfiles[0].name )
    AssertString(t, "one.html", ms.writefiles[0].name )
    AssertString(t, "one.html", ms.timedfiles[0])

    // TODO(rjkroege): might want to diff the stirngs?
    AssertString(t, generated_output_2, ms.writefiles[0].String())

    // Output production skipped by date comparison.
    ms = &mockSystem { test_header_2,
            make([]*mockReadCloser, 0, 4),
            make([]*mockWriteCloser, 0, 4),
            realisticdate2012.Time,
            make([]string, 0, 4)}

    md = MetaData{"one.md", "http://two", realisticdate1999, realisticdate2012, "What I want",  true, "hello"};
    md.WriteHtmlFile(ms)

    AssertInt(t, 0, len(ms.writefiles))
    AssertInt(t, 1, len(ms.readfiles))
    AssertInt(t, 1, int(ms.readfiles[0].closedCount))

    AssertString(t, "one.md", ms.readfiles[0].name )
    AssertString(t, "one.html", ms.timedfiles[0])

    // TODO(rjkroege): Add additional tests to support validating error handling, etc.
}

func Test_PrettyDate(t *testing.T) {
    statdate, _ := parseDateUnix("1999/03/21 17:00:00")
    tagdate, _ := parseDateUnix("2012/03/19 06:51:15")

    md := MetaData{"", "", statdate, never , "What I want 0", false, ""}
    AssertString(t, "Sunday, Mar 21, 1999", md.PrettyDate())

    md = MetaData{"", "", statdate, tagdate , "What I want 0", true, ""}
    AssertString(t, "Monday, Mar 19, 2012", md.PrettyDate())
}

type tEdMd struct {
    err error
    result string
    md MetaData
}


const json1 = `{"link":"url-here-0","start":"Sunday, Mar 21, 1999","title":"What I want 0"}`
const json2  = `{"link":"url-here-1","start":"Monday, Mar 19, 2012","title":"What I want 0"}`

func Test_JsonDate(t *testing.T) {
    statdate, _ := parseDateUnix("1999/03/21 17:00:00")
    tagdate, _ := parseDateUnix("2012/03/19 06:51:15")

    datas := []tEdMd {
        { nil, json1, MetaData{"", "url-here-0", statdate, never , "What I want 0", false, ""}},
        { nil, json2, MetaData{"", "url-here-1", statdate, tagdate , "What I want 0", true, ""}}}

    for _, m := range(datas) {
        b, e := m.md.MarshalJSON()
        if m.err != e {
            t.Errorf("error value wrong");
        }
        AssertString(t, m.result, string(b))
    }
}
