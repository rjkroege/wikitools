/*
  Creates the contents of the JavaScript file that lists the entries in the
  timeline.
*/

package article

import (
  "fmt";
  "bytes";
  "io";
//  "io";
//  "strings";
//  "strconv";
//  "regexp";
//  "time";
  "exp/datafmt";
  "go/token"
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
emitter =
`
article "liqui.org/article";
main "./main";
string = "'%s'";
titleField = "'title': '%s'";
urlField = "'link': '%s'";
dateField = "'start': '%s'";
article.MetaData = "  {" ( "    " >>  "\n"
    Title:titleField ",\n"
    Url:urlField ",\n"
    FinalDate:dateField
    ) "\n  }";
ptr = * : article.MetaData;
array = { * / ",\n" };
`;

)


func WriteTimeline(fd io.Writer, e []*MetaData) {
  io.WriteString(fd, timeline_header);
  
  // it is highly unclear to me how to do this
  df, err := datafmt.Parse(token.NewFileSet(), "article.go",
      bytes.NewBufferString(emitter).Bytes(), nil);
  if err != nil {
		fmt.Print("Something went wrong with the formatted output: ")
    fmt.Println(err);
  } else {
    // mind that you have no looping (will add repetition
    // TODO(rjkroege): add repetition to the format.
    _, err2  := df.Fprint(fd, nil, e);
    fmt.Println("supposedly, I have generated output here")
    if (err2 != nil) {
      fmt.Print(err2);
      return;
    }
  }
  io.WriteString(fd, timeline_footer);
}


