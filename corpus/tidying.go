package corpus


import (
	"fmt"
	"log"
	"os"
)

// Tidying is the interface implemented by each of the kinds of Tidying
// passes.
type Tidying interface {
	// EachFile is called by the filepath.Walk over each file in the wiki tree.
	// It could do something like (e.g.)
	EachFile(path string, info os.FileInfo, err error) error
	Summary() error
}


// ListAllWikiFiles is a boring implementation of Tidying that lists all files.
type listAllWikiFiles struct{}

func NewListAllWikiFilesTidying() (Tidying, error) {
	return &listAllWikiFiles{}, nil
}

func (_ *listAllWikiFiles) EachFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}
	log.Printf("%s: %v\n", path, info)
	return nil
}

func (_ *listAllWikiFiles) Summary() error {
	return nil
}
