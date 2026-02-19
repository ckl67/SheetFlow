package utils

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

var slugRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)

// Transforme une chaîne de caractères en ASCII, en supprimant les accents et autres caractères spéciaux
// Exemple : "Élève" devient "Eleve"
func ToASCII(input string) string {
	ascii := unidecode.Unidecode(input)
	return strings.TrimSpace(ascii)
}

// Transforme une chaîne de caractères en un slug URL-friendly
// Exemple : "Hello World!" devient "hello-world"
func ToSlug(input string) string {
	ascii := unidecode.Unidecode(input)
	ascii = strings.ToLower(ascii)
	ascii = slugRegex.ReplaceAllString(ascii, "-")
	ascii = strings.Trim(ascii, "-")
	return ascii
}
