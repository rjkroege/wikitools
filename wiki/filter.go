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
		'/':  ',',
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

func UniqueValidName(basepath string, filename string, extension string, system System) string {
	p := filepath.Join(basepath, filename+extension)
	if system.Exists(p) {
		return filepath.Join(basepath, filename+"-"+system.Now().Format(timeformat)+extension)
	}
	return p
}
