package tidy

import (
	"fmt"
	"html/template"
	"io"

	"log"
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
    <li>{{ filetourl .Path }} ({{.ModTime.Format "2006-01-02"}})</li>
{{end}}
    </ul>
  </div>
</body>
</html>
`

func (tidy *listAllWikiFiles) _htmlSummaryWrite(w io.Writer) error {
	// TODO(rjk): Stash all the templates in a central place where
	// the HTML can be refactored for rapid development.
	// Note that I need to mark filetourl as being safe to include.

	// Temporary logging to demonstrate that feature is correct.
	log.Println("_htmlSummaryWrite", tidy.tags)

	if tidy.tmpl == nil {
		tmpl, err := template.New("articlelist").Funcs(template.FuncMap{
			"filetourl": func(path string) template.HTML {
				return template.HTML(filetourl(tidy.settings.Wikidir, path))
			},
		}).Parse(listallwikitmpl)
		if err != nil {
			return fmt.Errorf("can't listallwikitmpl template%v", err)
		}
		tidy.tmpl = tmpl
	}
	return tidy.tmpl.ExecuteTemplate(w, "articlelist", tidy.Files)
}
