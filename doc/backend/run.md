# Étape obligatoire côté terminal

## Configurer le fichier config.go

Il faut utiliser un mot de passe d'application

Les mots de passe d'application ne s'utilisent pas avec Authenticator.

Utilisez la barre de recherche en haut de votre compte Google et tapez directement "Mots de passe d'application".

export SMTP_PASSWORD="mot_de_passe_application_gmail"

## Fichier environnement

.env

à la racine du backend.
Et y mettre uniquement ce dont tu as besoin.

Exemple minimal pour ton cas :

DB_DRIVER=sqlite
PORT=7373
DEV=1

# Run

make run
make build
