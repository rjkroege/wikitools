//go:build darwin
// +build darwin

package search

import (
	"fmt"
	"log"

	"github.com/rjkroege/wikitools/corpus"
)

type spotlightWikilinkIndexer struct {
}

func (_ *spotlightWikilinkIndexer) Allpaths(_ string) ([]string, error) {
	return nil, fmt.Errorf("spotlightWikilinkIndexer not implemented")
}

func MakeWikilinkNameIndex() corpus.WikilinkNameIndex {
	log.Println("hi from darwin")
	return &spotlightWikilinkIndexer{}
}
