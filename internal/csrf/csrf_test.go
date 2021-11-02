package csrf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashToken_Check(t *testing.T) {
	hashToken, _ := NewHMACHashToken("secret")

	const cookie = "12345"
	tests := []struct {
		name          string
		inputToken    string
		expectedValid bool
		expectedError bool
	}{
		{
			name: "valid token",
			inputToken: func() string {
				token, _ := hashToken.Create(cookie, int64(2000000000000))
				return token
			}(),
			expectedValid: true,
		},
		{
			name:          "bad token",
			inputToken:    "11",
			expectedError: true,
		},
		{
			name:          "expired token",
			inputToken:    "1:1",
			expectedError: true,
		},
		{
			name:          "token contains letters",
			inputToken:    "1:c1",
			expectedError: true,
		},
		{
			name:          "token contains negative time",
			inputToken:    "1:-2",
			expectedError: true,
		},
		{
			name:          "failed hex decode token",
			inputToken:    "ad1wd3wa1eaw2ef d98n12:30000000000",
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid, err := hashToken.Check(cookie, test.inputToken)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expectedValid, isValid)
				assert.NoError(t, err)
			}
		})
	}
}
