/*
  Creates the contents of the JavaScript file that lists the entries in the
  timeline.
*/

package article

import (
  "io";
  "text/template";
)


const (

// Insert at the top of each generated file.
  timeline_header =
`
var timeline_data = {  // save as a global variable
'dateTimeFormat': 'iso8601',
'wikiURL': "http://simile.mit.edu/shelf/",
'wikiSection': "Simile Cubism Timeline",

'events': [
`;

// Insert at the bottom of each generated file.
timeline_footer = 
`
] }
`;


// TODO(rjkroege): understand the way of the formatter: indent.
// I have this feeling that this interface can do more but I'm missing
// the boat somehow.

// Use to actually generate the output using the formatter.
// 


emittercore =
`{{ define "core" }}   'title': '{{.Title}}',
    'link': '{{.Url}}',
    'start': '{{.FinalDate}}'{{end}}
`

// There should be a nicer way to do this.
emitter =
`  {
 {{template "core" .}}
  },
`

lastemitter =
`  {
 {{template "core" .}}
  }`
)

// It feels like it there are nicer ways to write this.
func WriteTimeline(fd io.Writer, e []*MetaData) {
  io.WriteString(fd, timeline_header);

  te, err := template.New("lastemitter").Parse(lastemitter);
  if err != nil { panic(err) }
  te, err = te.Parse(emittercore);
  if err != nil { panic(err) }

  t, err := template.New("emitter").Parse(emitter);
  if err != nil { panic(err) }
  t, err = t.Parse(emittercore);
  if err != nil { panic(err) }

  for i := 0; i < len(e) - 1; i++ {
    err = t.Execute(fd, e[i])
    if err != nil { 
      panic(err)
    }
  }
  if len(e) > 0 {  
    err = te.Execute(fd, e[len(e) - 1])
    if err != nil { 
      panic(err)
    }
  }

  io.WriteString(fd, timeline_footer);
}


