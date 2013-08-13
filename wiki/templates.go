package wiki

const (
journaltmpl = 
`title: {{.Title}}
date: {{.DetailedDate}}
tags: {{.Tagstring}}

Yo dawg! Write stuff here.
`

booktmpl =
`title: {{.Title}}
date: {{.DetailedDate}}
tags: {{.Tagstring}}
bib-address:
bib-author:
bib-booktitle: {{.Title}}
bib-chapter:
bib-edition:
bib-editor:
bib-eprint:
bib-howpublished:
bib-institution:
bib-journal:
bib-month:
bib-number:
bib-organization:
bib-pages:
bib-publisher:
bib-school:
bib-series:
bib-title: {{.Title}}
bib-type:
bib-url:
bib-volume:
bib-year:
bib-isbn:
bib-dewey:

Yo dawg! Put the bookreview here. Delete the undesired tags.
`

)
