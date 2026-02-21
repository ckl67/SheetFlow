# Variables
BINARY_NAME=sheetflow-backend
GO=go

.PHONY: all build clean tidy reset help

## all: Nettoie, organise les modules et compile
all: tidy build

## build: Compile le projet en binaire statique (parfait pour Docker/Linux)
build:
	@echo "==> Compilation du binaire..."
	CGO_ENABLED=0 $(GO) build -o build/$(BINARY_NAME) main.go

## tidy: Organise et nettoie le fichier go.mod
tidy:
	@echo "==> Nettoyage des modules (tidy)..."
	$(GO) mod tidy

## clean: Supprime le binaire local
clean:
	@echo "==> Suppression du binaire..."
	rm -f build/$(BINARY_NAME)

## reset: Nettoie le cache des modules et retélécharge tout (ton go clean -modcache)
reset:
	@echo "==> Nettoyage complet du cache des modules..."
	$(GO) clean -modcache
	$(GO) mod tidy

## help: Affiche les commandes disponibles
help:
	@echo "Commandes disponibles :"
	@echo "  make build  : Compile le projet"
	@echo "  make tidy   : Nettoie les dépendances"
	@echo "  make reset  : Vide le cache et réinstalle tout"
	@echo "  make clean  : Supprime le binaire"
	@echo "  make help   : Affiche cette aide"
