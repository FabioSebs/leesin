package scraper

import "strings"

func removeExtraWhitespace(s string) string {
	words := strings.Fields(s)
	return strings.Join(words, " ")
}
