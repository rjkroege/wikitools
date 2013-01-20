package article

import (
    "bytes"
    "testing"
)

const expected =
`
var timeline_data = {  // save as a global variable
'dateTimeFormat': 'iso8601',
'wikiURL': "http://simile.mit.edu/shelf/",
'wikiSection': "Simile Cubism Timeline",

'events': 
[{"link":"url-here-a","start":"Monday, Mar 19, 2012","title":"What I want 0"},{"link":"url-here-b","start":"Sunday, Mar 21, 1999","title":"What I want 1"},{"link":"url-here-c","start":"Monday, Mar 19, 2012","title":"What I want 2"}]

 }
`

func Test_WriteTimeline(t *testing.T) {
    /* General idea: create a constant string. Read from it., validate the resulting output. */
    statdate, _ := parseDateUnix("1999/03/21 17:00:00")
    tagdate, _ := parseDateUnix("2012/03/19 06:51:15")

    metadatas := []*MetaData {
        &MetaData{"", "url-here-a", statdate, tagdate , "What I want 0", false, ""},
        &MetaData{"", "url-here-b", statdate, statdate , "What I want 1", false, ""},
        &MetaData{"", "url-here-c", statdate, tagdate , "What I want 2", false, ""}}

    buffy := make([]byte, 0, 5000)
    fd := bytes.NewBuffer(buffy)
    
    WriteTimeline(fd, metadatas)
    
    AssertString(t, expected, fd.String())

//    t.Errorf("output: {%s}", fd.String())    
}
