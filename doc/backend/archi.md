# Architecture SheetFlow

Architecture backend Go orientée API REST avec séparation des responsabilités.
Une architecture backend Go monolithique structurée par couches

On est sur une organisation proche d’une Clean Architecture simplifiée, avec une dominante MVC enrichie
MVC (Modèle-Vue-Contrôleur) simplifié

Composant Rôle concret

- Models La représentation des données et la connexion à la base de données.
- Controllers Le "chef d'orchestre". Il reçoit la requête, appelle les bons services et renvoie une réponse.
- Middleware Le "videur" à l'entrée. Il vérifie par exemple si tu es connecté ou si tu as les droits avant même d'atteindre le contrôleur.
- Forms La gestion du nettoyage et de la validation des données envoyées par l'utilisateur (pour éviter que le contrôleur ne fasse 200 lignes de vérifications if/else).

## API

Une API (application programming interface ou « interface de programmation d'application ») est une interface logicielle qui permet de « connecter » un logiciel ou un service à un autre logiciel ou service afin d'échanger des données et des fonctionnalités.

## REST

Une architecture REST (Representational State Transfer) est un ensemble de règles qui permettent à deux systèmes informatiques **_(généralement un client et un serveur)_** de communiquer via le protocole HTTP.

### Les 6 Contraintes Fondamentales

Pour qu'une API soit qualifiée de "RESTful", elle doit respecter certains principes :

- Client-Serveur : Le client (l'interface) et le serveur (les données) sont séparés. On peut modifier le code du client sans toucher à la base de données du serveur.
- Sans État (Stateless) : C'est le point crucial. Le serveur ne garde aucun souvenir des requêtes précédentes. Chaque requête doit contenir toutes les informations nécessaires pour être traitée (jeton d'authentification, paramètres, etc.).
- Mise en cache (Cacheable) : Les réponses doivent indiquer si elles peuvent être mises en cache par le client pour améliorer les performances.
- Interface Uniforme : C’est ce qui rend l'API prévisible. On utilise des URL pour identifier les ressources et des méthodes HTTP standards.
- Système en couches : Le client ne sait pas s'il est connecté directement au serveur final ou à un intermédiaire (comme un pare-feu ou un équilibreur de charge).
- Code à la demande (Optionnel) : Le serveur peut envoyer du code exécutable (comme du JavaScript) au client.

### Le fonctionnement concret : Ressources et Méthodes

Dans le monde REST, tout est une Ressource (un utilisateur, une photo, une commande). Chaque ressource est identifiée par une URI (une adresse unique).
Les Verbes HTTP

On utilise les méthodes standards du protocole HTTP pour agir sur ces ressources :
Méthode Action Exemple d'URL

- GET Récupérer une ressource GET /utilisateurs/42
- POST Créer une nouvelle ressource POST /utilisateurs
- PUT Remplacer/Modifier une ressource PUT /utilisateurs/42
- DELETE Supprimer une ressource DELETE /utilisateurs/42 3. Le format des données

Bien que REST ne l'impose pas techniquement, le format JSON (JavaScript Object Notation) est devenu la norme universelle car il est léger et facile à lire pour les humains comme pour les machines.

Exemple d'une réponse de l'API - JSON

```shell
{
"id": 42,
"nom": "Jean Dupont",
"email": "jean@exemple.com"
}
```

# Vue d’ensemble

## Racine :

```text
.
├── api/
├── config/
├── go.mod
├── go.sum
└── main.go

```

| Bloc      | Rôle                                    |
| --------- | --------------------------------------- |
| `main.go` | Point d’entrée                          |
| `api/`    | Cœur applicatif (logique métier + HTTP) |
| `config/` | Données persistées et assets            |

## Entry Point

main.go
Bootstrapping de l’application

Initialisation config

Connexion base de données

Démarrage serveur HTTP

Il délègue ensuite à :

api/server.go
