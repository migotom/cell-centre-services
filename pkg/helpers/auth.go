package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes password by bcrypt.
func HashPassword(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash)
}
