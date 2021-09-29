package utils

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetHash(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}
