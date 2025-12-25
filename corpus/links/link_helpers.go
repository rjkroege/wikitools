package links

import (
	"iter"

	"github.com/rjkroege/wikitools/corpus"
)

func appendStringVector[T corpus.Link](f func(T) string, linkies map[string]corpus.LinkMap[T], articles map[string][]string) {
	for k, v := range linkies {
		for u := range v {
			articles[k] = append(articles[k], f(u))
		}
	}
}

func (links *Links) AppendStringVectorForwardLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	appendStringVector(f, links.ForwardLinks, articles)
}

func (links *Links) AppendStringVectorBackLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	appendStringVector(f, links.BackLinks, articles)
}

func (links *Links) AppendStringVectorOutUrls(f func(corpus.Urllink) string, articles map[string][]string) {
	appendStringVector(f, links.OutUrls, articles)
}

func (links *Links) AppendStringVectorDamagedLinks(f func(corpus.Wikilink) string, articles map[string][]string) {
	appendStringVector(f, links.DamagedLinks, articles)
}

func (links *Links) BackLinksIterator() iter.Seq[corpus.LinkTuple] {
	return func(yield func(corpus.LinkTuple) bool) {
		for k, v := range links.BackLinks {
			if !yield(corpus.LinkTuple{From: k, To: v}) {
				return
			}
		}
	}
}
