package wiki

import (
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

const (
	timeformat = "20060102-150405"
)

// Given an array of strings, convert this into a single valid file name.
func ValidBaseName(words []string) string {
	var b strings.Builder

	for i, s := range words {
		if i > 0 {
			b.WriteRune('-')
		}
		for _, r := range s {
			switch {
			case unicode.IsDigit(r),
				unicode.IsLetter(r):
				b.WriteRune(r)
			case unicode.IsSpace(r):
				b.WriteRune('-')
			case unicode.IsPunct(r):
				b.WriteRune('_')
			}
		}
	}

	return b.String()
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
