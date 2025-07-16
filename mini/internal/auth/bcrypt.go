package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns a bcrypt hash using a given cost
// Cost of 0 uses bcrypt.DefaultCost
func HashPassword(password string, cost int) (string, error) {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(h), err
}

// ComparePassword returns true if plaintext matches bcrypt hash
func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}