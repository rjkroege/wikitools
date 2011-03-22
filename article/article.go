/*
  Constructs the actual note. Does what the current python buildPage.py
  script does.
  
  insert here
  
  This file is probably mis-named.
  
*/

/*
  Must be a different package
  But now, I have a circular dependency. And how do I plan on handling this?
*/
package article;

import (
  "bufio"
  "fmt"
  "io"
  "os"
  "github.com/knieriem/markdown"
  "template"
)

// TODO(rjkroege): make sure that each entry has a nice comment
// and clean up.
type MetaData struct {
  Name string;    // Relative file name
  Url string;     // Url of the generated file.
  DateFromStat int64;
  DateFromMetadata int64;
  Title string;
  FinalDate string;
  hadMetaData bool;
  PrettyDate string;
  SourceUrl string;
}

// Contains large string contants.
const (
header =
`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
   "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
<head>
   <title>{Title}</title>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

  <!-- date argument for centering -->
  <script language="JavaScript" type="text/javascript">
    var external_titledate = '{PrettyDate}';
    function modifyTheUrl(event) {.meta-left}
      event.currentTarget.href += Date.now() + '.md';
      return true; {.meta-right}
  </script>

  <!-- timeline CSS -->
  <link rel="stylesheet" href="timeline/resetfonts.css" type="text/css">
  <link rel="stylesheet" type="text/css" href="timeline/base.css">
  <link rel="stylesheet" href="timeline/timeline-bundle.css" type="text/css">

  <!-- timeline JavaScript -->
  <script src="timeline/simile-ajax-api.js" type="text/javascript"></script>
  <script src="timeline/timeline-bundle.js" type="text/javascript"></script>
  <script src="timeline/timeline.js" type="text/javascript"></script>

  <!-- timeline data -->
  <script src="note_list.js" type="text/javascript"></script>

  <!-- notes CSS -->   
  <link rel="stylesheet" type="text/css" media="all" href="notes.css" />
  <link rel="stylesheet" type="text/css" media="print" href="notes-print.css" />

  <!-- I don't use -->
  <script type="text/javascript" src="styleLineNumbers.js"></script>

</head>
<body onload="onLoad();" onresize="onResize();">
   <div id="container">
      <div id="title">
        <!-- Add editing functionality? -->
        <h1 class="left">{Title}</h1>
        <h1 class="right">{PrettyDate}</h1>
      </div> <!-- title -->

      <!-- Timeline -->
      <div id="doc3" class="yui-t7">
        <div id="bd" role="main">
          <div class="yui-g">
            <div id='tl'></div>
          </div>
        </div>
      </div>

      
      <div id="note">
`;

footer =
`
<hr />
<p class="info">
   Source: <a href="plumb://open?url=file://$mdpath">$mdpath</a><br />
   Last modified: $modldate at $modtime<br />
   <a href="plumb://new?url=file:///Users/rjkroege/Documents/wiki2&template=file:///Users/rjkroege/Documents/wiki2/template.md">New Article</a><br />
   <!-- This page built: $buildtime -->
</p>
</div> <!-- note -->
</div> <!-- container -->
</body>
</html>
`

textmateFooter =
`
<hr />
<p class="info">
   Source: <a href="txmt://open?url={SourceUrl}">{Name}</a><br />
   <!-- new ones not handled -->
   <a onclick="modifyTheUrl(event)" href="txmt://open?url=file:///Users/rjkroege/Documents/wiki2/ca_">New Article</a><br />
</p>
</div> <!-- note -->
</div> <!-- container -->
</body>

</html>
`

// Used for exploring how the template facility works.
test = "foo foo {Title} bar bar\n{PrettyDate}\n{SourceUrl}\n{Name}\n"

)

var headerTemplate = template.MustParse(header, nil);
var footerTemplate = template.MustParse(textmateFooter, nil);

// Converts an article name into its name as a formatted object.s
func (md *MetaData) FormattedName() string {
  oname := md.Name[0:len(md.Name) - len(".md")] + ".html";
  return oname
}

// Constructs a URL path equivalent to the given source file.
func (md *MetaData) UrlForName(path string) string {
  // Prefix file:///<path>/fname.html
  md.Url = "file://" + path + "/" + md.FormattedName();
  return md.Url;
}

func (md *MetaData) SourceForName(path string) string {
  md.SourceUrl = "file://" + path + "/" + md.Name
  return md.SourceUrl
}

// TODO(rjkroege): it might be desirable to divide this funciton
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
  

  fd, err := os.Open(md.Name, os.O_RDONLY, 0)
  defer fd.Close()
  if err != nil {
    fmt.Print(err)
    return
  }
  
  statinfo, serr := os.Stat(md.FormattedName())
  if serr != nil || statinfo.Mtime_ns < md.DateFromStat {
    // TODO(rjkroege): if the md file is not as new as the HTML file, 
	  // skip all of this work.
    fmt.Println("processing " + md.Name)
    ofd, werr := os.Open(md.FormattedName(), os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644);
    defer ofd.Close()
    if werr != nil {
      fmt.Print("one ", werr, "\n");
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
          fmt.Print("two ", werr, "\n");
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
			if rerr == os.EOF {
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
		doc := markdown.Parse(body, markdown.Extensions{Smart: true})
    doc.WriteHtml(w)

    // Footer with substitutions
    footerTemplate.Execute(w, md)
    fmt.Println("done " + md.Name)
  }
}

