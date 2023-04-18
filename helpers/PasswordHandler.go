package helpers

import "golang.org/x/crypto/bcrypt"

func HashPass(value string) string {
	var (
		salt     = 8
		password = []byte(value)
	)

	hash, _ := bcrypt.GenerateFromPassword(password, salt)

	return string(hash)
}

func CompareHashAndPassword(hash, pass []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, pass)

	return err == nil
}
