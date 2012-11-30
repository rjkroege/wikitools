/*
  Creates the contents of the JavaScript file that lists the entries in the
  timeline.
*/

package article

import (
  "io"
  "encoding/json"
)

const (
// Insert at the top of each generated file.
  timeline_header =
`
var timeline_data = {  // save as a global variable
'dateTimeFormat': 'iso8601',
'wikiURL': "http://simile.mit.edu/shelf/",
'wikiSection': "Simile Cubism Timeline",

'events': 
`

// Insert at the bottom of each generated file.
timeline_footer = 
`
 }
`)

func WriteTimeline(fd io.Writer, e []*MetaData) {
  io.WriteString(fd, timeline_header);

  enc := json.NewEncoder(fd);
  err := enc.Encode(e); 
  if err != nil { panic(err) }

  io.WriteString(fd, timeline_footer);
}
