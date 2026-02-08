Le projet utilise :

github.com/mattn/go-sqlite3
github.com/mattn/go-sqlite3 est le driver SQLite le plus utilisé en Go,
go-sqlite3 connecte le backend Go à la base SQLite

SQLite est écrit en C

Le driver fait le pont Go ↔ C

Ton code Go ne voit que l’API Go, jamais le C

Le prpjet utilise SQLite

./config/database.db

    Fichier .db → signature classique de SQLite
