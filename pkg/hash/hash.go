package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
)

func GetRandomString(l int) string {
	validCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = validCharacters[RandInt(0, len(validCharacters)-1)]
	}
	return string(bytes)
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func GetHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}