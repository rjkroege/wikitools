// TODO(rjk): It's arguable that this might be the wrong class name.
package wikiextension

import (
	"log"

	"github.com/rjkroege/wikitools/corpus"
	"github.com/rjkroege/wikitools/wiki"
	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/wikilink"
)

// I followed along with the goldmark code with this pattern but I'm not
// sure why it should be arranged that way.
// NB: there is a separate Linkminer instance per file.
type Linkminer struct {
	settings *wiki.Settings
	recorder corpus.UrlRecorder
	fpath    string
}

func NewLinkminer(settings *wiki.Settings, recorder corpus.UrlRecorder, fpath string) *Linkminer {
	return &Linkminer{
		settings: settings,
		recorder: recorder,
		fpath:    fpath,
	}
}

// Show that Linkminer is a goldmark.Extender
var _ goldmark.Extender = (*Linkminer)(nil)

func (e *Linkminer) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			// I don't understand how to correctly set priorties in goldmark's extensions.
			util.Prioritized(NewLinkMinerASTTransformation(e.settings, e.recorder, e.fpath), 999),
		),
	)
}

// It's arguable that I could have a single struct that implements all the stuff?
// And it would end up being simpler code?
// TODO(rjk): Consider opportunities for code simplification
type linkMinerASTTransformation struct {
	settings *wiki.Settings
	recorder corpus.UrlRecorder
	fpath    string
}

// Show that linkMinerASTTransformation is a goldmark.ASTTransformer
var _ parser.ASTTransformer = (*linkMinerASTTransformation)(nil)

// NewLinkMinerASTTransformation returns a new parser.ASTTransformer that
// can extract all of the Links found in a document.
func NewLinkMinerASTTransformation(settings *wiki.Settings, recorder corpus.UrlRecorder, fpath string) parser.ASTTransformer {
	return &linkMinerASTTransformation{
		settings: settings,
		recorder: recorder,
		fpath:    fpath,
	}
}

// Transform is an example of one way to extend goldmark.
func (a *linkMinerASTTransformation) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	// goldmark's Walk function starts at the given node and invokes the
	// provided function on every node. Nodes (well at least Link nodes)
	// invoke the function both on entering and leaving the node.
	gast.Walk(node, func(n gast.Node, entering bool) (gast.WalkStatus, error) {
		if n.Kind() == wikilink.Kind && !entering {
			link, ok := n.(*wikilink.Node)
			if !ok {
				// Conceivably this is an internal error in Goldmark because I should
				// have made sure that n is a WikiLink above.
				log.Println("WikiLink node not the right type!!!!")
				return gast.WalkContinue, nil
			}

			// WikiLinks are perhaps more complicated. They can also have a location.
			// Some parsing extensions might be needed.
			log.Printf("WikiLink Node %q %q %q", string(link.Target), string(link.Fragment), a.fpath)

			// TODO(rjk): Need to mark this kind of link differently so that they can be formatted
			// correctly.
			a.recorder.RecordWikilink(string(link.Fragment), string(link.Target), a.fpath)
		}

		if n.Kind() == gast.KindLink && !entering {
			log.Println(a.fpath, "node:", n.Type(), "kind:", n.Kind(), "entering:", entering)
			log.Println("node text:", string(n.Text(reader.Source())))

			// Dump all link nodes.
			n.Dump(reader.Source(), 1)

			link, ok := n.(*gast.Link)
			if !ok {
				// Conceivably this is an internal error in Goldmark because I should
				// have made sure that n is a Link above.
				return gast.WalkContinue, nil
			}

			// dest is the foofoo in [blah](foofoo).
			dest := string(link.Destination)

			// title is the blah in [blah](foofoo).
			title := string(link.Title)
			if title == "" {
				title = string(n.Text(reader.Source()))
			}
			a.recorder.RecordUrl(title, dest, a.fpath)

			// I can use a.settings.IsWikiMarkdownLink() to determine if a
			// URL points into the wiki. It is possible that I do not need this feature.
		}
		return gast.WalkContinue, nil
	})
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
