package article

import (
	"log"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"io"
	"bufio"

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
	missingmd [][]string
	tmpl *template.Template
}

func MakeBatchCleaner() (*BatchCleaner, error) {
	tmpl, err :=  template.New("newstylemetadata").Parse(iawritermetadataformat)
	if err != nil {
		return nil, fmt.Errorf("can't MakeBatchCleaner %v", err)
	}
	return &BatchCleaner{
		missingmd: make([][]string, MdModern + 1),
		tmpl: tmpl,
	}, nil
}

func (abc *BatchCleaner) GenerateReport() {
	// TODO(rjk): generate reports of articles with invalid metadata
	// with links and such
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
		updatedpth, err := abc.updateMetadata(path)
		if err != nil {
			return err
		}
		if updatedpth != "" {
			if err := replaceFile(updatedpth, path); err != nil {
				return fmt.Errorf("swapFile can't: %v", err)
			}
		}
	}
	return nil
}

func replaceFile(newpath, oldpath string) error {
	backup := oldpath + ".back"
	if err := os.Link(oldpath, backup); err != nil {
		return fmt.Errorf("replaceFile backup: %v", error)
	}

	if err := os.Remove(oldpath); err != nil {
		return fmt.Errorf("replaceFile remove: %v", error)
	}

	if err := os.Link(newpath, oldpath); err != nil {
		return fmt.Errorf("replaceFile emplace: %v", error)
	}

	if err := os.Remove(newpath); err != nil {
		return fmt.Errorf("replaceFile remove: %v", error)
	}
	if err := os.Remove(backup); err != nil {
		return fmt.Errorf("replaceFile remove: %v", error)
	}
	return nil
}

func (abc *BatchCleaner) updateMetadata(path string) (string, error) {
	d, err := os.Stat(path)
	if err != nil {
		log.Println("updateMetadata Stat error", err)
		return "", fmt.Errorf("can't DoMetadataUpdate Stat %s: %v", path, err)
	}

	ifd, err := os.Open(path)
	if err != nil {
		log.Println("updateMetadata Open error", err)
		return "", fmt.Errorf("can't DoMetadataUpdate Open %s: %v", path, err)
	}
	defer ifd.Close()
	fd := bufio.NewReader(ifd)

	// TODO(rjk): RootThroughFileForMetadata needs to return an error when it fails
	md := MakeMetaData(filepath.Base(path), d.ModTime())
	md.RootThroughFileForMetadata(fd)

	// Collect and categorize the articles by their metadata state
	abc.missingmd[md.mdtype] = append(abc.missingmd[md.mdtype], path)

	if md.mdtype != MdLegacy {
		// Nothing to do.
		return "", nil
	}

	tpath := path + "-updating"
	nfd, err := os.Create(tpath)
	if err != nil {
		log.Println("replaceLegacyMetadata Create error", err)
		return "", fmt.Errorf("can't updateMetadata Create %s: %v", tpath, err)
	}
	defer nfd.Close()

		if err := abc.writeUpdatedMetadata(path, fd, nfd, md); err != nil {
			log.Println("DoMetadataUpdate", err)
			return "", fmt.Errorf("can't updateMetadata: %v",  err)
		}

	return tpath, nil
}

// There are other transformations that I'll want to implement. Refactor when
// I need to. Assumes that ofd's read point is at the end of the metadata in original file.
// 
func (abc *BatchCleaner) writeUpdatedMetadata(path string, ofd io.Reader, nfd io.Writer, md *MetaData) error {
	// write new metadata to nfd
	nmd := &IaWriterMetadataOutput{
		Title: md.Title,
		Date: md.DetailedDate(),
		Tags: md.Tagstring(),
		Extrakeys: md.extraKeys,
	}

//	log.Printf("nmd: %#v\n", nmd)

	if err := abc.tmpl.Execute(nfd, nmd); err != nil {
		log.Println("oops, bad template write because", err)
		return fmt.Errorf("can't writeUpdatedMetadata Execute template: %v", err)
	}

	// write existing file minus its metadata to it (first line after the first blank line)
	_, err := io.Copy(nfd, ofd)
	return 	err
}

type IaWriterMetadataOutput struct {
	Title             string
	Date    string
	Tags string
	Extrakeys map[string]string
}

const iawritermetadataformat = `---
title: {{.Title}}
date: {{.Date}}{{if ne .Tags "" }}
tags: {{.Tags}}{{end}}{{range $key, $value := .Extrakeys}}
{{$key}}: {{$value}}{{end}}
---

`
