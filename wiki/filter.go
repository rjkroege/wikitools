package wiki

import (
	"path/filepath"
	"strings"
	"time"
)

const (
	timeformat = "20060102-150405"
)

func filter(r rune) rune {
	lut := map[rune]rune{
		' ':  '-',
		',':  '-',
		'/':  '_',
		'#':  ',',
		'\t': '-'}
	nr, ok := lut[r]
	if !ok {
		return r
	}
	return nr
}

// Given an array of strings, convert this into a single valid file name.
func ValidBaseName(words []string) string {
	s := strings.Join(words, " ")
	return strings.Map(filter, s)
}

type System interface {
	Exists(path string) bool
	Now() time.Time
}

// UniqueValidName creates new names by inserting the current time
// between the filename and the extension. Returns only the filename.
func UniqueValidName(basepath string, filename string, extension string, system System) string {
	fn := filename + extension
	if system.Exists(filepath.Join(basepath, fn)) {
		return filename + "-" + system.Now().Format(timeformat) + extension
	}
	return fn
}
