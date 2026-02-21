package config

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

// ========================================
// Principes de configuration du serveur
// ========================================
// Flus de lecture
// 1. NewConfig() initialise valeurs par défaut
// 2. DotEnv feeder charge .env
// 3. Env feeder surcharge variables système
// Cette dernière partie est gérée par
//	github.com/golobby/config/v3
//	Attention ne pas mélanger os.Getenv() et golobby/config

// Configuration du serveur SMTP pour l'envoi d'emails (ex: réinitialisation de mot de passe, notifications, etc.)
type SmtpConfig struct {
	Enabled        string `env:"SMTP_ENABLED"`
	From           string `env:"SMTP_FROM"`
	HostServerAddr string `env:"SMTP_HOST"`
	HostServerPort int    `env:"SMTP_PORT"`
	Username       string `env:"SMTP_USERNAME"`
	Password       string `env:"SMTP_PASSWORD"`
}

// Configuration de la base de données, utilisée pour se connecter à la base de données et effectuer des opérations CRUD.
type DatabaseConfig struct {
	Driver   string `env:"DB_DRIVER"`
	Host     string `env:"DB_HOST"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	Port     int    `env:"DB_PORT"`
}

// ServerConfig est la struct qui contient tous les paramètres de configuration du serveur.
type ServerConfig struct {
	AdminEmail    string `env:"ADMIN_EMAIL"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
	ApiSecret     string `env:"API_SECRET"`
	ServerUrl     string `env:"SERVER_URL"`
	ConfigPath    string `env:"CONFIG_PATH"`

	Dev  bool `env:"DEV"`
	Port int  `env:"PORT"`

	Database   DatabaseConfig
	Smtp       SmtpConfig
	CorsOrigin string `env:"CORS_ORIGIN"` //"https://app.sheetflow.com" ou "http://localhost:3000" pour dev, ou "*" pour autoriser toutes les origines
}

// configBuilder est une struct utilisée pour construire la configuration du serveur à partir d'un fichier .env et des variables d'environnement.
type configBuilder struct {
	dotenvFile           string
	errorOnMissingDotenv bool
}

// Singleton pattern pour la configuration du serveur
// Le sync.Once est l'outil standard en Go pour implémenter proprement le Pattern Singleton.
// Il garantit qu'une action spécifique ne soit exécutée qu'une seule fois,
// même si ton code essaie de l'appeler depuis plusieurs endroits en même temps
var (
	// La "boîte" qui contient les données
	serverConfig ServerConfig // 's' minuscule = privé au package

	// Le "gardien" qui se souvient si on a déjà rempli la boîte
	configOnce sync.Once
)

func (c ServerConfig) LogSafe() {
	log.Println("------ SERVER CONFIG (DEV MODE) ------")
	log.Printf("AdminEmail: %s\n", c.AdminEmail)
	// log.Printf("AdminPassword: %s\n", c.AdminPassword) // Affiche les secrets de l'administrateur, à éviter en production
	// log.Printf("ApiSecret: %s\n", c.ApiSecret)         // Affiche la configuration de base du serveur, y compris les secrets (à éviter en production)
	log.Printf("ServerUrl: %s\n", c.ServerUrl)
	log.Printf("ConfigPath: %s\n", c.ConfigPath)
	log.Printf("Dev mode: %v\n", c.Dev)
	log.Printf("Port: %d\n", c.Port)

	log.Println("Database:")
	log.Printf("  Driver: %s\n", c.Database.Driver)
	log.Printf("  Host: %s\n", c.Database.Host)
	log.Printf("  User: %s\n", c.Database.User)
	log.Printf("  DB Password: %s\n", c.Database.Password) // Affiche la configuration de la base de données, y compris le mot de passe (à éviter en production)
	log.Printf("  Name: %s\n", c.Database.Name)
	log.Printf("  Port: %d\n", c.Database.Port)
	log.Printf("	CORS Origin: %s\n", c.CorsOrigin)

	log.Println("SMTP:")
	log.Printf("  Enabled: %s\n", c.Smtp.Enabled)
	log.Printf("  From: %s\n", c.Smtp.From)
	log.Printf("  Host: %s\n", c.Smtp.HostServerAddr)
	log.Printf("  Port: %d\n", c.Smtp.HostServerPort)
	log.Printf("  Username: %s\n", c.Smtp.Username)
	log.Printf("  Password: %s\n", c.Smtp.Password) // Affiche la configuration SMTP, y compris le mot de passe (à éviter en production)
	log.Println("--------------------------------------")
}

// ConfigBuilder() retourne une instance de configBuilder qui est utilisée pour construire la configuration du serveur
// à partir d'un fichier .env et des variables d'environnement.
// Elle permet de spécifier le fichier .env à utiliser et si une erreur doit être levée si le fichier .env est manquant.
func ConfigBuilder() configBuilder {
	return configBuilder{}
}

// WithDotenvFile() permet de spécifier le fichier .env à utiliser pour charger la configuration du serveur.
// Par défaut, le fichier .env utilisé est ".env" à la racine du projet.
func (b configBuilder) WithDotenvFile(file string) configBuilder {
	b.dotenvFile = file
	return b
}

// PanicOnMissingDotenv() permet de spécifier si une erreur doit être levée si le fichier .env est manquant.
// Si errorOnMissingDotenv est vrai et que le fichier .env est manquant, la fonction Build() panique avec un message d'erreur.
// Sinon, elle continue à charger la configuration à partir des variables d'environnement.
func (b configBuilder) PanicOnMissingDotenv(status bool) configBuilder {
	b.errorOnMissingDotenv = status
	return b
}

// La fonction Config() utilise le pattern singleton pour s'assurer que la configuration est chargée une seule fois
// et est accessible globalement dans l'application.
// La fonction Config() est utilisée dans tout le code pour accéder à la configuration du serveur,
// par exemple dans server.go, users_controller.go, etc.
// Elle retourne une instance de ServerConfig qui contient tous les paramètres de configuration nécessaires pour le serveur.
func Config() ServerConfig {
	// La magie est ici : le code à l'intérieur de la fonction anonyme
	// ne s'exécutera qu'une seule fois dans toute la vie du programme.

	configOnce.Do(func() {
		fmt.Println("Chargement de la configuration...")
		// Ici on appelle le builder pour remplir serverConfig
		serverConfig = ConfigBuilder().Build()
	})
	return serverConfig
}

// La fonction Build() est appelée par Config() pour construire l'instance de ServerConfig en chargeant
// les paramètres de configuration à partir du fichier .env (si spécifié) et des variables d'environnement.
// Elle utilise la bibliothèque golobby/config pour charger la configuration dans la struct ServerConfig.
// Si le fichier .env est manquant et que errorOnMissingDotenv est vrai, elle panique avec un message d'erreur.
// Sinon, elle continue à charger la configuration à partir des variables d'environnement et retourne l'instance de ServerConfig.
func (b configBuilder) Build() ServerConfig {
	serverConfig = NewConfig()

	dotenvFile := ".env"
	if b.dotenvFile != "" {
		dotenvFile = b.dotenvFile
	}
	dotenvFeeder := feeder.DotEnv{Path: dotenvFile}
	envFeeder := feeder.Env{}

	err := config.New().AddStruct(&serverConfig).AddFeeder(dotenvFeeder).Feed()
	if err != nil {
		if strings.Contains(err.Error(), "no such file") && b.errorOnMissingDotenv {
			log.Fatalf("error loading config from dotenv file %s: %s", dotenvFile, err.Error())
		}
	}
	err = config.New().AddStruct(&serverConfig).AddFeeder(envFeeder).Feed()
	if err != nil {
		log.Fatalf("error loding config from environemnt: %s", err.Error())
	}
	return serverConfig
}

// Cette fonction NewConfig() est utilisée pour définir les valeurs par défaut de la configuration du serveur.
// Ces valeurs par défaut sont utilisées si les paramètres correspondants ne sont pas définis dans le fichier .env ou les variables d'environnement.
// Par exemple, si DB_DRIVER n'est pas défini dans le fichier .env ou les variables d'environnement,
// la valeur par défaut "sqlite" sera utilisée pour la configuration de la base de données.
func NewConfig() ServerConfig {
	return ServerConfig{
		AdminEmail:    "admin@admin.com",
		AdminPassword: "sheetflow",
		ApiSecret:     "sheetflow_secret_key",
		ServerUrl:     "http://localhost:8080",
		ConfigPath:    "./config/",
		CorsOrigin:    "",
		Database: DatabaseConfig{
			Driver: "sqlite",
		},
		Smtp: SmtpConfig{
			Enabled:        "1", // 1 = activé
			From:           "christian.klugesherz@gmail.com",
			HostServerAddr: "smtp.gmail.com",
			HostServerPort: 587,
			Username:       "christian.klugesherz@gmail.com",
			Password:       "", // récupéré depuis variable d'environnement : SMTP_PASSWORD
		},
	}
}
