package validation

import (
	"2021_2_LostPointer/internal/constants"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: дописать)

func TestValidateRegisterCredentials(t *testing.T) {
	tests := []struct {
		name           string
		expectedValid        bool
		email          string
		password       string
		nickname       string
		expectedInvalidMessage string
		expectedError  bool
	}{
		{
			name: "invalid nickname",
			email: "qwerty@qw.com",
			password: "Qwerty1111",
			nickname: "Qwerty!",
			expectedInvalidMessage: constants.InvalidNicknameMessage,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid, errMsg, err := ValidateRegisterCredentials(test.email, test.password, test.nickname)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expectedValid, isValid)
				assert.Equal(t, test.expectedInvalidMessage, errMsg)
				assert.NoError(t, err)
			}
		})
	}
}
