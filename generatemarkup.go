/*
  Metadata extraction  
  ; fn gogo  {6g -I . listnotes.go  metadata.go generatemarkup.go && 6l  listnotes.6  && ./6.out ; echo}

*/

package main

import (
  "fmt";
  "io";
//  "bufio";
//  "io";
//  "strings";
//  "strconv";
//  "regexp";
//  "time";
)

func writeMarkup(fd io.Writer, e []*FileMetaData) {
  fmt.Println("hello");
}
