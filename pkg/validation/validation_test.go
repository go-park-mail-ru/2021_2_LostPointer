package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"2021_2_LostPointer/internal/constants"
)

//nolint:scopelint
func TestValidateRegisterCredentials(t *testing.T) {
	tests := []struct {
		name                   string
		expectedValid          bool
		email                  string
		password               string
		nickname               string
		expectedInvalidMessage string
		expectedError          bool
	}{
		{
			name:                   "invalid nickname",
			email:                  "qwerty@qw.com",
			password:               "Qwerty1111",
			nickname:               "Qwerty!",
			expectedInvalidMessage: constants.NicknameInvalidSyntaxMessage,
		},
	}

	for _, current := range tests {
		t.Run(current.name, func(t *testing.T) {
			isValid, errMsg, err := ValidateRegisterCredentials(current.email, current.password, current.nickname)
			if current.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, current.expectedValid, isValid)
				assert.Equal(t, current.expectedInvalidMessage, errMsg)
				assert.NoError(t, err)
			}
		})
	}
}
