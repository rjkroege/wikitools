package article

import (
	"log"
	"strings"
	"text/template"
	"time"
)

type MetaDataForOutput struct {
	Name             string
	DateFromStat     time.Time
	DateFromMetadata time.Time
	Title            string
	Dynamicstring    string
	Mdtype      string
	Tags             []string
	ExtraKeys        map[string]string
	Datepath         string
}

const MDFOtemplate = `MetaData:
	Name: {{.Name}}
	DateFromStat:   {{.DateFromStat}}
	DateFromMetadata: {{.DateFromMetadata}}
	Title:          {{.Title}}
	Dynamicstring:    {{.Dynamicstring}}
	mdtype:      {{.Mdtype}}
	Tags:             {{.Tags}}
	ExtraKeys:        {{.ExtraKeys}}
	Datepath: {{.Datepath}}
`

var MdtypeNames = [...]string{
"MdInvalid ",
"MdLegacy",
"MdIaWriter",
"MdModern",
  }


func (md *MetaData) Dump() string {
	tmpl, err := template.New("dumper").Parse(MDFOtemplate)
	if err != nil {
		panic("oops!, bad template")
	}

	// TODO(rjk): Update for metadata type categorization
	yada := &MetaDataForOutput{
		Name:             md.filename,
		DateFromStat:     md.DateFromStat,
		DateFromMetadata: md.DateFromMetadata,
		Title:            md.Title,
		Dynamicstring:    md.Dynamicstring,
		Mdtype:      MdtypeNames[md.mdtype], 
		Tags:             md.tags,
		ExtraKeys:        md.extraKeys,
		Datepath:         md.datepath,
	}

	b := new(strings.Builder)
	if err := tmpl.Execute(b, yada); err != nil {
		log.Fatal("oops, bad template write because", err)
	}
	return b.String()
}
