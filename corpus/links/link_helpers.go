package links

import "github.com/rjkroege/wikitools/corpus"

func AppendStringVector[T corpus.Link](f func(T) string, linkies map[string]corpus.LinkMap[T], articles map[string][]string) {
	for k, v := range linkies {
		for u := range v {
			articles[k] = append(articles[k], f(u))
		}
	}
}
