package random

import (
	"crypto/rand"
)

func randomChars(length int, charset string) string {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}
	for i := range length {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}
	return string(randomBytes)
}

func GenerateInviteCode(length int) string {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	return randomChars(length, charset)
}

func GenerateRandomUsername(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	return randomChars(length, charset)
}
