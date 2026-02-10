package seed

import (
	"backend/api/models"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

func Load(db *gorm.DB, email string, password string) {
	err := db.AutoMigrate(&models.User{}, &models.Sheet{}, &models.Composer{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	// Essaie de créer un utilisateur avec email et password
	// Si l’utilisateur existe déjà, GORM renvoie une erreur
	// Tu attrapes l’erreur et tu fais juste :
	err = db.Model(&models.User{}).Create(&models.User{
		Email:    email,
		Password: password,
	}).Error
	if err != nil {
		fmt.Println("User already exists")
		return
	}
}
