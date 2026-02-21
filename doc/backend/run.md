# Étape obligatoire avant run

Principes de configuration du serveur

1. NewConfig() initialise valeurs par défaut
2. DotEnv charge fichier .env
3. Env feeder surcharge variables système

## Configurer le fichier api/config.config.go

A travers NewConfig()

Pour le mot de passe SMTP :
export SMTP_PASSWORD="ton_mot_de_passe_application_gmail"

Exemple Google :
Utilisez la barre de recherche en haut de votre compte Google et tapez directement "Mots de passe d'application".

## Fichier environnement

.env

A la racine du backend y mettre uniquement ce dont tu as besoin.

```shell
DB_DRIVER=sqlite
PORT=7373
DEV=1
```

# Run et Compilation

## Version

Nous suivons a Semantic Versioning :

- Major.Minor.Patch

## Injection

- ldflags : ("linker flags" ) Passe des arguments au programme de liaison (linker).
- -X : Indique au linker que l'on va définir la valeur d'une variable.
- main.var : Le chemin complet vers la variable.

Attention : Les drapeaux de compilation (comme -ldflags) doivent impérativement se trouver avant le nom du fichier ou du package.

```shell
go build -ldflags="-X main.Version=0.2.3" -o build/sf-backend  main.go

go run -ldflags="-X main.Version=0.2.3" main.go
# Plusieurs varibles
go run -ldflags="-X main.Version=1.2.3 -X main.BuildDate=2024-05-20" main.go
```

## Injection avec version git

### creation tag

Grace à la comme git

```shell
git describe --tags --abbrev=0
```

Exemple
Basculez sur Master :
Ramenez les modifications de Dev vers Master :
Créer le premier Tag (v0.1.0)
Affichage Tag
Pousser le tag vers le serveur

```shell
git checkout master
git merge dev
git tag v0.1.0 -m "Première version stable avec affichage de version"
git push origin v0.1.0

GIT_VER=$(git describe --tags --abbrev=0)
echo $GIT_VER
--> va afficher v0.1.0
```

### Run et Compilation avec Tag git

L'option --always est une sécurité : si aucun tag n'est trouvé, elle renvoie au moins le hash du commit au lieu d'une erreur).

```shell
go build -ldflags="-X main.Version=0.2.3" -o build/sf-backend main.go
go run -ldflags="-X main.Version=0.2.3" main.go

go build -ldflags="-X main.Version=$(git describe --tags --always)" -o build/sf-backend main.go
go run -ldflags="-X main.Version=$(git describe --tags --always)" main.go
```

# Run

Par la suite en lançant le build nous avons la bonne version
./build/sf-backend
