# CORS : Cross-Origin Resource Sharing

## Définition

En développement web, CORS signifie Cross-Origin Resource Sharing (Partage de ressources entre origines multiples).
Pour faire simple : c'est un mécanisme de sécurité qui permet (ou empêche) un site web d'accéder à des ressources situées sur un autre domaine.

## Pourquoi le CORS existe-t-il ?

Par défaut, les navigateurs appliquent la règle de la Same-Origin Policy (Politique de même origine).
Cette règle interdit à une page web de faire des requêtes vers un domaine différent du sien
(par exemple, site-a.com ne peut pas appeler l'API de site-b.com).

Sans cette protection, un site malveillant pourrait facilement lire les données de votre session Facebook ou de votre banque si vous les aviez ouvertes dans un autre onglet.

## Comment ça fonctionne ?

Le CORS utilise des en-têtes HTTP pour permettre au serveur de "donner sa permission" au navigateur.
Voici le scénario typique :

- La requête : Votre application (http://mon-app.fr) envoie une requête Fetch ou Axios vers http://mon-api.com.
- L'autorisation : Le serveur de l'API doit répondre avec un en-tête spécifique :
- Access-Control-Allow-Origin: http://mon-app.fr (ou \* pour autoriser tout le monde, bien que ce soit moins sécurisé).
- Le verdict : Si cet en-tête est présent et correspond à l'origine de l'appelant, le navigateur laisse passer les données.
- Sinon, il bloque la réponse et vous obtenez la célèbre erreur rouge dans la console.

Les requêtes "Preflight" (Options)

Pour les requêtes "sensibles" (comme celles avec des méthodes PUT, DELETE ou des en-têtes personnalisés), le navigateur envoie d'abord une petite requête de vérification appelée Preflight avec la méthode OPTIONS.
Il demande au serveur : "Est-ce que j'ai le droit de faire cette opération ?". Si le serveur dit oui, la vraie requête est envoyée.

## Erreurs CORS courantes

Si vous voyez un message d'erreur CORS, cela signifie généralement que :

- Le serveur n'est pas configuré pour accepter des requêtes venant de votre domaine.
- L'en-tête Access-Control-Allow-Origin est manquant ou mal configuré côté serveur.
- Note importante : Le CORS est une sécurité côté navigateur.
  - Un outil comme Postman ou un script Python n'est pas soumis aux restrictions CORS, car ils ne sont pas des navigateurs web.

## GO

En Go, comme le langage est souvent utilisé pour bâtir des microservices ou des API performantes, la gestion du CORS est une étape quasi-obligatoire dès que votre frontend (React, Vue, etc.) est séparé de votre backend.

Voici les deux méthodes principales pour implémenter le CORS en Go.

1. Utiliser le package standard (net/http)

Si vous ne voulez pas de dépendances externes, vous pouvez créer un middleware simple qui ajoute les en-têtes nécessaires à chaque réponse.
Go

```go
func enableCors(next http.Handler) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r _http.Request) {
// Autoriser une origine spécifique (recommandé) ou "_"
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

        // Autoriser les méthodes HTTP spécifiques
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

        // Autoriser les en-têtes personnalisés (ex: Content-Type ou Authorization)
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Gérer la requête de "Preflight" (OPTIONS)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })

}
```

Ensuite, vous enveloppez votre routeur avec ce middleware :

```go
http.ListenAndServe(":8080", enableCors(router)).
```

2. Utiliser la bibliothèque rs/cors (Recommandé)

Pour les projets de production, on utilise souvent le package rs/cors car il gère proprement tous les cas complexes (comme les Wildcards, le caching des options, etc.).

```go
package main

import (
"net/http"
"github.com/rs/cors"
)

func main() {
mux := http.NewServeMux()
mux.HandleFunc("/api/data", func(w http.ResponseWriter, r \*http.Request) {
w.Write([]byte(`{"message": "Succès !"}`))
})

    // Configuration personnalisée
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000", "https://mon-site.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowCredentials: true,
        Debug:            true, // Utile pour débugger dans la console en dev
    })

    // On applique le middleware
    handler := c.Handler(mux)
    http.ListenAndServe(":8080", handler)

}
```

Points de vigilance en Go

- L'ordre des Middlewares : Le middleware CORS doit généralement être l'un des premiers appliqués. Si un autre middleware (comme une authentification) bloque la requête avant que les en-têtes CORS ne soient ajoutés, le navigateur rejettera la réponse.

- Le cas des OPTIONS : N'oubliez jamais que le navigateur envoie une requête OPTIONS avant les requêtes POST ou PUT. Si votre code ne renvoie pas un statut 200 OK (ou 204 No Content) à cette requête, le navigateur annulera la "vraie" requête.

- Sécurité : Évitez d'utiliser AllowedOrigins: []string{"\*"} en production si votre API manipule des données sensibles ou utilise des cookies/sessions.

## Test

Pour simuler ce que fait un navigateur (la fameuse requête Preflight), on utilise la méthode OPTIONS.

```shell
curl -v -X OPTIONS http://localhost:8080/api/data \
 -H "Origin: http://localhost:3000" \
 -H "Access-Control-Request-Method: POST"
```

Ce qu'il faut vérifier dans la réponse

Si votre code Go est bien configuré, vous devriez voir ces en-têtes dans la section des résultats (commençant par <) :
En-tête attendu Signification
HTTP/1.1 200 OK (ou 204) Le serveur a accepté la vérification.
Access-Control-Allow-Origin Doit afficher http://localhost:3000 (ou \*).
Access-Control-Allow-Methods Doit contenir POST (la méthode qu'on a testée).

### Pourquoi utiliser -v (verbose) ?

Sans le drapeau -v, curl ne vous montrera que le corps de la réponse (souvent vide pour un OPTIONS). Le mode verbeux vous permet de voir les headers, qui sont la seule chose qui compte pour le CORS.
