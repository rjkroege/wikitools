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
	HadMetaData      bool
	Tags             []string
	ExtraKeys        map[string]string
	Datepath string
}

const MDFOtemplate = `MetaData:
	Name: {{.Name}}
	DateFromStat:   {{.DateFromStat}}
	DateFromMetadata: {{.DateFromMetadata}}
	Title:          {{.Title}}
	Dynamicstring:    {{.Dynamicstring}}
	HadMetaData:      {{.HadMetaData}}
	Tags:             {{.Tags}}
	ExtraKeys:        {{.ExtraKeys}}
	Datepath: {{.Datepath}}
`

func (md *MetaData) Dump() string {
	tmpl, err := template.New("dumper").Parse(MDFOtemplate)
	if err != nil {
		panic("oops!, bad template")
	}

	yada := &MetaDataForOutput{
		Name:             md.Name,
		DateFromStat:     md.DateFromStat,
		DateFromMetadata: md.DateFromMetadata,
		Title:            md.Title,
		Dynamicstring:    md.Dynamicstring,
		HadMetaData:      md.HadMetaData,
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
