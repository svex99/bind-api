package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/svex99/bind-api/utils/token"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *User) SaveUser() (*User, error) {
	if err := DB.Create(&u).Error; err != nil {
		return &User{}, err
	}

	return u, nil
}

func LoginUser(email, password string) (string, error) {
	user := User{}

	err := DB.Model(User{}).Where("email = ?", email).Take(&user).Error

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	newToken, err := token.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return newToken, nil
}
