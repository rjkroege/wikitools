//go:build !darwin
// +build !darwin

// All non-darwin platforms as darwin

package search

import (
	"fmt"
	"log"

	"github.com/rjkroege/wikitools/corpus"
)

type portableWikilinkNameIndex struct {
}

var _ corpus.LinkToFile = (*portableWikilinkNameIndex)(nil)

func (_ *portableWikilinkNameIndex) Path(_, _, _ string) (string, error) {
	return "", fmt.Errorf("portableWikilinkNameIndex.Path not implemented")
}
func (_ *portableWikilinkNameIndex) Allpaths(_, _, _ string) ([]string, error) {
	return nil, fmt.Errorf("portableWikilinkNameIndex.Allpaths not implemented")
}

func MakeWikilinkNameIndex(_ string) *portableWikilinkNameIndex {
	log.Println("hi not from darwin")
	return &portableWikilinkNameIndex{}
}

func (_ *portableWikilinkNameIndex) Wikitext( _, _ string) (string, error) {
	return "", fmt.Errorf("portableWikilinkNameIndex.Wikitext not implemented")
}
