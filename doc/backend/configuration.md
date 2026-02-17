# config.go

## Principe

Fichier backend/api/config/config.go

Dans config.go nous définissons

```go
func Config() ServerConfig {
    configOnce.Do(func() {
        serverConfig = ConfigBuilder().Build()
    })
    return serverConfig
}
```

C’est un singleton : configOnce.Do(...) garantit que la configuration n’est construite qu’une seule fois.
Build() va :

- Créer serverConfig avec NewConfig() → valeurs par défaut
- Tenter de surcharger avec .env
- Tenter de surcharger avec les variables d’environnement

Donc à chaque fois que l'on appelle config.Config(), on récupère la configuration initialisée.
config.go n’est pas “auto-exécuté” : ses fonctions sont appelées lorsqu’une autre partie du code fait Config().

```go
unc (server *Server) Initialize(version string) {
	var err error

	server.Version = version

	// Set Release Mode
	if !Config().Dev {
		gin.SetMode(gin.ReleaseMode)
	}

	DbDriver := Config().Database.Driver
	DbUser := Config().Database.User
	DbPassword := Config().Database.Password
	DbHost := Config().Database.Host
	DbPort := Config().Database.Port
	DbName := Config().Database.Name

     // Silence the logger
    server.DB.LogMode(false)

```

Dans le code (par ex. main.go), on as souvent :

```go
server := &Server{}
server.Initialize("v1.0")
```

À l’intérieur de Initialize() ou dans d’autres fonctions, tu utilises Config() pour récupérer des variables :

## Résumé du flux

- main.go → crée Server → appelle server.Initialize()
- Initialize() → lit Config().Database, Config().Dev, etc.
- Config() appelle ConfigBuilder().Build() une seule fois (singleton)
- Build() → NewConfig() → surcharge .env → surcharge variables d’environnement

Les valeurs configurées sont ensuite utilisées pour :

- Connexion à la base
- Configuration du port
- Logs
- Etc.

En pratique : toutes les lectures de config passent par Config(), et c’est le moment où .env et les valeurs par défaut sont combinées.

```
      ┌─────────────────────┐
      │      main.go        │
      │                     │
      │ server := &Server{} │
      │ server.Initialize() │
      └─────────┬───────────┘
                │
                ▼
      ┌───────────────────────────┐
      │ Server.Initialize(version)│
      │                           │
      │ 1. Lire Config().Database │
      │ 2. Lire Config().Dev      │
      │ 3. Connexion DB           │
      │ 4. Setup Router / Logs    │
      └─────────┬─────────────────┘
                │
                ▼
      ┌────────────────────────────┐
      │     config.Config()        │  <--- singleton
      │                            │
      │  - configOnce.Do()         │
      │  - Appelle ConfigBuilder() │
      └─────────┬──────────────────┘
                │
                ▼
      ┌────────────────────────────────────────────┐
      │     ConfigBuilder().Build()                │
      │                                            │
      │  1. serverConfig = NewConfig()             │
      │     → initialise valeurs par défaut        │
      │       (SQLite, API_SECRET, etc.)           │
      │  2. Surcharge .env (feeder.DotEnv)         │
      │  3. Surcharge variables environnement      │
      │  4. Retourne ServerConfig final            │
      └─────────┬──────────────────────────────────┘
                │
                ▼
      ┌────────────────────────────┐
      │   serverConfig singleton   │
      │   Accessible via Config()  │
      │   → Dev, Port, Database…   │
      └────────────────────────────┘

```

# Création fichier d'environnement

Crée un .env

```shell
DEV=1
PORT=8080
```

Lancer :

go run main.go
