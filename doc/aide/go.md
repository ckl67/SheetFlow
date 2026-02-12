# Go - Synthèse

# Introduction

Go est un langage de programmation expressif et concis.
Ses mécanismes de concurrence facilitent l’écriture de programmes exploitant au mieux les machines multicœurs et connectées en réseau.
Go se compile rapidement en code machine.
C’est un langage compilé, statiquement typé et rapide, mais qui reste simple et agréable à utiliser, à la manière d’un langage interprété.

# Documentation

[Documentation officielle de Go](https://go.dev/doc/)

# Rôle de `go.mod`

Le fichier `go.mod` indique à Go :
« Voici le module, voici ma version de Go, et voici mes dépendances. »
Même localement, Go se réfère à ce fichier pour résoudre les imports et vérifier les versions des packages.

## `go mod tidy`

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

## Configuration Linux

Ajoutez cette ligne dans le fichier de configuration du shell `~/.bashrc` :

```bash
export PATH="$PATH:/usr/local/go/bin"
```

Puis rechargez le shell :

```bash
source ~/.bashrc
```

# Notions de programmes

## Receiver

Exemple :

```go
func (u *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {}
```

- `func` → mot-clé pour définir une fonction en Go.
- `(u *User)` → **receiver**, la fonction appartient à la struct `User`. `u` est un pointeur vers `User`, ce qui permet à la fonction de modifier directement l’objet.

Exemple de struct :

```go
type User struct {
    ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
    Email     string    `gorm:"size:100;not null;unique" json:"email"`
    UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
```

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

# GORM

GORM est une librairie externe.

## Go standard (sans GORM)

```go
import "database/sql"
```

Il faut écrire ses requêtes SQL à la main, gérer le scan et les erreurs.

## Go avec GORM

Exemple :

```go
db.Where("email = ?", email).First(&user)
```

Equivalent Go pur :

```go
row := db.QueryRow("SELECT id, email FROM users WHERE email = ?", email)
row.Scan(&user.ID, &user.Email)
```

### Complément

ORM existe dans beaucoup de langages :

- Java → Hibernate
- Python → Django ORM / SQLAlchemy
- PHP → Eloquent

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

Gin est un framework HTTP pour Go, surcouche de `net/http`.

- Routes : `r.GET()`, `r.POST()`, `r.PUT()`, `r.DELETE()`
- Middleware :

```go
r.Use(AuthMiddleware())
```

Architecture typique :

```
Client → Gin → Services → GORM → Database
```

# JWT (JSON Web Token)

Mécanisme d’authentification **stateless** :

- Le serveur ne garde pas de session en mémoire.
- Toute l’information nécessaire est contenue dans le token.

Un JWT comporte trois parties :

```
HEADER.PAYLOAD.SIGNATURE
```
