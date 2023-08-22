package auth

import "golang.org/x/crypto/bcrypt"

// Encrypt encrypts string in hashed.
func Encrypt(passwd string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

// Compare compares hashedPassword store in db with password from user request.
func Compare(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
