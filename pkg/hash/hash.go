package hash

import (
	"golang.org/x/crypto/bcrypt"
)

const defaultCost = 12

// Hash creates a bcrypt hash from password
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
	return string(bytes), err
}

// Verify checks if password matches hash
func Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
