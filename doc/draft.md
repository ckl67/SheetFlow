au niveau du hook de gorm V2

Create n'a pas été modifié
err = db.Create(&u).Error
fonctionne donc toujours

OK pris note que BeforeSave demande une adapatation, même si je n'utilise pas tx chez moi

func (u *User) BeforeSave(tx *gorm.DB) error {
hashedPassword, err := Hash(u.Password)
if err != nil {
return err
}
u.Password = string(hashedPassword)
return nil
}
