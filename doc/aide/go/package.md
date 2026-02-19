# Rôle de `go.mod`

Le fichier `go.mod` dans le projet indique à Go :
« Voici le module, voici ma version de Go, et voici mes dépendances. »
Même localement, Go se réfère à ce fichier pour résoudre les imports et vérifier les versions des packages.

Le fichier pendant est le fichier mod.sum

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

# Incident : utilisation de go get -u

## Raison

La commande :

```shell
go get -u
```

- met à jour toutes les dépendances du module courant, y compris :
- les dépendances directes
- les dépendances indirectes (// indirect)
- les dépendances transitives (dépendances des dépendances)

Go tente alors de résoudre l’ensemble du graphe de dépendances avec les versions les plus récentes compatibles.

## Pourquoi c’est déconseillé en backend

Dans un projet backend stable :
Les dépendances doivent être verrouillées
Le build doit être reproductible
Les mises à jour doivent être contrôlées

```shell
go get -u
```

peut :

- Introduire des breaking changes
- Mettre à jour des modules indirects non maîtrisés
- Exiger une version plus récente de Go
- Modifier profondément go.mod et go.sum

C’est donc une commande risquée hors contexte de maintenance planifiée.

Il faut donc au plus vite revenir vers un mod.go ancien et repartir propre !!

## Où ont été stockés les modules téléchargés ?

Les modules téléchargés ont simplement été stockés dans le cache global.

Dans le cache global des modules :

go env GOMODCACHE
Généralement :
~/go/pkg/mod

Ce cache :

- Est partagé entre tous les projets
- N’influence pas les versions réellement utilisées
- Sert uniquement à éviter des téléchargements répétés

Les versions réellement utilisées sont déterminées uniquement par :

- go.mod
- go.sum

## Solution appliquée

revenir à une version antérieure de mod.go et mod.sum

Eventuellement par la suite mettre à jour

```shell
go get github.com/golang-jwt/jwt/v5@v5.3.1
go mod tidy
```

Sans utiliser -u.

## À propos du vidage du cache (go clean -modcache)

Commande :

```shell
go clean -modcache
```

Effet :

- Supprime tout le cache des modules
- Oblige Go à retélécharger les modules au prochain build

Important :

- Cela ne change pas les versions utilisées
- Cela ne corrige pas les erreurs de version
- Cela ne modifie pas go.mod

## Bonnes pratiques retenues

- Ne jamais utiliser go get -u globalement en production.
- Mettre à jour une dépendance à la fois.
- Toujours versionner go.mod et go.sum.
- Utiliser go mod tidy pour nettoyer, pas pour mettre à jour.
