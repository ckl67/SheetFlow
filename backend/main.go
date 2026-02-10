package main

import (
	"backend/api"
	"backend/api/utils"
)

// main.go appelle dans la package api, api.Run() (fichier server.go)
// api.Run() appelle server.Initialize() (fichier base.go)
// server.Initialize() appelle server.SetupRouter() (fichier routes.go)
// SetupRouter() (dans routes.go) définit toutes les routes Gin

func main() {
	var version string

	version = "v1.0.0"
	utils.PrintAsciiVersion(version) // affiche une bannière ASCII avec la version du serveur fichier version.go
	api.Run(version)                 // appelle api.Run() dans server.go pour démarrer le serveur Gin
}
