//go:build !darwin
// +build !darwin

// All non-darwin platforms as darwin

package search

import (
	"fmt"
	"log"

	"github.com/rjkroege/wikitools/corpus"
)

type portableWikilinkIndexer struct {
}

func (_ *portableWikilinkIndexer) Allpaths(_ string) ([]string, error) {
	return nil, fmt.Errorf("portableWikilinkIndexer not implemented")
}

func MakeWikilinkNameIndex() corpus.WikilinkNameIndex {
	log.Println("hi not from darwin")
	return &portableWikilinkIndexer{}
}
