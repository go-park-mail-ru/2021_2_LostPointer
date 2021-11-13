package utils

import (
	"crypto/rand"
	"math/big"
)

func GetRandomString(l int) (string, error) {
	validCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lengthValidCharacters := int64(len(validCharacters))
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(lengthValidCharacters-1))
		if err != nil {
			return "", err
		}
		bytes[i] = validCharacters[randomIndex.Int64()]
	}
	return string(bytes), nil
}
