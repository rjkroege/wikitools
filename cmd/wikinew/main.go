package main

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	// "fmt"			// needeebugging.

	"9fans.net/go/acme"
	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/rjkroege/wikitools/config"
)

func Makearticle(args []string, tags []string) *article.MetaData {
	s := strings.Join(args, " ")
	
	tmpmd := article.NewArticle("", "", []string{})
	destpath := filepath.Join(config.Basepath, tmpmd.RelativeDateDirectory())
	filename := wiki.UniqueValidName(destpath, wiki.ValidBaseName(args), config.Extension, wiki.SystemImpl(0))
	return article.NewArticle(filename, s, tags)
}

type ExpandedArticle struct {
	*article.MetaData
	buffy *bytes.Buffer
}

func (md *ExpandedArticle) Plumb() {
	if err := os.MkdirAll(filepath.Join(config.Basepath, md.RelativeDateDirectory()), 0777); err != nil {
		log.Fatal(err)
	}

	win, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}

	err = win.Name(md.FullPathName(config.Basepath))
	if err != nil {
		log.Fatal(err)
	}

	_, err = win.Write("body", md.buffy.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	err = win.Fprintf("tag", "mdpreview tableflip")
	if err != nil {
		log.Fatal(err)
	}
}

/*
func (md* Article) Tagstring() string {
    return strings.Join(md.Tags,  " ")
}
*/

func Expand(a *article.MetaData, tpl wiki.Template) *ExpandedArticle {
	f := template.Must(template.New("newwiki").Parse(tpl.Template))
	a.Dynamicstring = tpl.Custombody

	b := new(bytes.Buffer)
	f.Execute(b, a)
	return &ExpandedArticle{a, b}
}

// TODO(rjkroege): add usage output on failure.
func main() {
	config := wiki.ReadConfiguration()
	tmpls := wiki.NewTemplatePalette()
	tmpls.AddDynamcTemplates(config)

	args, tags := wiki.Split(os.Args[1:])
	tm, args, tags := tmpls.Picktemplate(args, tags)
	Expand(Makearticle(args, tags), tm).Plumb()
}
