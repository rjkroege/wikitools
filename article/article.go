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
)


type MetaData struct {
  Name string;
  Url string;
  DateFromStat int64;
  DateFromMetadata int64;
  Title string;
  FinalDate string;
}

// Contains large string contants.
const (
header =
`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
   "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
<head>
   <title>$entrytitle</title>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

  <!-- date argument for centering -->
  <script language="JavaScript" type="text/javascript">
    var external_titledate = '$titledate';
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
        <h1 class="left">$entrytitle</h1>
        <h1 class="right">$titledate</h1>
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
)




/*
  How to proceed?
  
  We need some testing infrastructure. And other good stuff.
  
  I want to have less cruft in the wiki directory proper. It should be
  clean. Which behooves placing the template data inline here (in the
  go source).
  
  This is somewhat offensive but because compiles are fast, I can
  largely deal. It is however (in the long term) not the right approach.
  A better way might be to store them as html for proper editing and
  then suck them into Go code during the build process.

*/


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



// Idea: it is desirable to not have enormous side-effect
// intense functions such as this one.

/**
  Given a article.MetaData object containing some paths and stuff,
  does appropriate transformations to construct the HTML form.
*/
func (md *MetaData) Build() {
  // TODO(rjkroege): it is silly to re-open these files when I 
  // have had them open before them. And to re-read chunks of
  // them when I have already done so. But this is easier. And
	// it probably doesn't matter given that most files don't need
	// to be regenerated.
  fd, err := os.Open(md.Name, os.O_RDONLY, 0)
  
  if err != nil {
    fmt.Print(err)
    return
  }

  // TODO(rjkroege): if the md file is not as new as the HTML file, 
	// skip all of this work.
  ofd, werr := os.Open(md.FormattedName(), os.O_WRONLY | os.O_CREATE, 0644);
  
  if werr != nil {
    fmt.Print(werr);
    return
  } else {
    body := "";
    rd := bufio.NewReader(io.Reader(fd));
  
    for {
      line, rerr := rd.ReadString('\n');
      body += line;
			if rerr != nil {
				break
			}
    }
  
		// so... low road is to run a command here. We have still removed
		// some forks.
		fmt.Println("processing a file...")

		// Convert the markdown file into a HTML
		doc := markdown.Parse(body, markdown.Extensions{Smart: true})
    ofd.WriteString(header)
    w := bufio.NewWriter(ofd)
    doc.WriteHtml(w)
    w.Flush()
    ofd.WriteString(footer)

    // TODO(rjkroege):     
    // 4. replace special symbols with some properties.
    // setup the stuff that we are inserting.
    // result = result modified...
    // ofd.WriteString(result);
  }
}

