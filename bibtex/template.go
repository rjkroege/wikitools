package bibtex

const (

bibtextmpl =
`@{{.EntryType}}({{.RefKey}},
{{range $fieldname, $fieldvalue := .Fields}}	{{$fieldname}} = "{{$fieldvalue}}",
{{end}})
`
)
