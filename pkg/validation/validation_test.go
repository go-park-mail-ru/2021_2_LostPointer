package validation

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name             string
		password         string
		expectedValid    bool
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name:             "wrong password length",
			password:         "q1!Q",
			expectedErrorMsg: constants.PasswordValidationInvalidLengthMessage,
		},
		{
			name:             "no digit in password",
			password:         "qQw!eeQW!",
			expectedErrorMsg: constants.PasswordValidationNoDigitMessage,
		},
		{
			name:             "no uppercase in password",
			password:         "q2w!ee11!",
			expectedErrorMsg: constants.PasswordValidationNoUppercaseMessage,
		},
		{
			name:             "no lowercase in password",
			password:         "Q2W!EE11!",
			expectedErrorMsg: constants.PasswordValidationNoLowerCaseMessage,
		},
		{
			name:             "no special symbols in password",
			password:         "q123Wew2eq",
			expectedErrorMsg: constants.PasswordValidationNoSpecialSymbolMessage,
		},
		{
			name:             "correct password",
			password:         "QwertyY1#",
			expectedValid:    true,
			expectedErrorMsg: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid, errMsg, err := ValidatePassword(test.password)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expectedValid, isValid)
				assert.Equal(t, test.expectedErrorMsg, errMsg)
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRegisterCredentials(t *testing.T) {
	tests := []struct {
		name             string
		userData         models.User
		expectedValid    bool
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "valid credentials",
			userData: models.User{
				Email:    "test@test.com",
				Password: "Qwerty123$",
				Nickname: "Kucheryavenko",
			},
			expectedValid:    true,
			expectedErrorMsg: "",
		},
		{
			name: "name too short",
			userData: models.User{
				Email:    "test@test.com",
				Password: "Qwerty123$",
				Nickname: "fa",
			},
			expectedErrorMsg: constants.InvalidNicknameMessage,
		},
		{
			name: "name too long",
			userData: models.User{
				Email:    "test@test.com",
				Password: "Qwerty123$",
				Nickname: "faawdaaecsdefvsrvsfvsfgbdfbg",
			},
			expectedErrorMsg: constants.InvalidNicknameMessage,
		},
		{
			name: "no @ in email",
			userData: models.User{
				Email:    "testtest.com",
				Password: "Qwerty123$",
				Nickname: "test",
			},
			expectedErrorMsg: constants.InvalidEmailMessage,
		},
		{
			name: "no domain in email",
			userData: models.User{
				Email:    "test@test",
				Password: "Qwerty123$",
				Nickname: "test",
			},
			expectedErrorMsg: constants.InvalidEmailMessage,
		},
		{
			name: "no domain in email",
			userData: models.User{
				Email:    "test@test",
				Password: "Qwerty123$",
				Nickname: "test",
			},
			expectedErrorMsg: constants.InvalidEmailMessage,
		},
		{
			name: "wrong password",
			userData: models.User{
				Email:    "test@test",
				Password: "Qwerty",
				Nickname: "test",
			},
			expectedErrorMsg: constants.InvalidEmailMessage,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid, errMsg, err := ValidateRegisterCredentials(&test.userData)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.expectedValid, isValid)
				assert.Equal(t, test.expectedErrorMsg, errMsg)
				assert.NoError(t, err)
			}
		})
	}
}
