# Test backend

Permet de tester toutes les fonctions du backend

# Content Type

-H "Content-Type: application/json"

-H dans curl signifie ajouter un header HTTP.
Donc cette instruction ajoute l’en-tête :
Content-Type: application/json
Autrement dit : comment interpréter les données envoyées.

url -X POST http://localhost:8080/api/login \
 -H "Content-Type: application/json" \
 -d '{"email":"test@mail.com","password":"1234"}'

Ici :

    -d envoie un body brut
    Content-Type: application/json
    Gin utilisera ShouldBindJSON() ou BindJSON()

Le serveur va parser le body comme du JSON

| Content-Type                        | Usage                   | Méthode Gin associée |
| ----------------------------------- | ----------------------- | -------------------- |
| `application/json`                  | API REST classique      | `ShouldBindJSON()`   |
| `multipart/form-data`               | Upload fichier + champs | `ShouldBind()`       |
| `application/x-www-form-urlencoded` | Form HTML simple        | `ShouldBind()`       |

# Public

| Méthode | URL        | Fonction     |
| ------- | ---------- | ------------ |
| GET     | `/health`  | health check |
| GET     | `/version` | version info |
| GET     | `/api`     | API root     |

## curl avec formattage json

Installe
sudo apt-get install jq

## Commande

curl http://localhost:8080/health
curl http://localhost:8080/version | jq
curl http://localhost:8080/api | jq

# Login

| Méthode | URL          | Fonction |
| ------- | ------------ | -------- |
| POST    | `/api/login` | login    |

0. Login
   TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@admin.com","password":"sheetflow"}' | tr -d '"')

echo "$TOKEN"

Va retourner le token
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NzEyNzAxMjcsInVzZXJfaWQiOjZ9.GmN6ksFwjMq63Y3DaMv62IS8NsnxbhO3awWaX5rVPU4"

Ce token servira ensuite pour tester toutes les routes sécurisées (/users, /sheets, /search, etc.)

Login → récupérer JWT
Passer le token dans le header :
-H "Authorization: Bearer <JWT_TOKEN>"

## Reset Password

| Méthode | URL                           | Fonction               |
| ------- | ----------------------------- | ---------------------- |
| POST    | `/api/request_password_reset` | request password reset |

curl -X POST http://localhost:8080/api/request_password_reset \
 -H "Content-Type: application/json" \
 -d '{"email":"christian.klugesherz@gmail.com"}'

Doit retourner par mail
Hey there was a password reset request to your accout. Go to

http://localhost:8080/reset-password/MElwDftAwhkkixCQgaAUHpsJeFvOPPQQKoTARWfm

to update your password

### Architecture correcte d’un reset password

Il y a toujours 2 étapes distinctes.

#### Étape 1 — L’utilisateur clique sur le lien

Lien reçu par email : http://localhost:3000/reset-password/<token>
⚠ Ça devrait être le FRONTEND, pas le backend.
Le frontend affiche :
Champ nouveau mot de passe
Champ confirmation

#### Étape 2 — Le frontend appelle le backend

Quand l’utilisateur valide :
Le frontend fait :
POST /api/reset_password
Avec JSON :
{
"token": "PIRrlBJxGPdZBoAAZALTibpIPgcYooRsJejMnnII",
"password": "nouveauMotDePasse"
}
Et là ton backend traite la demande.

#### Test

| Méthode | URL                   | Fonction       |
| ------- | --------------------- | -------------- |
| POST    | `/api/reset_password` | reset password |

curl -X POST http://localhost:8080/api/reset_password \
 -H "Content-Type: application/json" \
 -d '{"passwordResetId":"MElwDftAwhkkixCQgaAUHpsJeFvOPPQQKoTARWfm","password":"nouveauPassword123"}'

le token ici n’est PAS un JWT d’authentification.

Réponse :
{"id":2,"email":"christian.klugesherz@gmail.com","role":1,"password":"$2a$10$HWwKJlFToR1ienkITTWDCu0ZjxV9Y0NB4YDEG3bMTscy8XoQWIEH.","password_reset":"MElwDftAwhkkixCQgaAUHpsJeFvOPPQQKoTARWfm","password_reset_expire":"2026-02-11T14:08:55.056487833+01:00","created_at":"2026-01-28T11:12:36.01464077+01:00","updated_at":"2026-02-11T11:44:25.655026361+01:00"}

Vérification

curl -s -X POST http://localhost:8080/api/login \
 -H "Content-Type: application/json" \
 -d '{"email":"christian.klugesherz@gmail.com","password":"nouveauPassword123"}' | tr -d '"'

## Possible de message

"SMTP backend not configured. Go take a look at the docs to get started with emails."

Penser à configurer config.go

# Tester toutes les routes sécurisées

curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/users | jq

| Méthode  | URL                                    | Fonction                 |     |
| -------- | -------------------------------------- | ------------------------ | --- |
| GET      | `/api/users`                           | get all users            | 1   |
| POST     | `/api/users`                           | create user              | 4   |
| PUT      | `/api/users/:id`                       | update user              | 5   |
| DELETE   | `/api/users/:id`                       | delete user              |     |
| GET      | `/api/sheets`                          | get sheets page          | 2   |
| POST     | `/api/sheets`                          | get sheets page / search |     |
| PUT      | `/api/sheet/:sheetName`                | update sheet             |     |
| DELETE   | `/api/sheet/:sheetName`                | delete sheet             |     |
| POST     | `/api/upload`                          | upload file              |     |
| PUT/POST | `/api/sheet/:sheetName/info`           | update sheet info text   |     |
| POST     | `/api/tag/sheet/:sheetName`            | append tag               |     |
| DELETE   | `/api/tag/sheet/:sheetName`            | delete tag               |     |
| GET/POST | `/api/tag`                             | find sheets by tag       |     |
| GET/POST | `/api/composers`                       | get composers page       |     |
| PUT      | `/api/composer/:composerName`          | update composer          |     |
| DELETE   | `/api/composer/:composerName`          | delete composer          |     |
| GET      | `/api/users/:id`                       | get user by id           | 3   |
| GET      | `/api/sheet/:sheetName`                | get sheet by name        |     |
| GET      | `/api/sheet/pdf/:composer/:sheetName`  | get PDF                  |     |
| GET      | `/api/sheet/thumbnail/:name`           | get thumbnail            |     |
| GET      | `/api/search/:searchValue`             | search sheets            |     |
| GET      | `/api/search/composers/:searchValue`   | search composers         |     |
| GET      | `/api/composer/portrait/:composerName` | serve portraits          |     |

Nécessité de 0. pour la suite

1.  get all users
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/users | jq
2.  get sheets page
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/sheets | jq
3.  get user by id
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/users/1 | jq
4.  create user
    curl -X POST http://localhost:8080/api/users -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{ "email": "christian.klugesherz@gmail.com", "password": "Password123!"}'

    curl -X POST http://localhost:8080/api/users -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{ "email": "use1@user.com", "password": "PassWorduser1!"}'

    curl -X POST http://localhost:8080/api/users -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{ "email": "use2@user.com", "password": "PassWorduser2!"}'

5.  update user
    Exemple : utilisateur normal peut changer son email et password
    curl -X PUT http://localhost:8080/api/users/8 -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{ "email": "newmail@example.com", "password": "Password123!"}'
    Exemple : admin peut changer change mail password et role
    curl -X PUT http://localhost:8080/api/users/2 -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{ "email": "christian.klugesherz@gmail.com", "password": "Password123!", "role": 0}'

    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/sheet/:sheetName | jq
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/sheet/pdf/:composer/:sheetName | jq
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/sheet/thumbnail/:name | jq
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/search/:searchValue | jq
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/search/composers/:searchValue | jq
    curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/composer/portrait/:composerName | jq

/api/users/:id
/api/sheets
/api/sheet/:sheetName
/api/sheet/:sheetName
/api/upload
/api/sheet/:sheetName/info
/api/tag/sheet/:sheetName
/api/tag/sheet/:sheetName
/api/tag
/api/composers
/api/composer/:composerName
/api/composer/:composerName
