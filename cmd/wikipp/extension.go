package main

import (
	"log"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// doesn't belong here
// an extension for goldmark to extract 

//    Why this structure? It keeps the internals private?
type linkminer struct {
}

var Linkminer = &linkminer{}

func (e *linkminer) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(NewLinkMinerASTTransformation(), 999),
		),
	)
}

type linkMinerASTTransformation struct {
}

// NewFootnoteASTTransformer returns a new parser.ASTTransformer that
// insert a footnote list to the last of the document.
func NewLinkMinerASTTransformation() parser.ASTTransformer {
	return &linkMinerASTTransformation{}
}

// called to do stuff to the AST
func (a *linkMinerASTTransformation) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if n.Kind() == gast.KindLink  && !entering {
			log.Println("node:", n.Type(), "kind:", n.Kind(), "entering:", entering)
			log.Println("node text:", string(n.Text(reader.Source())))

			// TODO(rjk): I'd like to be able to dump the node here.
			
		}
		return gast.WalkContinue, nil
	})
}
