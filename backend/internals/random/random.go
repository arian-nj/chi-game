package random

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

func GenerateRandomUsername(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	for i := 0; i < length; i++ {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}

	return string(randomBytes)
}

func GenerateRandomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes)[:length]
}

func GenerateRandomNumber(max_num int) int {
	max := big.NewInt(int64(max_num))
	random_index, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}
	return int(random_index.Int64())

}
