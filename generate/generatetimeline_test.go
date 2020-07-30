package generate

import (
	"bytes"
	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/testhelpers"
	"testing"
)

const expected = `
var timeline_data = {  // save as a global variable
'dateTimeFormat': 'iso8601',
'wikiURL': "http://simile.mit.edu/shelf/",
'wikiSection': "Simile Cubism Timeline",

'events': 
[{"link":"file:///foo/bar0.html","start":"2012-03-19T06:51:15-04:00","title":"What I want 0"},{"link":"file:///foo/bar1.html","start":"1999-03-21T17:00:00-05:00","title":"What I want 1"},{"link":"file:///foo/bar2.html","start":"2012-03-19T06:51:15-04:00","title":"What I want 2"}]

 }
`

func Test_WriteTimeline(t *testing.T) {
	/* General idea: create a constant string. Read from it., validate the resulting output. */
	statdate, _ := article.ParseDateUnix("1999/03/21 17:00:00")
	tagdate, _ := article.ParseDateUnix("2012/03/19 06:51:15")
	article.SetPathForContent("/foo")

	metadatas := []*article.MetaData{
		article.NewArticleWithTime("bar0.md", statdate, tagdate, "What I want 0", article.MdInvalid),
		article.NewArticleWithTime("bar1.md", statdate, statdate, "What I want 1", article.MdInvalid),
		article.NewArticleWithTime("bar2.md", statdate, tagdate, "What I want 2", article.MdInvalid)}

	buffy := make([]byte, 0, 5000)
	fd := bytes.NewBuffer(buffy)

	WriteTimeline(fd, metadatas)

	testhelpers.AssertString(t, expected, fd.String())

	//    t.Errorf("output: {%s}", fd.String())
}
