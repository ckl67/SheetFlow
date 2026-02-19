package api

import (
	"backend/api/config"
	"backend/api/controllers"
	"backend/api/seed"
	"fmt"
)

var server = controllers.Server{}

// La fonction Run() est le point d'entrée principal pour démarrer le serveur Gin.
// Elle initialise le serveur, charge les données de seed, et démarre le serveur sur le port spécifié dans la configuration ou par défaut (8080).
func Run(version string) {
	// Fonction server.Initialize(version) est définie dans base.go
	// server.Initialize() appelle server.SetupRouter() (fichier routes.go)
	server.Initialize(version)

	// La fonction seed.Load() est définie dans seeder.go
	// Elle effectue une migration automatique des tables User, Sheet et Composer,
	// puis crée un utilisateur administrateur avec les informations fournies (email et mot de passe) dans la configuration.
	seed.Load(server.DB, config.Config().AdminEmail, config.Config().AdminPassword)

	port := 8080
	if config.Config().Port != 0 {
		port = config.Config().Port
	}

	server.Run(fmt.Sprintf("0.0.0.0:%d", port), config.Config().Dev)
}

func RunWithPort(port int, version string) {
	// To run modules from cloud-backend-services controller

	server.Initialize(version)

	seed.Load(server.DB, config.Config().AdminEmail, config.Config().AdminPassword)

	server.Run(fmt.Sprintf("0.0.0.0:%d", port), config.Config().Dev)
}
