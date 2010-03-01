/*
  Metadata extraction  
  ; fn gogo  {6g -I . listnotes.go  metadata.go generatemarkup.go && 6l  listnotes.6  && ./6.out ; echo}

*/

package main

import (
  "fmt";
  "io";
//  "io";
//  "strings";
//  "strconv";
//  "regexp";
//  "time";
  "exp/datafmt";
  "strings";
)


const (

// Insert at the top of each generated file.
  header =
`
var timeline_data = {  // save as a global variable
'dateTimeFormat': 'iso8601',
'wikiURL': "http://simile.mit.edu/shelf/",
'wikiSection': "Simile Cubism Timeline",

'events': [
`;

// Insert at the bottom of each generated file.
footer = 
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
main "./main";
string = "'%s'";
titleField = "'title': '%s'";
urlField = "'link': '%s'";
dateField = "'start': '%s'";
main.FileMetaData = "  {" ( "    " >>  "\n"
    Title:titleField ",\n"
    Url:urlField ",\n"
    FinalDate:dateField
    ) "\n  }";
ptr = * : main.FileMetaData;
array = { * / ",\n" };
`;

)


func writeMarkup(fd io.Writer, e []*FileMetaData) {

  // Might not need...
  // fmap := make(FormatterMap);
  

  df, err := datafmt.Parse("listnotes.go", strings.Bytes(emitter), nil);
  if err != nil {
    fmt.Print(err);
  } else {
    io.WriteString(fd, header);

    // mind that you have no looping (will add repetition
    // TODO(rjkroege): add repetition to the format.
    _, err2  :=df.Fprint(fd, nil, e);
    if (err2 != nil) {
      fmt.Print(err2);
      return;
    }
    io.WriteString(fd, footer);
  }
}


