package helpers

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		// TODO add error management
		panic(e)
	}
	return output
}

func NormalizeString(str string) string {
	// lower case
	str = strings.ToLower(str)
	// remove accents
	str = RemoveAccents(str)
	// remove spaces at beginning and ending of string
	str = strings.TrimSpace(str)

	return str
}
