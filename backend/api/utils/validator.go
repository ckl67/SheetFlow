package utils

// Fonctions génériques

import (
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	once     sync.Once
)

// Validator retourne une instance singleton du validateur
// Exemple d'utilisation : utils.Validator().Struct(myStruct)
func Validator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
	})
	return validate
}

// ValidateStruct valide une structure en utilisant le validateur singleton
// Exemple d'utilisation : utils.ValidateStruct(myStruct)
func ValidateStruct(s interface{}) error {
	return Validator().Struct(s)
}

// SanitizeUser applique les normalisations basiques
// exemple : "  John.Doe@EXAMPLE.COM  " devient "john.doe@example.com"
func SanitizeUserEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
