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

# Run

make run
make build

ou
go run -ldflags="-X main.Version=1.2.3" main.go
