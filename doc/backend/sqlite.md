# sqlitebrowser

## Installation

sudo apt update
sudo apt install sqlitebrowser

sqlitebrowser

# Utilisation avec ton projet SheetFlow

Sélectionne ton fichier SQLite

./config/database.db

# Suppression admin

Il peut être utile de supprimer le compte admin dans le cas ou

curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"email":"admin@admin.com","password":"sheetable"}'
ne renvoie pas de token

Onglet Browse Data
Ouvre la base dans DB Browser for SQLite
Onglet Browse Data

Table users (ou user)
Sélectionne la ligne :
email = admin@admin.com
Efface complètement la valeur de la colonne password
Clique Write Changes

Ensuite, relance le backend
le code fait typiquement :

if admin n'existe pas → create

Dans ce cas, il faut redémarrer le serveur !!
