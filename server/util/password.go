package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GeneratePassword() string {
	var passwordLength = 10
	password := make([]byte, passwordLength)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	println(string(password))
	return string(password)
}
