/*
  Rips out the meta data from a single Journal file.
  {make &&  ./extractmeta *.md ; echo}
*/

package main

import (
//  "./article"
  "bytes";
  "exp/datafmt";
  "flag";
  "fmt";
  "os";
  "strings";
  "time";
)

const (
singleEmitter =
`
main "./main";
string = "'%s'";
titleField = "entrytitle = %s";
urlField = "link = %s";
dateField = "titledate = %s";
main.article.MetaData =
    Title:titleField "\n"
    Url:urlField "\n"
    FinalDate:dateField "\n";
ptr = * : main.article.MetaData;
`;

titleTimeFormat = "Mon Jan _2, 2006";
)



/**
 * Turns a time in ns since epoch into a string
 */
func dateToString(ti int64) string {
  t := time.SecondsToLocalTime(int64(ti / 1e9));
  // return t.Format(time.ISO8601);
  return t.Format(titleTimeFormat);
}


func main() {
  flag.Parse();
  pwd, _ := os.Getwd();
  
  df, err1 := datafmt.Parse("extractone.go",
      bytes.NewBufferString(singleEmitter).Bytes(), nil);
  if err1 != nil {
    fmt.Print(err1);
  }
  
  for i := 0; i < flag.NArg(); i++ {
    fname := flag.Arg(i);
  
    // Skip files of the wrong metadata
    fi, err := os.Stat(fname);
    if strings.HasSuffix(fname, ".md") && err == nil {
      e := new(article.MetaData);
      e.Name = fname;
      e.Url = makeUrlFromName(e.Name, pwd);
      
      e.DateFromStat = fi.Mtime_ns;
      e.DateFromMetadata, e.Title = rootThroughFileForMetadata(fname);

      if (e.DateFromMetadata > int64(0)) {
        e.FinalDate = dateToString(e.DateFromMetadata);
      } else {
        e.FinalDate = dateToString(e.DateFromStat);
      }

        
      
      _, err2 := df.Fprint(os.Stdout, nil, e);
      if err2 != nil {
        fmt.Print(err2);
      }

    } else if err != nil {
      fmt.Print(err);
    } else {
      fmt.Print("given file: " + fname + " is not a .md file\n");
    }
  }
}
