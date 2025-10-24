package password

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func Equal(password string, targetHash string) bool {
	byteHash := []byte(targetHash)
	bytePassword := []byte(password)

	result := bcrypt.CompareHashAndPassword(byteHash, bytePassword)

	return result == nil
}
