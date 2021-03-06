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
)

func Makearticle(args []string, tags []string) *article.MetaData {
	s := strings.Join(args, " ")

	tmpmd := article.NewArticle("", "", []string{})
	destpath := filepath.Join(wiki.Basepath, tmpmd.RelativeDateDirectory())
	filename := wiki.UniqueValidName(destpath, wiki.ValidBaseName(args), wiki.Extension, wiki.SystemImpl(0))
	return article.NewArticle(filename, s, tags)
}

type ExpandedArticle struct {
	*article.MetaData
	buffy *bytes.Buffer
}

func (md *ExpandedArticle) Plumb() {
	if err := os.MkdirAll(filepath.Join(wiki.Basepath, md.RelativeDateDirectory()), 0777); err != nil {
		log.Fatal(err)
	}

	win, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}

	err = win.Name(md.FullPathName(wiki.Basepath))
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

func Expand(a *article.MetaData, tpl wiki.Template) *ExpandedArticle {
	// TODO(rjk): Better error handling here. The templates can come from
	// user data. I need some kind of validation.
	f := template.Must(template.New("newwiki").Parse(tpl.Template))
	a.Dynamicstring = tpl.Custombody

	b := new(bytes.Buffer)
	f.Execute(b, a)
	return &ExpandedArticle{a, b}
}

// TODO(rjkroege): add usage output on failure.
// TODO(rjkroege): support editors other than Acme/Edwood.
func main() {
	config := wiki.ReadConfiguration()
	tmpls := wiki.NewTemplatePalette()
	tmpls.AddDynamcTemplates(config)

	args, tags := wiki.Split(os.Args[1:])
	tm, args, tags := tmpls.Picktemplate(args, tags)
	// TODO(rjk): This is too cute. Don't do things like this.
	Expand(Makearticle(args, tags), tm).Plumb()
}
