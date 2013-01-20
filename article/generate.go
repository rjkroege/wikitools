/*
    The code needed to actually construct an article's HTML representation.
    Moved here as a prelude to separating generation into a different package.
*/

package article;

import (
   "time"
  "bufio"
  "fmt"
  "github.com/knieriem/markdown"
  "io"
  "os"
  "strings"
  "text/template"
)

/*
    Dependency injection interfaces. (Could move this somewhere else?)
    How to write this with the least pain?

    Editorial: writing things in terms of interfaces makes it easy to
    do dependency injection for testing or re-purposing. The way of Go
    hacking is clearly to code stupid, then up-level against an
    interface and make the interface configurable.

    Or something like that.
*/
type System interface {
    OpenFileForReading(name string) (rd io.ReadCloser, err error)
    ModTime(name string) (modtime time.Time, err error)
    OpenFileForWriting(name string) (wr io.WriteCloser, err error)
}

func (md* MetaData) OpenFileForReading(name string) (rd io.ReadCloser, err error) {
    rd, err = os.OpenFile(name, os.O_RDONLY, 0)
    return
}

func (md* MetaData) ModTime(name string) (modtime time.Time, err error) {
    statinfo, err := os.Stat(name)
    if err != nil {
         modtime = statinfo.ModTime()
    }
    return
}

func (md* MetaData) OpenFileForWriting(name string) (wr io.WriteCloser, err error) {
    wr, err = os.OpenFile(name, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
    return
}

var headerTemplate = template.Must(template.New("header").Parse(header));
var footerTemplate = template.Must(template.New("footer").Parse(plumberfooter));

// TODO(rjkroege): it might be desirable to divide this function
// up.
/**
  Given a article.MetaData object containing some paths and stuff,
  does appropriate transformations to construct the HTML form.
*/
func (md *MetaData) WriteHtmlFile(sys System) {
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

  fd, err := sys.OpenFileForReading(md.Name)
  defer fd.Close()
  if err != nil {
    fmt.Print(err)
    return
  }
  
  modtime, serr := sys.ModTime(md.FormattedName())
   // This might be suspect?
  if serr != nil || modtime.Before(md.DateFromStat.Time) {
    // TODO(rjkroege): if the md file is not as new as the HTML file, 
    // skip all of this work.
    // fmt.Println("processing " + md.Name)
    ofd, werr := sys.OpenFileForWriting(md.FormattedName());
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
