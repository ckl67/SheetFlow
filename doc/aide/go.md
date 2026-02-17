# Go - Synthèse

# Introduction

Go est un langage de programmation expressif et concis.
Ses mécanismes de concurrence facilitent l’écriture de programmes exploitant au mieux les machines multicœurs et connectées en réseau.
Go se compile rapidement en code machine.
C’est un langage compilé, statiquement typé et rapide, mais qui reste simple et agréable à utiliser, à la manière d’un langage interprété.

## La compilation

La commande go build compile les fichiers sources de ton projet (ainsi que leurs dépendances) pour les transformer en un fichier exécutable binaire.

Ce qu'il se passe selon le contexte :

- Si tu es dans un "main package" : Go crée un fichier exécutable dans le répertoire courant.
  Sous Windows, ce sera nom_du_projet.exe, et sous Linux/macOS, simplement nom_du_projet.

- Si tu compiles une bibliothèque (package seul) : Go vérifie simplement que le code compile sans erreur,
  mais il ne génère pas de fichier exécutable (car il n'y a pas de fonction main à lancer).
  Exemple : go build api/forms/common.go

## Attention éviter de compiler un seul fichier !!!

En Go :
Un package peut contenir plusieurs fichiers
Les méthodes peuvent être réparties dans plusieurs fichiers
Mais elles appartiennent au même type si dans le même package

On ne fait presque jamais : go build fichier.go
On fait : go build
ou go build ./...

## run

La commande go run main.go : Compile et lance le programme immédiatement sans laisser de fichier derrière lui (pratique pour le développement).
go build : Crée le fichier final pour la distribution ou la mise en production.

# Documentation

[Documentation officielle de Go](https://go.dev/doc/)

# Installation

```bash
sudo apt remove golang-go
sudo apt purge golang-go
sudo apt autoremove

sudo rm -rf /usr/local/go

wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz
```

## Configuration Linux

Ajoutez cette ligne dans le fichier de configuration du shell `~/.bashrc` :

```bash
export PATH="$PATH:/usr/local/go/bin"
source ~/.bashrc

go version
```

Puis rechargez le shell :

```bash
source ~/.bashrc
```

## Point Important

Après upgrade :

```bash
go clean -modcache
go mod tidy
go build
```

Pour être certain que tout est reconstruit proprement avec la nouvelle toolchain.

```bash
go version
go env GOROOT
go env GOPATH
go env GOMOD
```

## Modifier dans le fichier go.md

```bash
go mod edit -go=1.24
go mod tidy
```

# Formatage

```bash
go fmt ./...
```

1. Le rôle de go fmt

La commande go fmt ./... est l'un des outils les plus appréciés de l'écosystème Go. Elle sert à formater automatiquement ton code source pour qu'il respecte strictement les standards officiels de style Go.

Go est un langage très d'opinion sur la présentation du code. Au lieu de débattre pour savoir s'il faut mettre des espaces ou des tabulations, ou où placer les accolades, la communauté utilise un outil unique.
go fmt réécrit tes fichiers pour :

- Remplacer les espaces par des tabulations pour l'indentation.
- Aligner correctement les colonnes dans les déclarations de structures ou de variables.
- Placer les accolades au bon endroit.
- Ajuster les espaces autour des opérateurs.

2. Que signifie ./... ?

C'est un "wildcard" (caractère générique) spécifique à l'outil Go :

- . : Partir du répertoire courant.
- /... : Chercher de manière récursive dans tous les sous-dossiers.
  Donc, go fmt ./... formate tout ton projet d'un seul coup, peu importe la profondeur des dossiers (comme ton api/forms/).

# Rôle de `go.mod`

Le fichier `go.mod` dans le projet indique à Go :
« Voici le module, voici ma version de Go, et voici mes dépendances. »
Même localement, Go se réfère à ce fichier pour résoudre les imports et vérifier les versions des packages.

La commande :

```bash
go mod tidy
```

fait deux choses principales :

1. **Ajoute les dépendances manquantes**
   - Si un package externe (ex. : `github.com/gin-gonic/gin`) n’est pas dans `go.mod`, `tidy` l’ajoute automatiquement.
2. **Supprime les dépendances inutilisées**
   - Si `go.mod` contient des packages qui ne sont plus importés dans le code, `tidy` les supprime.

# Installation

go get <package>

# Suppression

go get <package>@none

Le suffixe @none est une instruction spéciale pour l'outil Go. Elle signifie : "Supprime complètement cette dépendance de mon projet."

```shell
go get github.com/jinzhu/gorm@none
go get github.com/mattn/go-sqlite3@none
```

- github.com/mattn/go-sqlite3 (Le Pilote)
  C'est un driver (un pilote). C'est la couche logicielle "bas niveau" qui permet à Go de parler directement au fichier de base de données SQLite.

- github.com/jinzhu/gorm (Le Traducteur / ORM)
  C'est un ORM (Object-Relational Mapping). C'est une couche "haut niveau" qui se place au-dessus du driver.
  Son rôle : Il te permet de manipuler ta base de données en utilisant des structures Go au lieu d'écrire du SQL à la main.
  Au lieu d'écrire SELECT \* FROM users;, tu écris db.Find(&users).

# Package

En Go, le package est l’unité de base de l’organisation du code. C'est l'équivalent d'un "module" ou d'une "bibliothèque" dans d'autres langages.

## La notion de Package

Chaque fichier .go doit commencer par la déclaration package nom_du_package.
Le package main : C'est le seul qui génère un exécutable. Il doit contenir une fonction func main().
Les autres packages : Ce sont des bibliothèques de fonctions, structures ou variables réutilisables.
Visibilité (Encapsulation) :

- Si un nom commence par une Majuscule (User), il est exporté (public).
- S'il commence par une minuscule (user), il est privé au package.

## L'importation classique (La norme)

Normalement, quand tu importes un package, tu accèdes à son contenu via son nom :

```Go
import "fmt"

func main() {
    fmt.Println("Hello") // On utilise le préfixe 'fmt'
}
```

### L'utilisation du . (Dot Import)

Le "dot import" consiste à ajouter un point devant le chemin de l'import :

```Go
import . "fmt"

func main() {
    Println("Hello") // Plus besoin de 'fmt.' !
}
```

Mais c'est déconseillé et c'est une très mauvaise pratique pour plusieurs raisons :

- En important avec un ., tu "verses" toutes les fonctions du package directement dans ton fichier actuel.
- La perte de lisibilité

# Force le compilateur Go à ne pas utiliser les bibliothèques C du système hôte.

Par défaut, Go essaie parfois d'utiliser des bibliothèques C (via cgo) pour certaines fonctionnalités comme la résolution de noms de domaine (DNS) ou les certificats X.509.
Le problème ? Si tu compiles ton programme sur une machine avec une version spécifique de la bibliothèque C (glibc), ton binaire risque de ne pas démarrer sur une autre machine !!

Avec

```shell
CGO_ENABLED=0
```

On désactive CGO, tu demandes à Go de générer un binaire statique pur.
Autonomie totale : Le binaire contient tout ce dont il a besoin. Il ne dépend plus d'aucune bibliothèque .so ou .dll externe.
Portabilité maximale : Tu peux copier ce binaire sur n'importe quel système Linux (peu importe la distribution) et il s'exécutera sans erreur de dépendance.
Sécurité : Réduit la surface d'attaque en évitant de lier des bibliothèques système potentiellement vulnérables.

Il est donc important de positionner ce paramètre

C'est la commande standard pour créer des images Docker ultra-légères :
On compile en mode statique

```shell
CGO_ENABLED=0 GOOS=linux go build main.go
```

Le binaire fonctionnera car il n'a besoin d'aucune lib système

## Importation anonyme

Le tiret bas (\_) devant un import est ce qu'on appelle un "blank import" (importation anonyme).

En Go, si tu importes un package sans l'utiliser dans ton code, le compilateur génère une erreur.
Le \_ permet de dire à Go : "Je sais que je n'appelle aucune fonction de ce package directement, mais je veux quand même l'importer pour ses effets secondaires."

# Mise à jour Gorm

Le projet de base utilisait

## Gorm v1

```shell
"github.com/jinzhu/gorm"
"github.com/mattn/go-sqlite3"
```

Problème
github.com/jinzhu/gorm = GORM v1
mattn/go-sqlite3 = dépend de CGO (librairie C SQLite)

## Passage à Gorm v2

Suppression

```shell
go get github.com/jinzhu/gorm@none
go get github.com/mattn/go-sqlite3@none
```

Objectif

```shell
"gorm.io/gorm"
"gorm.io/driver/sqlite"
"gorm.io/driver/mysql"
"gorm.io/driver/postgres"

_ "modernc.org/sqlite"
```

Problème :
Le package gorm.io/driver/sqlite est maintenu par l'équipe officielle de GORM.
Le souci, c'est qu'il est codé pour utiliser par défaut github.com/mattn/go-sqlite3 (qui utilise du C).

Certains développeurs ajoutent \_ "modernc.org/sqlite" en espérant que cela "forcera" le programme à utiliser la version Pure Go.
Mais attention : Dans la majorité des cas, cela ne suffit pas.
Si tu gardes gorm.io/driver/sqlite dans ton code, GORM va quand même essayer d'appeler le code C au moment de la compilation.

**L'import anonyme** de modernc va simplement enregistrer un deuxième driver SQLite dans ton programme,
mais GORM continuera de pointer vers celui qui utilise C.

## Finallement la solution retenue

```shell
import (
  "github.com/glebarez/sqlite"
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "gorm.io/driver/postgres"
```

Toutes les librairies sont en Go pur !

Ce que fait réellement github.com/glebarez/sqlite
glebarez/sqlite est :

- Un driver SQLite pure Go
- Basé sur modernc.org/sqlite
- Adapté spécifiquement pour GORM v2
- Sans CGO

```shell
go get github.com/glebarez/sqlite
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get gorm.io/driver/postgres
go mod tidy
```

## Modification entre V1 et V2

❌ Ancien code (v1)
if gorm.IsRecordNotFoundError(err) {
return &Composer{}, errors.New("Composer not found")
}

✅ Nouveau code (v2)

gorm.IsRecordNotFoundError n’existe plus en GORM v2.

import "errors"

if errors.Is(err, gorm.ErrRecordNotFound) {
return &Composer{}, errors.New("Composer not found")
}

# Mise à jour gin

Tu as actuellement :

github.com/gin-gonic/gin v1.7.4

## Problèmes avec cette version :

Ancienne API : certaines méthodes ou signatures ont changé depuis v1.7 → GORM v2 et Go 1.24 pourraient exposer des incompatibilités subtiles.
Sécurité : correctifs de vulnérabilités sur JSON binding, middleware et context.
Performance : le framework a été optimisé depuis v1.7.
Compatibilité Go 1.24 : les anciennes versions peuvent générer des warnings ou comportements inattendus.

## Version recommandée

La dernière version stable est v1.9.x (2025+), compatible Go 1.24.
Pour Go 1.24, v1.9.0 ou supérieure est safe.

## Comment mettre à jour

go get github.com/gin-gonic/gin@latest
go mod tidy

## Nettoyer le module

go mod tidy -v

# Notions de programmes

## Receiver

Exemple :

```go
type User struct {
    ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
    Email     string    `gorm:"size:100;not null;unique" json:"email"`
    UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}


func (u *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {}
```

Nous sommes sur la même notion que sur une classe et méthode

- `func` → mot-clé pour définir une fonction en Go.
- `(u *User)` → **receiver**, la fonction appartient à la struct `User`.
- `u` est un pointeur vers `User`, ce qui permet à la fonction de modifier directement l’objet.

- `FindUserByEmail` → nom de la fonction. Les majuscules signifient que la fonction est exportée.
- `(db *gorm.DB, email string)` → paramètres :
  - `db *gorm.DB` → connexion à la base de données Gorm
  - `email string` → email à rechercher

- `(*User, error)` → valeurs de retour : pointeur vers `User` et erreur éventuelle.

Exemple d’utilisation :

```go
var u User
foundUser, _ := u.FindUserByEmail(db, "bob@example.com")
```

# Slice

Le problème avec un tableau c’est qu’il a une taille fixe, il faut donc absolument connaître sa taille au moment de sa déclaration
ce qui n’est pas vraiment évident car on peut très vite s’apercevoir plus tard dans notre code que la taille allouée au départ à notre tableau est insuffisante.

C’est là qu’interviennent les Slices dans Go.
Ils vont nous permettre d’avoir un tableau flexible et le dimensionner de façon dynamique sans se soucier de sa taille pendant sa déclaration.

```go
  var nombres = []int{0, 0, 0, 0, 0} // création d'une slice avec 5 éléments
```

Pour rajouter un élément dans votre slice il faut utiliser la fonction append(), qui prend comme paramètres d'abord votre slice et
ensuite l'élément que vous voulez rajouter et elle vous retournera une nouvelle Slice avec l'élément rajouté.

# GORM

GORM est une librairie externe.

## Go standard (sans GORM)

Il faut écrire ses requêtes SQL à la main, gérer le scan et les erreurs.

```go
row := db.QueryRow("SELECT id, email FROM users WHERE email = ?", email)
row.Scan(&user.ID, &user.Email)

```

## Go avec GORM

Exemple :

```go
db.Where("email = ?", email).First(&user)
```

### Complément

ORM existe dans beaucoup de langages :

- Java → Hibernate
- Python → Django ORM / SQLAlchemy
- PHP → Eloquent

## Couches orthogonales :

Exemple

```go
// Ici GORM doit savoir :
// que SafeSheetName est la clé primaire
// que UploaderID ne peut pas être NULL
// que Tags est un text[] PostgreSQL !! ATTENTION A CETTE NOTION QUI LIMITE LE TYPE DE BASE !

type Sheet struct {
SafeSheetName string `gorm:"primary_key" json:"safe_sheet_name"`
SheetName string `json:"sheet_name"`
SafeComposer string `json:"safe_composer"`
Composer string `json:"composer"`
ReleaseDate time.Time
PdfUrl string `json:"pdf_url"`
UploaderID uint32 `gorm:"not null" json:"uploader_id"`
CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
Tags pq.StringArray `gorm:"type:text[]" json:"tags"` // C’est un champ PostgreSQL text[].
Categories pq.StringArray `gorm:"type:text[]" json:"categories"`
InformationText string `json:"information_text"`
}
```

Il faut bien distinguer deux aspects

- gorm:"..." → mapping base de données (ORM → schéma SQL)
- json:"..." → mapping API HTTP (struct Go ↔ JSON)

Ce sont deux tags totalement indépendants.

### gorm:"..." → contrat avec la base de données

Le tag GORM influence :
le type SQL généré
les contraintes (PRIMARY KEY, NOT NULL, UNIQUE…)
la taille (size:255)
les valeurs par défaut
les index
le nom de colonne

### json:"..." → contrat avec ton API

Le tag JSON sert uniquement pour :

c.JSON(...)
c.BindJSON(...)

Exemple :

SheetName string `json:"sheet_name"`

Cela signifie :
côté API → la clé JSON sera sheet_name
sans tag → ce serait SheetName
Donc json concerne la sérialisation HTTP, pas la base.

### Pourquoi certains champs n’ont pas gorm ?

Parce que :
➜ GORM sait déjà quoi faire par convention
Exemple :
SheetName string `json:"sheet_name"`
Composer string `json:"composer"`

Par défaut GORM :
type Go string → type SQL varchar(255) (selon driver)
nom colonne → sheet_name (snake_case automatique)

### Pourquoi certains champs n’ont pas json ?

Exemple :

ReleaseDate time.Time
Sans tag JSON :
la clé sera ReleaseDate dans le JSON (camel case exact)
pas release_date

## Hooks

GORM détecte automatiquement certaines méthodes :

- `BeforeCreate`
- `AfterCreate`
- `BeforeSave`
- `AfterSave`

Exemple :

```go
db.Create(&u)
```

Appelle automatiquement :

```
BeforeSave()
BeforeCreate()
INSERT
AfterCreate()
AfterSave()
```

## Take

```go
user := User{}
db.Where("email = ?", email).Take(&user)
user.PasswordReset = ...
db.Save(&user)
```

`Take()` est nécessaire pour modifier l’objet chargé.

## Vision

```go
Model(&User{}).Where("email = ?", email)
```

Équivalent à :

```go
_, err := user.FindUserByEmail(db, email)
Model(&user) // GORM connaît déjà l’ID pour générer WHERE id = ...
```

# Gin

## Définition

Gin est un framework HTTP pour Go, surcouche de `net/http`.

- Routes : `r.GET()`, `r.POST()`, `r.PUT()`, `r.DELETE()`
- Middleware :

```go
r.Use(AuthMiddleware())
```

## Architecture typique :

Client → Gin → Services → GORM → Database

## JWT (JSON Web Token)

Mécanisme d’authentification **stateless** :

- Le serveur ne garde pas de session en mémoire.
- Toute l’information nécessaire est contenue dans le token.

Un JWT comporte trois parties :

```
HEADER.PAYLOAD.SIGNATURE
```

## Bind, ShouldBind, BindJSON et ShouldBindWith (Gin)

Dans Gin, le binding permet de parser les données entrantes d’une requête HTTP (JSON, formulaire, multipart, etc.) vers une structure Go.

### Exemple

```go
// Un champ est obligatoire seulement si tu ajoutes :binding:"required"
// Exemple : SheetName string `form:"sheetName" binding:"required"`
type UploadRequest struct {
  File \*multipart.FileHeader `form:"uploadFile"`
  Composer string `form:"composer"`
  SheetName string `form:"sheetName"`
  ReleaseDate string `form:"releaseDate"`
  Categories string `form:"categories"`
  Tags string `form:"tags"`
  InformationText string `form:"informationText"`
}

// Fonction gin
// ShouldBind :
// Parse la requête HTTP selon le Content-Type
// Remplit la struct
// Retourne une erreur uniquement si le binding échoue
var uploadForm forms.UploadRequest
	if err = c.ShouldBind(&uploadForm); err != nil {
		utils.DoError(c, http.StatusBadRequest, fmt.Errorf("bad upload request: %v", err))
		return
	}


// Signifie :
// Gin va chercher un champ nommé "composer" dans le multipart/form-data
// Si il trouve un champ "composer", il va convertir sa valeur en string
// Si la conversion réussit, il va assigner la valeur convertie à uploadForm.Composer

```

### ShouldBind

```go
   if err := c.ShouldBind(&obj); err != nil {
   // gestion manuelle de l'erreur
   }
```

Fonctionnement

Choisit automatiquement le binder en fonction du Content-Type.
Ne modifie pas automatiquement la réponse HTTP.
Retourne une erreur si le parsing échoue.
Quand l’utiliser ?

- Recommandé en production
- Quand on veut contrôler la gestion des erreurs

### Bind

```go
   c.Bind(&obj)
```

Fonctionnement

Choisit automatiquement le binder selon Content-Type.
En cas d’erreur :
Écrit automatiquement un HTTP 400
Stoppe le traitement du handler

Inconvénient

- Ne laisse pas le contrôle sur la réponse d’erreur.
  Quand l’utiliser ?
- Usage simple ou rapide, mais déconseillé en production.

### BindJSON / ShouldBindJSON

```go
   c.ShouldBindJSON(&obj)
```

Fonctionnement

Force l’utilisation du binder JSON.
Ignore la détection automatique du Content-Type.
Différence
BindJSON → écrit automatiquement HTTP 400 en cas d’erreur.
ShouldBindJSON → retourne l’erreur sans écrire la réponse.

Quand l’utiliser ?

- API REST pure JSON
- Quand on veut forcer le parsing JSON

### ShouldBindWith

```go
   c.ShouldBindWith(&obj, binding.JSON)
```

Fonctionnement

Permet de choisir explicitement le binder.
Ignore la détection automatique du Content-Type.

Exemples
c.ShouldBindWith(&obj, binding.JSON)
c.ShouldBindWith(&obj, binding.Form)
c.ShouldBindWith(&obj, binding.FormMultipart)

Quand l’utiliser ?

- Cas avancés
- Quand le Content-Type ne peut pas être fiable

### En production, privilégier :

ShouldBind
ShouldBindJSON
ShouldBindWith
afin de garder le contrôle total sur la gestion des erreurs et les réponses HTTP.
