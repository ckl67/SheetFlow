package utils

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

// SanitizeName transforme une chaîne en une version "safe" pour fichiers, URLs, etc.
// Exemple : "École #1.pdf" → "Ecole-1.pdf"
func SanitizeName(name string) string {
	// 1️⃣ Translittération Unicode → ASCII
	s := unidecode.Unidecode(name)

	// 2️⃣ Minuscules
	s = strings.ToLower(s)

	// 3️⃣ Retirer les espaces en début/fin
	s = strings.TrimSpace(s)

	// 4️⃣ Remplacer les espaces et underscores par un tiret
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// 5️⃣ Garder uniquement [a-zA-Z0-9-.] → retirer tout le reste
	re := regexp.MustCompile(`[^a-zA-Z0-9\-.]+`)
	s = re.ReplaceAllString(s, "")

	// 6️⃣ Remplacer plusieurs tirets consécutifs par un seul
	re2 := regexp.MustCompile(`-+`)
	s = re2.ReplaceAllString(s, "-")

	// 5. Éventuellement, limiter la longueur pour éviter problème FS/DB
	if len(s) > 255 {
		s = s[:255]
	}

	return s
}
