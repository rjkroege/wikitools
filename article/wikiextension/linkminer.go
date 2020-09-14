package wikiextension

import (
	"log"
	"path"
	"strings"

	"github.com/rjkroege/wikitools/wiki"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

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
	outboundwikilinks := make([][]byte, 0)

	// goldmark's Walk function starts at the given node and invokes the
	// provided function on every node. Nodes (well at least Link nodes)
	// invoke the function both on entering and leaving the node.
	gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if n.Kind() == gast.KindLink && !entering {
			log.Println("node:", n.Type(), "kind:", n.Kind(), "entering:", entering)
			log.Println("node text:", string(n.Text(reader.Source())))

			// Dump all link nodes.
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

			// The actual writing of the links to the KV-store can happen
			// asynchronously so I could send them off to a service Go routine here?
			// Maybe not.

			// Filter the link here
			if isWikiMarkdownLink(link.Destination) {
				outboundwikilinks = append(outboundwikilinks, link.Destination)
			}
		}
		return gast.WalkContinue, nil

	})

	// Dump everything (i.e. the root of the AST). For debugging.
	// node.Dump(reader.Source(), 0)

	// Dump the links.
	for _, d := range outboundwikilinks {
		log.Println("link", string(d))
	}

	// TODO(rjk): Insert link map instead
	log.Println("attempting to insert link map")
	a.insertLinkMap(node, reader, pc)

}

// isWikiLink returns true if the provided dest is a link inside of the
// wiki. Links are "inside" the wiki if they are relative or absolute
// with the root of the wiki as prefix.
// TODO(rjk): Should I make sure that there's a file at the end of the
// link? wikipp shoudln't but wikiclean should probably check link
// validity for all of the wiki articles and generate a report if they
// contain invalid links.
func isWikiLink(dest []byte) bool {
	pth := path.Clean(string(dest))
	if !path.IsAbs(pth) || strings.HasPrefix(pth, wiki.Basepath) {
		return true
	}
	return false
}

func isWikiMarkdownLink(dest []byte) bool {
	pth := path.Clean(string(dest))
	if path.Ext(pth) == wiki.Extension && (!path.IsAbs(pth) || strings.HasPrefix(pth, wiki.Basepath)) {
		return true
	}
	return false
}

// testInsertNode is an experiment to add stuff to the AST. Appends HelloWorld. This block of
// code worked as desired in that I got a h1 "Hello World" appended. This
// shows how I would add contents to the AST. Here would I generate the
// link map for this document and insert into place. Presume that the
// link map generator makes an SVG. Then, I should make an SVG here and insert.
func (a *linkMinerASTTransformation) testInsertNode(node *gast.Document, reader text.Reader, pc parser.Context) {
	nhs := gast.NewString([]byte("Hello World"))
	nh := gast.NewHeading(1)
	nh.AppendChild(nh, nhs)
	node.AppendChild(node, nh)

}

func (a *linkMinerASTTransformation) insertLinkMap(node *gast.Document, reader text.Reader, pc parser.Context) {
	// figure out how to make a simple svg object here
	svgs := gast.NewString(makeLinkMapTest())
	svgs.SetCode(true)

	nlm := gast.NewHTMLBlock(gast.HTMLBlockType7)
	nlm.AppendChild(nlm, svgs)
	node.AppendChild(node, nlm)
}

const simplesvg = `<svg height="100" width="100">
  <circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="red" />
  Sorry, your browser does not support inline SVG.  
</svg>`

// makeLinkMapTest is a stub implementation of generating a linkmap to
// show that the insertion is working. It inserts a red circle.
func makeLinkMapTest() []byte {
	return []byte(simplesvg)
}
