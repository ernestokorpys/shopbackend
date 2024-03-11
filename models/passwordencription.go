package models

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func (user *User) EncryptedPassword(password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Println("Error hashing password:", err)
		return
	}
	user.Password = hashedPassword
}

func (user *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
}
