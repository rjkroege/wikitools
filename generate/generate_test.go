package generate

import (
	"bytes"
	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/testhelpers"
	"io"
	"strings"
	"testing"
	"time"
)

type closedCount int

func (cc *closedCount) Close() error {
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
	input      string
	readfiles  []*mockReadCloser
	writefiles []*mockWriteCloser
	modtime    time.Time
	timedfiles []string
}

func (ms *mockSystem) OpenFileForReading(name string) (rd io.ReadCloser, err error) {
	p := &mockReadCloser{0, strings.NewReader(ms.input), name}
	ms.readfiles = append(ms.readfiles, p)
	return p, nil
}

func (ms *mockSystem) ModTime(name string) (modtime time.Time, err error) {
	ms.timedfiles = append(ms.timedfiles, name)
	return ms.modtime, nil
}

func (ms *mockSystem) OpenFileForWriting(name string) (wc io.WriteCloser, err error) {
	buffy := make([]byte, 0, 5000)
	mwc := &mockWriteCloser{0, bytes.NewBuffer(buffy), name}
	ms.writefiles = append(ms.writefiles, mwc)
	return mwc, nil
}

const generated_output_2 = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
   "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
<head>
   <title>What I want</title>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
   <link rel="icon" href="icon.png">

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
   Source: <a href="plumb:/am-a-path/one.md">one.md</a><br />
   <a href="plumb:/Users/rjkroege/gda/wiki2/New">New Article</a><br />
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
const test_header_2 = `title: What I want
date: 2012/03/19 06:51:15
tags: @journal

I need to figure out what I want. 
`

func Test_WriteHtmlFile(t *testing.T) {
	realisticdate1999, _ := article.ParseDateUnix("1999/03/21 17:00:00")
	realisticdate2012, _ := article.ParseDateUnix("2012/03/19 06:51:15")
	article.SetPathForContent("/am-a-path")

	// Produce output case.
	ms := &mockSystem{test_header_2,
		make([]*mockReadCloser, 0, 4),
		make([]*mockWriteCloser, 0, 4),
		time.Time{},
		make([]string, 0, 4)}

	md := article.NewArticleWithTime("one.md", realisticdate1999, realisticdate2012, "What I want", true)
	WriteHtmlFile(ms, md)

	testhelpers.AssertInt(t, 1, len(ms.writefiles))
	testhelpers.AssertInt(t, 1, len(ms.readfiles))
	testhelpers.AssertInt(t, 1, int(ms.writefiles[0].closedCount))
	testhelpers.AssertInt(t, 1, int(ms.readfiles[0].closedCount))

	testhelpers.AssertString(t, "one.md", ms.readfiles[0].name)
	testhelpers.AssertString(t, "one.html", ms.writefiles[0].name)
	testhelpers.AssertString(t, "one.html", ms.timedfiles[0])

	// TODO(rjkroege): might want to diff the stirngs?
	testhelpers.AssertString(t, generated_output_2, ms.writefiles[0].String())

	// Output production skipped by date comparison.
	ms = &mockSystem{test_header_2,
		make([]*mockReadCloser, 0, 4),
		make([]*mockWriteCloser, 0, 4),
		realisticdate2012,
		make([]string, 0, 4)}

	md = article.NewArticleWithTime("one.md", realisticdate1999, realisticdate2012, "What I want", true)
	WriteHtmlFile(ms, md)

	testhelpers.AssertInt(t, 0, len(ms.writefiles))
	testhelpers.AssertInt(t, 1, len(ms.readfiles))
	testhelpers.AssertInt(t, 1, int(ms.readfiles[0].closedCount))

	testhelpers.AssertString(t, "one.md", ms.readfiles[0].name)
	testhelpers.AssertString(t, "one.html", ms.timedfiles[0])

	// TODO(rjkroege): Add additional tests to support validating error handling, etc.
}
