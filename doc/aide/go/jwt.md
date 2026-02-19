## Qu’est-ce qu’un JWT ?

JWT = **JSON Web Token**, c’est un **jeton sécurisé** qui contient des informations (claims) que tu peux transmettre entre un client et ton serveur.
Un JWT a trois parties :

- HEADER.PAYLOAD.SIGNATURE

* **Header** → indique l’algorithme de signature (`HS256`, par exemple).
* **Payload** → contient les claims : des données comme `user_id`, `exp` (expiration), `authorized`, etc.
* **Signature** → une signature HMAC avec ta clé secrète (`apiSecret`) pour que personne ne puisse falsifier le token.

Exemple concret:

```go
claims := jwt.MapClaims{}
claims["authorized"] = true
claims["user_id"] = user_id
claims["exp"] = time.Now().Add(time.Hour * 168).Unix() // expiration 1 semaine
```

On crée ensuite le token :

```go
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
signedToken, err := token.SignedString([]byte(apiSecret))
```

Le résultat est une **chaîne longue** que le client stocke (souvent dans le header `Authorization: Bearer <token>`).

## Création du token (`CreateToken`)

Quand l’utilisateur se connecte :

1.  Il fournit son email + mot de passe.
2.  On valide le mot de passe en base de données.
3.  Si tout est OK, on appelle `CreateToken(user_id, apiSecret)` → on obtient une chaîne JWT.
4.  On envoie ce token au client (front-end ou API caller).

C’est ce token qui servira pour **authentifier les requêtes suivantes**.

## Middleware Gin (`AuthMiddleware`)

Le middleware s’exécute **avant chaque handler protégé** :

```go
func AuthMiddleware() gin.HandlerFunc {
    secret := config.Config().ApiSecret

    return func(c *gin.Context) {
        tokenString := utils.ExtractToken(c)  // prend le token du header
        err := auth.TokenValid(tokenString, secret)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }
        c.Next() // Token valide → continue vers le handler
    }
}
```

Explication :

1.  On récupère le token depuis le header HTTP (souvent `Authorization: Bearer ...`).
2.  On vérifie que le token est **valide et non falsifié** grâce à `apiSecret`.
3.  Si c’est valide, Gin passe au handler. Sinon, on renvoie `401 Unauthorized`.

## Extra : récupérer `user_id` depuis le token

Souvent on souhaite savoir **quel utilisateur fait la requête**. Avec JWT v5 :

```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
    }
    return []byte(apiSecret), nil
})

claims := token.Claims.(jwt.MapClaims)
userID := uint32(claims["user_id"].(float64)) // On récupère l’ID
```

Ensuite

```go
c.Set("user_id", userID)
```

Ensuite dans le handler :

```go
userID := c.GetUint32("user_id")
```

Ça évite de reparser le token à chaque fois dans tes routes.

### 🔹 Résumé du flux complet :

1.  Client fait POST `/api/login` → envoie email + mot de passe.
2.  Serveur valide → crée JWT avec `user_id` et `apiSecret`.
3.  Serveur renvoie le JWT au client.
4.  Client stocke JWT (localStorage, session, etc.).
5.  Client fait requête protégée → envoie header `Authorization: Bearer <token>`.
6.  Middleware Gin récupère le token → vérifie signature + expiration.
7.  Middleware injecte `user_id` dans le contexte si token valide.
8.  Handler récupère `user_id` → effectue opérations sécurisées.
