package utils

import "strings"

// SafeFileName normalise un nom de fichier
func SafeFileName(name string) string {
	return sanitizeName(ToASCII(name)) // <-- sans "utils."
}

// sanitizeName applique un nettoyage basique
func sanitizeName(name string) string {
	// exemple : supprime les caractères spéciaux
	return strings.Map(func(r rune) rune {
		if r == '_' || r == '-' || r == '.' || r == ' ' || ('0' <= r && r <= '9') || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			return r
		}
		return -1
	}, name)
}
