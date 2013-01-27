package wiki

import (
    "strings"
)

func filter(r rune) rune {
    lut := map[rune]rune { 
        ' ':  '-',
        '/':  ',',
        '#':  ',',
        '\t': '-'  }
    nr, ok := lut[r]
    if !ok {
        return r
    }
    return nr
}

// Given an array of strings, convert this into a single valid file name.
func Validname(words []string) string {
    s := strings.Join(words, " ");
    return strings.Map(filter, s)
}
