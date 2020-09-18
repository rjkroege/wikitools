package tidy

import (
	"github.com/rjkroege/wikitools/corpus"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark"
	"github.com/rjkroege/wikitools/article"
)
	

// find all of the outlinks from all of the articles
// use this as the context? A Context has a lot of methods. I'd rather not.
type findOutLinks struct {
	md goldmark.Markdown
	context parser.Context
	linkdbkey parser.ContextKey
}

// NewFindOutLinks creates a new Tidying implementation to update
// metadata.
// We need a bunch a bunch of state here.
func NewFindOutLinks() (corpus.Tidying, error) {
	// todo create the database if it doesn't exist
	// upgrade the database format if it exists and is out of date with the binary
	// fail if the database is newer than I am
	return newFindOutLinksImpl()
}

func newFindOutLinksImpl() (*metadataUpdater, error) {
	
	linkdbkey := parser.NewContextKey()
	context := parser.NewContext()

	// TODO(rjk): make a linkdb here.
	context.Set(linkdbkey, linkdb)

	return &findOutLinks{
		md: article.NewDefaultMarkdownConverter(),
		linkdbkey: linkdbkey,
		context: context,
	}
}

func (f *findOutLinks) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}

	if wiki.IsWikiArticle(path, info) {
		return nil
	}

	// I haz code for this mostly now
	updatedpth, err := f.findAllLinks(path)
	if err != nil {
		return err
	}

	return nil
}

func (f *findOutLinks) findAllLinks(path string) (string, error) {
	ifd, err := os.Open(path)
	if err != nil {
		log.Println("updateMetadata Open error", err)
		return "", fmt.Errorf("can't DoMetadataUpdate Open %s: %v", path, err)
	}

	// TODO(rjk): read the metadata
	// skip article if it's not a supported kind

	// We are not interested in writing out state. figure out how to skip that.

	return nil
}

func (f *findOutLinks) Summary() error {
	// write database here
	// but start just by dumping all the out-links
	return nil
}
