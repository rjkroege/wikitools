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
    .auto-column-list {
      column-width: 45ch;
      column-gap: 2rem;
    }
    .auto-column-list li {
      break-inside: avoid;
    }
  </style>
</head>
<body>

<ul class="auto-column-list">
{{range .}}
    <li>{{.RelPath}} ({{.ModTime.Format "2006-01-02 15:04:05"}})</li>
{{end}}
</ul>
  </ul>
</body>
</html>
`

func (tidy *listAllWikiFiles) _htmlSummaryWrite(w io.Writer) error {
	t := template.Must(template.New("list").Parse(listallwikitmpl))
	return t.Execute(w, tidy.Files)
}
