package userRepository

import (
	"bytes"
	"golang.org/x/crypto/argon2"
)

func hashPass(salt []byte, plainPassword string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	return append(salt, hashedPass...)
}

func checkPass(passHash []byte, plainPassword string) bool {
	salt := make([]byte, 8)
	copy(salt, passHash[:8])
	userPassHash := hashPass(salt, plainPassword)
	return bytes.Equal(userPassHash, passHash)
}
