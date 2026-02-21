# Variables
BINARY_NAME=sf-backend
GO=go

# Récupération de la version via Git
# --tags: cherche les tags
# --always: si pas de tag, utilise le hash du commit
# --dirty: ajoute "-dirty" si tu as des modifs non committées
# Assure-toi d'avoir des tags dans ton repo pour que ça fonctionne bien, sinon tu auras juste le hash du commit.
# La version affiché se base toujours sur le dernier Tag réalisé à travers
# Exemple
#		git tag v0.1.0 -m "Première version stable avec affichage de version"
VERSION=$(shell git describe --tags --always --dirty)

# Le chemin vers ta variable dans le main.go
# Si ton fichier main.go est à la racine, c'est main.Version
LDFLAGS=-ldflags="-X main.Version=$(VERSION)"

# Règles de Makefile
.PHONY: all build run clean tidy reset help version

## all: Nettoie, organise les modules et compile
all: tidy build

## version: Affiche la version que le Makefile détecte (pratique pour débugger)
version:
	@echo "Version détectée : $(VERSION)"

## build: Compile le projet avec injection de la version
build:
	@echo "==> Compilation du binaire (Version: $(VERSION))..."
	CGO_ENABLED=0 $(GO) build $(LDFLAGS) -o build/$(BINARY_NAME) main.go

## Run : Lance le programme
run:
	@echo "==> Lance go.main (Version: $(VERSION))..."
	$(GO) run $(LDFLAGS) main.go

## tidy: Organise et nettoie le fichier go.mod
tidy:
	@echo "==> Nettoyage des modules (tidy)..."
	$(GO) mod tidy

## clean: Supprime le binaire local
clean:
	@echo "==> Suppression du binaire..."
	rm -f build/$(BINARY_NAME)

## reset: Nettoie le cache des modules et retélécharge tout
reset:
	@echo "==> Nettoyage complet du cache des modules..."
	$(GO) clean -modcache
	$(GO) mod tidy

## help: Affiche les commandes disponibles
help:
	@echo "Commandes disponibles :"
	@echo "  make build   : Compile le projet avec injection de version"
	@echo "  make run     : Lance le run avec injection de version"
	@echo "  make version : Affiche la version qui sera injectée"
	@echo "  make tidy    : Nettoie les dépendances"
	@echo "  make reset   : Vide le cache et réinstalle tout"
	@echo "  make clean   : Supprime le binaire"
	@echo "  make help    : Affiche cette aide"
