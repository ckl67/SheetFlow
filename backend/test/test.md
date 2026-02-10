# Test backend

Permet de tester toutes les fonctions du backend

# Public

| Méthode | URL        | Fonction     |
| ------- | ---------- | ------------ |
| GET     | `/health`  | health check |
| GET     | `/version` | version info |
| GET     | `/api`     | API root     |

curl http://localhost:8080/health
curl http://localhost:8080/version
curl http://localhost:8080/api

# Login

curl -X POST http://localhost:8080/api/login \
 -H "Content-Type: application/json" \
 -d '{"email":"admin@admin.com","password":"sheetflow"}'

Va retourner le token
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NzEyNzAxMjcsInVzZXJfaWQiOjZ9.GmN6ksFwjMq63Y3DaMv62IS8NsnxbhO3awWaX5rVPU4"

Ce token servira ensuite pour tester toutes les routes sécurisées (/users, /sheets, /search, etc.)

Login → récupérer JWT
Passer le token dans le header :
-H "Authorization: Bearer <JWT_TOKEN>"

# Sauvegarde TOKEN

TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
 -H "Content-Type: application/json" \
 -d '{"email":"admin@admin.com","password":"sheetflow"}' | tr -d '"')

echo "$TOKEN"

# Tester toutes les routes sécurisées

curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/users

| Méthode | URL                           | Fonction               |
| ------- | ----------------------------- | ---------------------- |
| POST    | `/api/login`                  | login                  |
| POST    | `/api/request_password_reset` | request password reset |

| POST | `/api/reset_password` | reset password |
| GET | `/api/users` | get all users |
| GET | `/api/users/:id` | get user by id |
| POST | `/api/users` | create user |
| PUT | `/api/users/:id` | update user |
| DELETE | `/api/users/:id` | delete user |
| GET | `/api/sheets` | get sheets page |
| POST | `/api/sheets` | get sheets page / search |
| GET | `/api/sheet/:sheetName` | get sheet by name |
| PUT | `/api/sheet/:sheetName` | update sheet |
| DELETE | `/api/sheet/:sheetName` | delete sheet |
| POST | `/api/upload` | upload file |
| PUT/POST | `/api/sheet/:sheetName/info` | update sheet info text |
| GET | `/api/sheet/pdf/:composer/:sheetName` | get PDF |
| GET | `/api/sheet/thumbnail/:name` | get thumbnail |
| GET | `/api/search/:searchValue` | search sheets |
| GET | `/api/search/composers/:searchValue` | search composers |
| POST | `/api/tag/sheet/:sheetName` | append tag |
| DELETE | `/api/tag/sheet/:sheetName` | delete tag |
| GET/POST | `/api/tag` | find sheets by tag |
| GET/POST | `/api/composers` | get composers page |
| PUT | `/api/composer/:composerName` | update composer |
| DELETE | `/api/composer/:composerName` | delete composer |
| GET | `/api/composer/portrait/:composerName` | serve portraits |
