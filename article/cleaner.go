package article

import (
	"log"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rjkroege/wikitools/config"

)

func ShowFileInfo(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}
	log.Printf("%s: %v\n", path, info)
	return nil
}


type BatchCleaner struct {
	filesmissingmetadata []string
}

// skipper returns true for files that we don't want to process
func skipper(path string, info os.FileInfo) bool {
	relp, err := filepath.Rel(config.Basepath, path)
	if err != nil {
		return true // Always skip bad paths
	}

	switch {
	case strings.HasPrefix(relp, "templates"):
		return true
	case info.Name() == "README.md":
		return true
	}
	return false
}

func (abc *BatchCleaner) ModernizeMetadata(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println("couldn't read ", path, ": ", err)
		return fmt.Errorf("couldn't read %s: %v", path, err)
	}
	
	if skipper(path, info) {
		return nil
	}

	if !info.IsDir() && filepath.Ext(info.Name()) == ".md" {
		abc.DoMetadataUpdate(path)
	}
	return nil
}

func (abc *BatchCleaner) DoMetadataUpdate(path string) error {
		log.Println("should process markdown file", path)
	// figure out if we have Metadata
	// figure out if it's wrapped in ---
	// log the result nicely (to unsorted)

	d, err := os.Stat(path)
	if err != nil {
		log.Println("DoMetadataUpdate Stat error", err)
		return fmt.Errorf("can't DoMetadataUpdate Stat %s: %v", path, err)
	}

	fd, err := os.Open(path)
	if err != nil {
		log.Println("DoMetadataUpdate Open error", err)
		return fmt.Errorf("can't DoMetadataUpdate Open %s: %v", path, err)
	}
	defer fd.Close()

	// TODO(rjk): Return an error if it fails.
	md := MakeMetaData(filepath.Base(path), d.ModTime())
	md.RootThroughFileForMetadata(fd)

	if md.HadMetaData {
		log.Println(path, "had metadata")
	} else {
		log.Println(path, "no metadata")
		abc.filesmissingmetadata := append(abc.filesmissingmetadata, path)
	}

	return nil
}

