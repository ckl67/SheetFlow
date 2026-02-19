package controllers

import (
	"backend/api/config"
	"backend/api/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gin-contrib/cors"
)

type Server struct {
	DB      *gorm.DB
	Router  *gin.Engine
	Version string
}

func (server *Server) Initialize(version string) {
	var err error

	server.Version = version

	// Charger la configuration du serveur
	cfg := config.Config()

	DbDriver := cfg.Database.Driver
	DbUser := cfg.Database.User
	DbPassword := cfg.Database.Password
	DbHost := cfg.Database.Host
	DbPort := cfg.Database.Port
	DbName := cfg.Database.Name

	// Logger configuration
	var gormLogger logger.Interface
	if cfg.Dev {
		log.Println("Running in DEV mode - SQL and GIN in logs enabled")
		// logger.Info affiche les requêtes SQL avec les valeurs des paramètres, ce qui est utile pour le développement et le débogage.
		gormLogger = logger.Default.LogMode(logger.Info)
		// gin.DebugMode affiche les requêtes HTTP entrantes, les paramètres, les en-têtes, etc.,
		// ce qui est également utile pour le développement et le débogage.
		gin.SetMode(gin.DebugMode)
		// Affichage de la configuration du serveur, y compris les secrets (à éviter en production)
		cfg.LogSafe()

	} else {
		log.Println("Running in PRODUCTION mode")
		// logger.Silent n'affiche aucune requête SQL, ce qui est recommandé en production pour éviter de divulguer des informations sensibles dans les logs.
		gormLogger = logger.Default.LogMode(logger.Silent)

		// gin.ReleaseMode désactive les logs de requêtes HTTP, ce qui est également recommandé en production pour éviter de divulguer des informations sensibles dans les logs.
		gin.SetMode(gin.ReleaseMode)

	}

	switch DbDriver {

	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			DbUser, DbPassword, DbHost, DbPort, DbName,
		)

		server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
			DbHost, DbPort, DbUser, DbName, DbPassword,
		)

		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})

	default: // sqlite
		if _, err := os.Stat(cfg.ConfigPath); os.IsNotExist(err) {
			_ = os.Mkdir(cfg.ConfigPath, os.ModePerm)
		}

		dbPath := path.Join(cfg.ConfigPath, "database.db")

		server.DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: gormLogger,
		})
	}

	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	log.Println("Database connected successfully")

	// ✅ Migration
	if err := server.DB.AutoMigrate(&models.User{}, &models.Sheet{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	server.SetupRouter()
}

func (server *Server) Run(addr string, dev bool) {
	fmt.Printf("Listening to port %v\n", addr)

	// Logger Gin (équivalent LoggingHandler)
	server.Router.Use(gin.Logger())
	server.Router.Use(gin.Recovery())

	// Activer CORS uniquement en mode dev
	// Frontend: http://localhost:3000
	// Backend:  http://localhost:8080
	// Origines différentes ⇒ CORS requis.

	if config.Config().CorsOrigin != "" {
		server.Router.Use(cors.New(cors.Config{
			AllowOrigins: []string{config.Config().CorsOrigin},
			AllowMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowHeaders: []string{
				"Origin",
				"Content-Type",
				"Authorization",
			},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	srv := &http.Server{
		Handler:      server.Router, // plus besoin de wrapper
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
