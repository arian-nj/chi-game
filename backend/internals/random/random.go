package random

import (
	"crypto/rand"
)

func GenerateRandomUsername(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}

	for i := range length {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}

	return string(randomBytes)
}
