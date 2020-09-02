package main

import (
	"log"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// TODO(rjk): Move this functionality out of the command.

// I followed along with the goldmark code with this pattern but I'm not
// sure why it should be arranged that way.
type linkminer struct {
}

var Linkminer = &linkminer{}

func (e *linkminer) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			// I don't understand how to correctly set priorties in goldmark's extensions.
			util.Prioritized(NewLinkMinerASTTransformation(), 999),
		),
	)
}

type linkMinerASTTransformation struct {
}

// NewLinkMinerASTTransformation returns a new parser.ASTTransformer that
// can extract all of the Links found in a document.
func NewLinkMinerASTTransformation() parser.ASTTransformer {
	return &linkMinerASTTransformation{}
}

// Transform is an example of one way to extend goldmark.
func (a *linkMinerASTTransformation) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {

	// goldmark's Walk function starts at the given node and invokes the
	// provided function on every node. Nodes (well at least Link nodes)
	// invoke the function both on entering and leaving the node.
	gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if n.Kind() == gast.KindLink && !entering {
			log.Println("node:", n.Type(), "kind:", n.Kind(), "entering:", entering)
			log.Println("node text:", string(n.Text(reader.Source())))

			// TODO(rjk): I'd like to be able to dump the node here.
			n.Dump(reader.Source(), 1)

			link, ok := n.(*gast.Link)
			if !ok {
				// Conceivably this is an internal error in Goldmark
				return gast.WalkContinue, nil
			}

			log.Println("Link des:", string(link.Destination))
			// link.Destination here is the actual link (URL). In the case of
			// intra-wiki links, these are the []byte that I want to store and use fo
			// generating link maps, etc.

		}
		return gast.WalkContinue, nil

	})

	// TODO(rjk): I would generate linkmaps here for insertion into the document.
}
