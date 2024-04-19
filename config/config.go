package config

import "golang.org/x/crypto/bcrypt"

// CreateToken is a function that takes a password and returns a hashed version of it
// It uses the bcrypt.GenerateFromPassword function to hash the password and returns the hashed password and an error if there is one.
func CreateToken(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// ComparePasswords is a function that takes a hashed password and a password and returns an error
// if the password does not match the hashed password or nil if the password matches the hashed password.
func ComparePasswords(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}