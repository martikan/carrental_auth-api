package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	PasswordUtils passwordUtilsInterface = &passwordUtils{}
)

type passwordUtils struct{}

type passwordUtilsInterface interface {
	HashPassword(string) (string, error)
	CheckPassword(string, string) error
}

// HashPassword returns the bcrypt hash of the password
func (p *passwordUtils) HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// CheckPassword check if the provided password is correct or not
func (p *passwordUtils) CheckPassword(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
