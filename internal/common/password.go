package common

import "golang.org/x/crypto/bcrypt"

// HashPassword creates a salted bcrypt hash. Store only the returned hash;
// callers must never persist the plaintext password.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword verifies a plaintext password against its stored bcrypt hash.
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
