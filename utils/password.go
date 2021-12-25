package utils

import (
	"gitlab.com/canco1/canco-kit/common_error"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", common_error.ErrCanNotHashPassword
	}
	return string(hashedPassword), nil
}

// VerifyPassword checks if the provided password is correct or not
func VerifyPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
