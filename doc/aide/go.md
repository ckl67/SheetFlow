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

# Slice

Le problème avec un tableau c’est qu’il a une taille fixe, il faut donc absolument connaître sa taille au moment de sa déclaration ce qui n’est pas vraiment évident car on peut très vite s’apercevoir plus tard dans notre code que la taille allouée au départ à notre tableau est insuffisante.

Et c’est là qu’interviennent les Slices dans Go, Ils vont nous permettre d’avoir un tableau flexible et le dimensionner de façon dynamique sans se soucier de sa taille pendant sa déclaration.

Il existe deux façons pour créer une Slice.

- Soit à partir de la même syntaxe qu'un tableau sauf que cette fois-ci on ne spécifie pas la taille du tableau :
  var nombres = []int{0, 0, 0, 0, 0} // création d'une slice avec 5 éléments

Pour rajouter un élément dans votre slice il faut utiliser la fonction append(), qui prend comme paramètres d'abord votre slice et ensuite l'élément que vous voulez rajouter et elle vous retournera une nouvelle Slice avec l'élément rajouté.

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

### Exemple

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

Il faut bien distinguer deux couches orthogonales :

- gorm:"..." → mapping base de données (ORM → schéma SQL)
- json:"..." → mapping API HTTP (struct Go ↔ JSON)

Ce sont deux tags totalement indépendants.

1️⃣ gorm:"..." → contrat avec la base de données

Le tag GORM influence :
le type SQL généré
les contraintes (PRIMARY KEY, NOT NULL, UNIQUE…)
la taille (size:255)
les valeurs par défaut
les index
le nom de colonne

Exemples dans ta struct :

SafeSheetName string `gorm:"primary_key" json:"safe_sheet_name"`
UploaderID uint32 `gorm:"not null" json:"uploader_id"`
Tags pq.StringArray `gorm:"type:text[]" json:"tags"`
CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

Ici GORM doit savoir :

que SafeSheetName est la clé primaire

que UploaderID ne peut pas être NULL

que Tags est un text[] PostgreSQL

que CreatedAt a une valeur par défaut SQL

Sans ces tags, GORM utiliserait ses conventions par défaut.

2️⃣ json:"..." → contrat avec ton API

Le tag JSON sert uniquement pour :

c.JSON(...)

c.BindJSON(...)

json.Marshal / json.Unmarshal

Exemple :

SheetName string `json:"sheet_name"`

Cela signifie :

côté API → la clé JSON sera sheet_name

sans tag → ce serait SheetName

Donc json concerne la sérialisation HTTP, pas la base.

3️⃣ Pourquoi certains champs n’ont pas gorm ?

Parce que :

➜ GORM sait déjà quoi faire par convention

Exemple :

SheetName string `json:"sheet_name"`
Composer string `json:"composer"`

Par défaut GORM :

type Go string → type SQL varchar(255) (selon driver)

nom colonne → sheet_name (snake_case automatique)

4️⃣ Pourquoi certains champs n’ont pas json ?

Exemple :

ReleaseDate time.Time

Sans tag JSON :

la clé sera ReleaseDate dans le JSON (camel case exact)

pas release_date

Donc si tu veux contrôler le contrat API, tu ajoutes json.

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

# Différences entre Bind, ShouldBind, BindJSON et ShouldBindWith (Gin)

Dans Gin, le binding permet de parser les données entrantes d’une requête HTTP (JSON, formulaire, multipart, etc.) vers une structure Go.

1. ShouldBind
   if err := c.ShouldBind(&obj); err != nil {
   // gestion manuelle de l'erreur
   }

Fonctionnement

Choisit automatiquement le binder en fonction du Content-Type.
Ne modifie pas automatiquement la réponse HTTP.
Retourne une erreur si le parsing échoue.
Quand l’utiliser ?

✅ Recommandé en production
✅ Quand on veut contrôler la gestion des erreurs

2. Bind
   c.Bind(&obj)

Fonctionnement

Choisit automatiquement le binder selon Content-Type.
En cas d’erreur :
Écrit automatiquement un HTTP 400
Stoppe le traitement du handler

Inconvénient

❌ Ne laisse pas le contrôle sur la réponse d’erreur.

Quand l’utiliser ?

Usage simple ou rapide, mais déconseillé en production.

3. BindJSON / ShouldBindJSON
   c.ShouldBindJSON(&obj)

Fonctionnement

Force l’utilisation du binder JSON.
Ignore la détection automatique du Content-Type.
Différence
BindJSON → écrit automatiquement HTTP 400 en cas d’erreur.
ShouldBindJSON → retourne l’erreur sans écrire la réponse.

Quand l’utiliser ?

✅ API REST pure JSON
✅ Quand on veut forcer le parsing JSON

4. ShouldBindWith
   c.ShouldBindWith(&obj, binding.JSON)

Fonctionnement

Permet de choisir explicitement le binder.

Ignore la détection automatique du Content-Type.

Exemples
c.ShouldBindWith(&obj, binding.JSON)
c.ShouldBindWith(&obj, binding.Form)
c.ShouldBindWith(&obj, binding.FormMultipart)

Quand l’utiliser ?

Cas avancés

Quand le Content-Type ne peut pas être fiable

Tests spécifiques

Résumé
Méthode Binder auto Gestion auto erreur Recommandé
Bind Oui Oui (HTTP 400) ❌
ShouldBind Oui Non ✅
BindJSON Non (JSON) Oui (HTTP 400) ⚠️
ShouldBindJSON Non (JSON) Non ✅
ShouldBindWith Non (manuel) Non Avancé
Recommandation

En production, privilégier :

ShouldBind

ShouldBindJSON

ShouldBindWith

afin de garder le contrôle total sur la gestion des erreurs et les réponses HTTP.
