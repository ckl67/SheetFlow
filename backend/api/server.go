package api

import (
	"backend/api/controllers"
	"backend/api/seed"
	"fmt"

	. "backend/api/config"
)

var server = controllers.Server{}

func Run(version string) {
	// Fonction server.Initialize(version) est définie dans base.go
	// server.Initialize() appelle server.SetupRouter() (fichier routes.go)
	server.Initialize(version)

	seed.Load(server.DB, Config().AdminEmail, Config().AdminPassword)

	port := 8080
	if Config().Port != 0 {
		port = Config().Port
	}

	server.Run(fmt.Sprintf("0.0.0.0:%d", port), Config().Dev)
}

func RunWithPort(port int, version string) {
	// To run modules from cloud-backend-services controller

	server.Initialize(version)

	seed.Load(server.DB, Config().AdminEmail, Config().AdminPassword)

	server.Run(fmt.Sprintf("0.0.0.0:%d", port), Config().Dev)
}
