/*
    The code needed to actually construct an article's HTML representation.
    Moved here as a prelude to separating generation into a different package.
*/

package article;

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "github.com/knieriem/markdown"
  "text/template"
  "strings"
)


var headerTemplate = template.Must(template.New("header").Parse(header));
var footerTemplate = template.Must(template.New("footer").Parse(plumberfooter));

// TODO(rjkroege): it might be desirable to divide this function
// up.
/**
  Given a article.MetaData object containing some paths and stuff,
  does appropriate transformations to construct the HTML form.
*/
func (md *MetaData) WriteHtmlFile() {
  // TODO(rjkroege): it is silly to re-open these files when I 
  // have had them open before them. And to re-read chunks of
  // them when I have already done so. But this is easier. And
	// it probably doesn't matter given that most files don't need
	// to be regenerated.

  // Handle  failures in the library.
  defer func() {
   if r := recover(); r != nil {
      fmt.Println(md.Name, "WriteHtmlFile, failed", r)
      }
    }()  

  fd, err := os.OpenFile(md.Name, os.O_RDONLY, 0)
  defer fd.Close()
  if err != nil {
    fmt.Print(err)
    return
  }
  
  statinfo, serr := os.Stat(md.FormattedName())
   // This might be suspect?
  if serr != nil || statinfo.ModTime().Before(md.DateFromStat) {
    // TODO(rjkroege): if the md file is not as new as the HTML file, 
    // skip all of this work.
    // fmt.Println("processing " + md.Name)
    ofd, werr := os.OpenFile(md.FormattedName(), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644);
    defer ofd.Close()
    if werr != nil {
      // fmt.Print("one ", werr, "\n");
      return
    }

    // TODO(rjkroege): using a byte slice might be faster?
    // learn to improve perf on Go.
    body := "";
    rd := bufio.NewReader(io.Reader(fd));

    // Trim the metadata here.
    if md.hadMetaData {
      for {
        line, rerr := rd.ReadString('\n');
        if rerr != nil {
          // fmt.Print("two ", werr, "\n");
          return
        }
        if line == "\n" {
          break
        }
      }
    }
    
    // TODO(rjkroege): don't read the file into memory.
    // Read errors will wipe out previously generated output. Do I care?
    for {
      line, rerr := rd.ReadString('\n');
			if rerr == io.EOF {
        break
			} else if rerr != nil {
			  fmt.Print("WriteHtmlFile: read error ", rerr, "\n")
			  return
			}
      body += line;
    }

    w := bufio.NewWriter(ofd)
    defer w.Flush()

    // Header with substitutions
    headerTemplate.Execute(w, md)

    // Convert the markdown file into a HTML
    p := markdown.NewParser(&markdown.Extensions{Smart: true});
    p.Markdown(strings.NewReader(body), markdown.ToHTML(w));

    // Footer with substitutions
    footerTemplate.Execute(w, md)
    // fmt.Println("done " + md.Name)
  }
}
