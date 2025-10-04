package corpus

import (
	"html/template"
	"io"
)

const listallwikitmpl = `
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Wiki Article List</title>
 <style>
    /* --- Container --- */
    .list-wrapper {
      width:  90%;             /* whatever width you need */
      margin: 2rem auto;    /* center the block */
      border: 1px solid #ccc;
      padding: 1rem;
    }

    /* --- List reset --- */
    ul.fill-across {
      list-style: none;
      margin: 0;
      padding: 0;
      display: flex;        /* put items in a row */
      flex-wrap: wrap;        /* allow wrapping to next line */
      gap: 1rem;              /* space between items */
    }

    /* --- List items --- */
    ul.fill-across li {
      flex: 1 1 200px;        /* grow, shrink, base width 200px */
     padding: 1rem;
    }
  </style>
</head>
<body>
<h1>Article List</h1>
 <div class="list-wrapper">
    <ul class="fill-across">
{{range .}}
    <li>{{.RelPath}} ({{.ModTime.Format "2006-01-02 15:04:05"}})</li>
{{end}}
    </ul>
  </div>
</body>
</html>
`

func (tidy *listAllWikiFiles) _htmlSummaryWrite(w io.Writer) error {
	t := template.Must(template.New("list").Parse(listallwikitmpl))
	return t.Execute(w, tidy.Files)
}
