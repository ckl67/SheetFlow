# JSON Web Tokens

curl -X POST http://localhost:8080/api/login \
 -H "Content-Type: application/json" \
 -d '{"email":"admin@admin.com","password":"sheetflow"}'

Un JWT est défini par : https://github.com/dgrijalva/jwt-go?tab=readme-ov-file

Un JWT = 3 parties Base64URL séparées par des points :
HEADER.PAYLOAD.SIGNATURE

"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NzEyNzM3NzksInVzZXJfaWQiOjZ9.BcgkGDIjwe6qfcNz_k4YDSU0yJuqSsZxrqMWCFYgKRQ"

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
.
eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NzEyNzAxMjcsInVzZXJfaWQiOjZ9
.
GmN6ksFwjMq63Y3DaMv62IS8NsnxbhO3awWaX5rVPU4

👉 Le payload est lisible sans la clé secrète.

# Décodage réel de TON token : PAYLOAD

echo "eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NzEyNzAxMjcsInVzZXJfaWQiOjZ9" | base64 -d

Résultat :

{"authorized":true,"exp":1771270127,"user_id":6}

Parce que ton backend les met explicitement dans le JWT.

Dans le code fichier token.go

claims := jwt.MapClaims{
"authorized": true,
"user_id": user.ID,
"exp": time.Now().Add(time.Hour \* 24).Unix(),
}

Utilisé par ton middleware :

if claims["authorized"] != true {
return errors.New("unauthorized")
}

userID := claims["user_id"]
Timestamp Unix : date -d @1771270127

➡️ Le token expire automatiquement
➡️ Gin + jwt refusent le token après cette date

⚠️ Tout le monde peut lire le payload
Ce qui est sécurisé, c’est la signature :

# Debug

Dans config.go on teste uid qui doit ADMIN_UID = 1

    if uid != ADMIN_UID {
    	c.String(http.StatusUnauthorized, "Only admins are able to persue this command")
    	return
    }
