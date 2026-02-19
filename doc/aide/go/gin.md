# Servery Go

```go
 http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("Hello"))}
      )
```

Cette ligne de code est la brique de base pour créer un serveur web en Go.
Elle définit une route (un "endpoint") et le comportement à adopter quand quelqu'un la visite.

Voici le détail de ce qui se passe :

- http.HandleFunc("/", ...)
  Cette fonction enregistre un handler (un gestionnaire) auprès du serveur par défaut de Go.
  "/" : C'est le chemin (le "pattern").
  Ici, la racine.
  En Go, le slash seul agit comme une "catch-all" :
  Il répondra à toutes les URLs qui ne correspondent pas à une autre route plus précise (ex: /contact, /api).

- func(w ..., r ...) :
  C'est une fonction anonyme qui sera exécutée à chaque fois qu'une requête arrive sur ce chemin.

  Les deux arguments clés (w et r)

  Le serveur Go vous donne deux outils pour travailler :

  w http.ResponseWriter : C'est votre "stylo". Vous l'utilisez pour construire la réponse que vous renvoyez au client (le navigateur). Vous pouvez y écrire du texte, du JSON, changer le code HTTP (200, 404), etc.

  r \*http.Request : C'est la "lettre" reçue. Elle contient tout ce que le client a envoyé : les paramètres d'URL, les headers (CORS !), les cookies, le corps du message (Body), etc.

- w.Write([]byte("Hello"))

  w.Write : Envoie des données au client.

  []byte("Hello") : La méthode Write n'accepte pas directement des string. Elle demande un "slice de bytes". On convertit donc la chaîne de caractères "Hello" en données brutes pour qu'elles puissent être envoyées sur le réseau.

Pour que cette ligne fonctionne, elle doit être placée dans une fonction main et suivie par le démarrage du serveur :

```go
package main

import "net/http"

func main() {
    // 1. On définit la route
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World"))
    })

    // 2. On lance le serveur sur le port 8080
    // (Nil signifie qu'on utilise le routeur par défaut qu'on a configuré au-dessus)
    http.ListenAndServe(":8080", nil)
}
```

# Wrapper

Un wrapper est une fonction qui prend un handler et en retourne un autre, en ajoutant un comportement autour.

Schéma mental :

Request → Logging → CORS → Router → Response

En code classique :

```go
handlerAvecCors := c.Handler(server.Router)
handlerFinal := handlers.LoggingHandler(os.Stdout, handlerAvecCors)
```

- server.Router est un http.Handler
- c.Handler(...) le wrappe avec CORS
- LoggingHandler(...) le wrappe encore

C’est une chaîne de handlers imbriqués.

Visuellement :

LoggingHandler(
CorsHandler(
Router
)
)

# Gin et middleware

Gin introduit le concept de middleware interne.

Quand tu fais :

```go
router.Use(middleware)
```

- Gin enregistre ce middleware dans une chaîne interne.
  Quand une requête arrive :

Request
↓
Middleware 1
↓
Middleware 2
↓
Route Handler
↓
Response

La différence clé :

👉 Avec Gin, tu n’as plus besoin de wrapper manuellement.

Gin construit la chaîne pour toi.

server.Router.Use(cors.New(...))
server.Router.Use(gin.Logger())

est conceptuellement équivalent à :

LoggingHandler(
CorsHandler(
Router
)
)

En conclusion

Avec Gin natif

Tu déclares les middlewares sur le router :

router.Use(...)

Puis :

http.Server{
Handler: router
}

🔴 Ancien modèle (wrapping)

```text
  ServeHTTP
    └── LoggingHandler.ServeHTTP
          └── CorsHandler.ServeHTTP
                  └── GinRouter.ServeHTTP
                        └── Route handler
```

Appels imbriqués.

🟢 Nouveau modèle (Gin middleware chain)

```text
  ServeHTTP
    ↓
  GinRouter
    ↓
  for each middleware:
    call middleware(context)
    if context.Next():
        continue
    else:
        stop
    ↓
  Route handler
```
