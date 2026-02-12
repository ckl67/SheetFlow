Le projet utilise :

github.com/mattn/go-sqlite3
github.com/mattn/go-sqlite3 est le driver SQLite le plus utilisé en Go,
go-sqlite3 connecte le backend Go à la base SQLite

SQLite est écrit en C

Le driver fait le pont Go ↔ C

Le code Go ne voit que l’API Go, jamais le C

Le propjet utilise SQLite

./config/database.db

    Fichier .db → signature classique de SQLite

#

Au lancement Le serveur Backend créée

- config
  database.db

Attention, il faut que admin ait l'id=1 !!

# Mot de passe

Architecture actuelle validée

Ton flux est maintenant :

POST /request_password_reset

Génération token

Stockage en base + expiration 1h

Envoi mail

POST /reset_password avec :

passwordResetId

password

Vérification expiration

Hash bcrypt

Invalidation du token

Refus si réutilisation

C’est un flux standard production-ready.
