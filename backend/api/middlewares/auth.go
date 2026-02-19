package middlewares

import (
	"backend/api/auth"
	"backend/api/config"
	"backend/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware vérifie que le JWT fourni est valide.
// Si oui, il injecte user_id dans le contexte Gin pour le handler.
func AuthMiddleware() gin.HandlerFunc {
	secret := config.Config().ApiSecret

	return func(c *gin.Context) {
		tokenString := utils.ExtractToken(c) // récupère le token depuis le header Authorization

		// Vérifie que le token est valide
		if err := auth.TokenValid(tokenString, secret); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Récupère l'user_id et l'ajoute au contexte
		userID, err := auth.ExtractTokenID(tokenString, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
