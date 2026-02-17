package seed

import (
	"backend/api/models"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Cette fonction est appelée dans server.go après l'initialisation du serveur et la connexion à la base de données.
// Elle utilise GORM pour effectuer une migration automatique des tables User, Sheet et Composer.
// Ensuite, elle crée un utilisateur administrateur avec les informations fournies (email et mot de passe).
func Load(db *gorm.DB, email string, password string) {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Sheet{},
		&models.Composer{},
	); err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	var existing models.User

	// Vérification de l'existence de l'utilisateur administrateur
	// Si l'utilisateur existe déjà, on ne fait rien pour éviter les doublons et on sor de la fonction.
	err := db.Where("email = ?", email).First(&existing).Error
	if err == nil {
		fmt.Println("User already exists")
		return
	}

	// Si l'erreur n'est pas "record not found", cela signifie qu'il y a un problème avec la base de données, et on panique avec un message d'erreur.
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Fatalf("unexpected DB error: %v", err)
	}

	// Création de l'utilisateur administrateur
	result := db.Create(&models.User{
		Email:    email,
		Password: password,
	})

	if result.Error != nil {
		log.Fatalf("cannot create admin user: %v", result.Error)
	}

	fmt.Println("Admin user created")
}
