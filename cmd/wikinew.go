package cmd

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"9fans.net/go/acme"
	"github.com/rjkroege/wikitools/article"
	"github.com/rjkroege/wikitools/wiki"
)

func Makearticle(settings *wiki.Settings, args []string, tags []string) *article.MetaData {
	s := strings.Join(args, " ")
	vfn := wiki.ValidBaseName(args)
	return article.NewArticle(vfn, s, tags)
}

type ExpandedArticle struct {
	*article.MetaData
	buffy *bytes.Buffer
}

func (md *ExpandedArticle) Plumb(settings *wiki.Settings) {
	if err := os.MkdirAll(filepath.Join(settings.Wikidir, md.RelativeDateDirectory()), 0777); err != nil {
		log.Fatal(err)
	}

	win, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}

	reldir := md.RelativeDateDirectory()
	filename := md.PreferredFileName(settings)
	ufn := settings.UniquingExtension(reldir, filename)

	fullpathname := filepath.Join(settings.Wikidir, reldir, settings.ExtensionedFileName(filename + ufn))

	err = win.Name(fullpathname)
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

// TODO(rjkroege): return error state.
func Wikinew(settings *wiki.Settings, args []string) {
	tmpls := wiki.NewTemplatePalette()
	tmpls.AddDynamcTemplates(settings.TemplateForTag)

// TODO(rjk): refactor this.
	args, tags := wiki.Split(args)
	tm, args, tags := tmpls.Picktemplate(args, tags)
	// TODO(rjk): This is too cute. Don't do things like this.
	Expand(Makearticle(settings, args, tags), tm).Plumb(settings)
}
