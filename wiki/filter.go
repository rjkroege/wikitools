package wiki

import (
	"strings"
	"unicode"
)

// ValidName converts a given string to a valid filename without extension.
func ValidName(s string) string {
	var b strings.Builder
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

	return b.String()
}

// Given an array of strings, convert this into a single valid file name.
// TODO(rjk): Can have a better name.
func ValidBaseName(words []string) string {
	var b strings.Builder

	for i, s := range words {
		if i > 0 {
			b.WriteRune('-')
		}
		b.WriteString(ValidName(s))
	}

	return b.String()
}
