package article

import (
    "bytes"
    "testing"
    "time"
)

const expected =
`
var timeline_data = {  // save as a global variable
'dateTimeFormat': 'iso8601',
'wikiURL': "http://simile.mit.edu/shelf/",
'wikiSection': "Simile Cubism Timeline",

'events': 
[{"link":"url-here-a","title":"What I want 0","start":"start-a"},{"link":"url-here-b","title":"What I want 1","start":"start-b"},{"link":"url-here-c","title":"What I want 2","start":"start-cc"}]

 }
`

func Test_WriteTimeline(t *testing.T) {
    /* General idea: create a constant string. Read from it., validate the resulting output. */

    notime := time.Time{}
    metadatas := []*MetaData {
        &MetaData{"", "url-here-a", notime, notime , "What I want 0", "start-a", false, ""},
        &MetaData{"", "url-here-b", notime, notime , "What I want 1", "start-b", false, ""},
        &MetaData{"", "url-here-c", notime, notime , "What I want 2", "start-cc", false, ""}}

    buffy := make([]byte, 0, 5000)
    fd := bytes.NewBuffer(buffy)
    
    WriteTimeline(fd, metadatas)
    
    AssertString(t, expected, fd.String())

//    t.Errorf("output: {%s}", fd.String())    
}
