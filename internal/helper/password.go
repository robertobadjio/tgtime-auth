package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword ???
func VerifyPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}
