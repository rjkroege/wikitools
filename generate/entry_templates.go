package generate;

const (
header =
`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
   "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html>
<head>
   <title>{{.Title}}</title>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />

  <!-- date argument for centering -->
  <script language="JavaScript" type="text/javascript">
    var external_titledate = '{{.PrettyDate}}';
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
        <h1 class="left">{{.Title}}</h1>
        <h1 class="right">{{.PrettyDate}}</h1>
      </div> <!-- title -->
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

plumberfooter =
`
<hr />
<p class="info">
   Source: <a href="plumb:{{.SourceForPath}}">{{.Name}}</a><br />
   <a href="plumb:/Users/rjkroege/gd/wiki2/New">New Article</a><br />
</p>
</div> <!-- note -->

      <!-- Timeline -->
      <div id="doc3" class="yui-t7">
        <div id="bd" role="main">
          <div class="yui-g">
            <div id='tl'></div>
          </div>
        </div>
      </div>

</div> <!-- container -->
</body>

</html>
`

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
`

)