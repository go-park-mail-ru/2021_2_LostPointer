package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type HashToken struct {
	Secret []byte
}

var Tokens *HashToken

//nolint:gochecknoinits
func init() {
	Tokens = NewHMACHashToken(os.Getenv("CSRF_SECRET"))
}

func NewHMACHashToken(secret string) *HashToken {
	return &HashToken{Secret: []byte(secret)}
}

func (tk *HashToken) Create(cookie string, tokenExpTime int64) (string, error) {
	h := hmac.New(sha256.New, tk.Secret)
	data := fmt.Sprintf("%s:%d", cookie, tokenExpTime)
	if _, err := h.Write([]byte(data)); err != nil {
		return "", err
	}
	token := hex.EncodeToString(h.Sum(nil)) + ":" + strconv.FormatInt(tokenExpTime, 10)
	return token, nil
}

func (tk *HashToken) Check(cookie string, inputToken string) (bool, error) {
	tokenData := strings.Split(inputToken, ":")
	if len(tokenData) != 2 {
		return false, fmt.Errorf("bad token data")
	}

	tokenExp, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("error token exp operation")
	}

	if tokenExp < time.Now().Unix() {
		return false, errors.New("Session expires")
	}

	hash := hmac.New(sha256.New, tk.Secret)
	data := fmt.Sprintf("%s:%d", cookie, tokenExp)
	_, err = hash.Write([]byte(data))
	if err != nil {
		return false, err
	}
	expectedMAC := hash.Sum(nil)
	messageMAC, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, fmt.Errorf("can't hex decode token")
	}

	return hmac.Equal(messageMAC, expectedMAC), nil
}
