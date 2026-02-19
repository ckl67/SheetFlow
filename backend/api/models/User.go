package models

import (
	"backend/api/config"
	"backend/api/utils"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                  uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Email               string    `gorm:"size:100;not null;unique" json:"email"`
	Role                uint8     `json:"role"` // 0=admin 1=normal
	Password            string    `gorm:"size:100;not null;" json:"password"`
	PasswordReset       string    `gorm:"size:64;" json:"password_reset"`
	PasswordResetExpire time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"password_reset_expire"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//
// ========================
// PASSWORD
// ========================
//

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

//
// ========================
// GORM HOOK
// ========================
//

func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password == "" {
		return nil
	}

	// Si déjà hashé
	if strings.HasPrefix(u.Password, "$2") {
		return nil
	}

	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	return nil
}

//
// ========================
// PREPARE
// ========================
//

func (u *User) PrepareForCreate() {
	u.ID = 0
	u.Email = utils.SanitizeUserEmail(u.Email)
	u.Role = 1
	u.PasswordReset = ""
	u.PasswordResetExpire = time.Now()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) PrepareForUpdate() {
	u.Email = utils.SanitizeUserEmail(u.Email)
	u.UpdatedAt = time.Now()
}

//
// ========================
// CRUD
// ========================
//

func (u *User) Save(db *gorm.DB) (*User, error) {
	if err := db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) FindByID(db *gorm.DB, id uint32) (*User, error) {
	if err := db.First(u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return u, nil
}

func (u *User) FindByEmail(db *gorm.DB, email string) (*User, error) {
	email = utils.SanitizeUserEmail(email)

	if err := db.Where("email = ?", email).First(u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return u, nil
}

func (u *User) Update(db *gorm.DB, id uint32, admin uint32) (*User, error) {
	u.UpdatedAt = time.Now()

	if admin == config.ADMIN_UID {
		if err := db.Model(&User{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"email":      u.Email,
				"password":   u.Password,
				"role":       u.Role,
				"updated_at": u.UpdatedAt,
			}).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Model(&User{}).
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"email":    u.Email,
				"password": u.Password,
				//	"role":       u.Role,
				"updated_at": u.UpdatedAt,
			}).Error; err != nil {
			return nil, err
		}
	}

	return u.FindByID(db, id)
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//
// ========================
// ROLE MANAGEMENT
// ========================
//

func SetUserRoleByEmail(db *gorm.DB, email string, role uint8) error {
	email = utils.SanitizeUserEmail(email)

	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if user.ID == 1 {
		return errors.New("cannot change role of admin")
	}

	return db.Model(&user).Update("role", role).Error
}

//
// ========================
// PASSWORD RESET
// ========================
//

func (u *User) GeneratePasswordResetToken() {
	u.PasswordReset = utils.CreateRandString(40)
	u.PasswordResetExpire = time.Now().Add(time.Hour)
}

func RequestPasswordReset(db *gorm.DB, email string) (string, error) {
	email = utils.SanitizeUserEmail(email)

	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("email not found")
	}

	user.GeneratePasswordResetToken()

	if err := db.Model(&user).Updates(map[string]interface{}{
		"password_reset":        user.PasswordReset,
		"password_reset_expire": user.PasswordResetExpire,
	}).Error; err != nil {
		return "", err
	}

	return user.PasswordReset, nil
}

func ResetPassword(db *gorm.DB, token string, newPassword string) (*User, error, int) {
	var user User

	if err := db.Where("password_reset = ?", token).First(&user).Error; err != nil {
		return nil, errors.New("invalid reset token"), http.StatusNotFound
	}

	if user.PasswordResetExpire.Before(time.Now()) {
		return nil, errors.New("reset token expired"), http.StatusForbidden
	}

	user.Password = newPassword
	user.PasswordResetExpire = time.Now()
	user.PasswordReset = ""

	if err := db.Save(&user).Error; err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &user, nil, 0
}

func (u *User) NormalizeEmail() {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
}

func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Limit(100).Find(&users).Error
	return users, err
}
