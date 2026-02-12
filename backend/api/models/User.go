package models

import (
	"backend/api/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// gorm:"size:10" sert à GORM pour créer la colonne correspondante en base de données
// avec un VARCHAR(10) (ou équivalent selon le moteur SQL).
// Une contrainte UNIQUE garantit que toutes les valeurs d’une colonne (ou combinaison de colonnes) sont distinctes.
// email,
// si mis sur PasswordReset , va interdir que 2 valeurs soient vide "" !!
//
// json:"password_reset"
// 👉 Interprété uniquement par le package encoding/json
// Cela contrôle : le nom du champ dans le JSON ; le binding (BindJSON) ; la sérialisation (c.JSON())

// User struct représente un utilisateur dans la base de données.
// ID : identifiant unique de l'utilisateur (clé primaire, auto-incrémenté).
// Email : adresse email de l'utilisateur (doit être unique).
// Role : rôle de l'utilisateur (0 pour admin, 1 pour normal).
// Password : mot de passe de l'utilisateur (stocké sous forme de hash).
// PasswordReset : token utilisé pour la réinitialisation du mot de passe.
// PasswordResetExpire : date d'expiration du token de réinitialisation.
// CreatedAt : date de création de l'utilisateur.
// UpdatedAt : date de dernière mise à jour de l'utilisateur.
type User struct {
	ID                  uint32    `gorm:"primary_key;auto_increment" json:"id"` // If ID == 1: user = admin
	Email               string    `gorm:"size:100;not null;unique" json:"email"`
	Role                uint8     `json:"role"` // 0=admin 1=normal,
	Password            string    `gorm:"size:100;not null;" json:"password"`
	PasswordReset       string    `gorm:"size:64;" json:"password_reset"` /* Random char string for resetting the password (prob not the best implementation of a password reset so it could be redone)*/
	PasswordResetExpire time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"password_reset_expire"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// PrepareForCreate() doit rester une fonction de création, pas d’update.
func (u *User) PrepareForCreate() {
	u.ID = 0
	u.Email = strings.TrimSpace(u.Email) // Supprime : espaces au début + espaces à la fin + retours ligne+ tabulations EscapeString - protection XSS côté frontend
	u.Role = 1                           // Par défaut 0=admin 1 = normal
	u.PasswordReset = ""                 // Pas de Génération Générer une chaîne aléatoire de 40 caractères pour reset de password.
	u.PasswordResetExpire = time.Now()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) PrepareForUpdate() {
	u.Email = strings.TrimSpace(u.Email)
	u.UpdatedAt = time.Now()
}

func (u *User) NormalizeEmail() {
	u.Email = strings.TrimSpace(u.Email) // Supprime : espaces au début + espaces à la fin + retours ligne+ tabulations
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil

	default:
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}
		return nil
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	// GORM Hook
	// db.Create(&u)
	// GORM appelle automatiquement :  BeforeSave()  BeforeCreate()  INSERT AfterCreate()  AfterSave()

	err = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// Méthode attachée à User pour trouver un utilisateur par son email
func (u *User) FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	var err error // Déclare une variable err pour stocker les erreurs potentielles de Gorm

	// Requête Gorm pour récupérer l'utilisateur
	// - db.Model(User{}) : on indique qu'on travaille sur la table correspondant à User
	// - .Where("email = ?", email) : condition WHERE email = email fourni
	// - .Take(&u) : récupère le premier résultat correspondant et le stocke dans u
	// - .Error : récupère l'erreur éventuelle de la requête
	// Après cette ligne, u contiendra l'utilisateur trouvé (si trouvé) et err contiendra une erreur si la requête a échoué
	err = db.Model(User{}).Where("email = ?", email).Take(&u).Error
	// Si une erreur est survenue (problème de base de données)
	if err != nil {
		return &User{}, err // renvoie un User vide et l'erreur
	}

	// Cette ligne est censée vérifier si l'erreur était "record not found"
	// Mais ici elle ne sert à rien car on renvoie déjà err juste avant si err != nil
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}

	// Si tout s'est bien passé, on retourne l'utilisateur trouvé et nil pour l'erreur
	return u, err
}

// GetUserRoleByEmail récupère le rôle d'un utilisateur via son email
func (u *User) GetUserRoleByEmail(db *gorm.DB, email string) (uint8, error) {
	var user User // variable pour stocker l'utilisateur trouvé

	// Requête Gorm pour récupérer l'utilisateur correspondant à l'email
	err := db.Model(User{}).Where("email = ?", email).Take(&user).Error
	// Si une erreur survient
	if err != nil {
		if gorm.IsRecordNotFoundError(err) { // aucun utilisateur trouvé
			return 0, errors.New("User not found")
		}
		return 0, err // autre erreur (DB etc.)
	}

	// Retourne le rôle de l'utilisateur
	return user.Role, nil
}

// SetUserRoleByEmail met à jour le rôle d'un utilisateur via son email
func (u *User) SetUserRoleByEmail(db *gorm.DB, email string, role uint8) error {
	var user User // variable pour stocker l'utilisateur trouvé

	// Recherche de l'utilisateur par email
	err := db.Model(User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) { // aucun utilisateur trouvé
			return errors.New("User not found")
		}
		return err
	}

	// Empêche de changer le rôle de l'admin
	if user.ID == 1 {
		return errors.New("Cannot change role of admin")
	}

	// Met à jour le rôle dans l'instance et dans la base
	user.Role = role
	err = db.Model(&user).Where("email = ?", email).Update("role", role).Error
	if err != nil {
		return err
	}

	return nil
}

// Principe
// PasswordReset string `gorm:"size:10;unique" json:"password_reset"`
// Ce champ est utilisé comme token de réinitialisation de mot de passe.
// Typiquement :
//
//			L’utilisateur demande un reset.
//	   Tu génères un token aléatoire.
//	   Tu l’enregistres en base.
//	   Tu envoies un lien du type :
//				https://site.com/reset?token=ABCD1234...
//	   L’utilisateur clique.
//	   Tu retrouves l’utilisateur avec :FindUserByPasswordResetId(...)
func (u *User) FindUserByPasswordResetId(db *gorm.DB, passwordResetId string) (*User, error) {
	var err error
	err = db.Model(User{}).Where("password_reset = ?", passwordResetId).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

func (u *User) UpdateAUserAndRole(db *gorm.DB, uid uint32) (*User, error) {
	// To hash the password
	// Car Avec UpdateColumns() GORM bypass les hooks.
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Updating user with ID:", uid, "to email:", u.Email, "and role:", u.Role)

	// Ici, GORM ne prend pas u.ID pour déterminer quelle ligne mettre à jour, mais utilise Where("id = ?", uid).
	// Avec UpdateColumns() GORM bypass les hooks.
	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"email":      u.Email,
			"role":       u.Role,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is to display the updated user
	err = db.Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// Méthode avec receiver pointeur (*User).
func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {
	// To hash the password
	// Car Avec UpdateColumns() GORM bypass les hooks.
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":   u.Password,
			"email":      u.Email,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is to display the updated user
	err = db.Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) GeneratePasswordResetToken() {
	u.PasswordReset = utils.CreateRandString(40)
	u.PasswordResetExpire = time.Now().Add(time.Hour)
}

func RequestPasswordReset(db *gorm.DB, email string) (string, error) {
	user := User{}
	// On charge user sur basd email
	_, err := user.FindUserByEmail(db, email)
	if gorm.IsRecordNotFoundError(err) {
		return "", errors.New("Email doesn't exist in the server.")
	}

	user.GeneratePasswordResetToken()

	// Model(&User{}).Where("email = ?", email)
	// 	→ tu travailles sur le type, tu dois préciser manuellement quelle ligne tu veux mettre à jour.
	// Model(&user)
	// 	→ tu passes l’instance, GORM connaît déjà l’ID de la ligne et l’utilise automatiquement pour générer le WHERE id = ....

	db = db.Model(&user).
		UpdateColumns(map[string]interface{}{
			"password_reset":        user.PasswordReset,
			"password_reset_expire": time.Now().Add(time.Hour * time.Duration(1)), // Set new expire date to 1h
		})

	return user.PasswordReset, nil
}

func ResetPassword(db *gorm.DB, passwordResetId string, updatedPassword string) (*User, error, int) {
	user := User{}
	_, err := user.FindUserByPasswordResetId(db, passwordResetId)
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("This passwordResetId is invalid."), http.StatusNotFound
	}
	if user.PasswordResetExpire.Before(time.Now()) {
		return &User{}, errors.New("PasswordResetId has already expired."), http.StatusForbidden
	}

	user.Password = updatedPassword

	user.BeforeSave() /* This will hash the password */

	db = db.Model(&User{}).Where("password_reset = ?", passwordResetId).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":              user.Password,
			"updated_at":            time.Now(),
			"password_reset_expire": time.Now(), /* So it cannot be used a 2nd time */
		},
	)
	if db.Error != nil {
		return &User{}, db.Error, 0
	}

	fmt.Println("Update user passsword")

	return &user, nil, 0
}
