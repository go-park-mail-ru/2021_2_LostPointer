package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"2021_2_LostPointer/internal/constants"
)

type User struct {
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
	Nickname string `json:"nickname" form:"nickname" query:"nickname"`
}

//nolint:scopelint
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
			expectedErrorMsg: constants.PasswordInvalidLengthMessage,
		},
		{
			name:             "no digit in password",
			password:         "qQw!eeQW!",
			expectedErrorMsg: constants.PasswordNoDigitMessage,
		},
		{
			name:             "no letters in password",
			password:         "12312123123!",
			expectedErrorMsg: constants.PasswordNoLetterMessage,
		},
		{
			name:             "correct password",
			password:         "QwertyY1#",
			expectedValid:    true,
			expectedErrorMsg: "",
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			isValid, errMsg, err := ValidatePassword(currentTest.password)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, currentTest.expectedValid, isValid)
				assert.Equal(t, currentTest.expectedErrorMsg, errMsg)
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateRegisterCredentials(t *testing.T) {
	tests := []struct {
		name             string
		userData         User
		expectedValid    bool
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			name: "valid credentials",
			userData: User{
				Email:    "test@test.com",
				Password: "Qwerty123$",
				Nickname: "Kucheryavenko",
			},
			expectedValid:    true,
			expectedErrorMsg: "",
		},
		{
			name: "wrong name",
			userData: User{
				Email:    "test@test.com",
				Password: "Qwerty123$",
				Nickname: "1321!#23@#E123123",
			},
			expectedErrorMsg: constants.NicknameInvalidSyntaxMessage,
		},
		{
			name: "name too short",
			userData: User{
				Email:    "test@test.com",
				Password: "Qwerty123$",
				Nickname: "fa",
			},
			expectedErrorMsg: constants.NicknameInvalidLengthMessage,
		},
		{
			name: "name too long",
			userData: User{
				Email:    "test@test.com",
				Password: "Qwerty123$",
				Nickname: "faawdaaecsdefvsrvsfvsfgbdfbg",
			},
			expectedErrorMsg: constants.NicknameInvalidLengthMessage,
		},
		{
			name: "no @ in email",
			userData: User{
				Email:    "testtest.com",
				Password: "Qwerty123$",
				Nickname: "test",
			},
			expectedErrorMsg: constants.EmailInvalidSyntaxMessage,
		},
		{
			name: "no domain in email",
			userData: User{
				Email:    "test@test",
				Password: "Qwerty123$",
				Nickname: "test",
			},
			expectedErrorMsg: constants.EmailInvalidSyntaxMessage,
		},
		{
			name: "no domain in email",
			userData: User{
				Email:    "test@test",
				Password: "Qwerty123$",
				Nickname: "test",
			},
			expectedErrorMsg: constants.EmailInvalidSyntaxMessage,
		},
		{
			name: "wrong password",
			userData: User{
				Email:    "test@test.ru",
				Password: "Qwert1",
				Nickname: "test",
			},
			expectedErrorMsg: constants.PasswordInvalidLengthMessage,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			isValid, errMsg, err := ValidateRegisterCredentials(currentTest.userData.Email, currentTest.userData.Password, currentTest.userData.Nickname)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, currentTest.expectedValid, isValid)
				assert.Equal(t, currentTest.expectedErrorMsg, errMsg)
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePlaylistTitle(t *testing.T) {
	tests := []struct {
		name                 string
		title                string
		expected             bool
		expectedErrorMessage string
		expectedError        bool
	}{
		{
			name:     "valid title length",
			title:    "test",
			expected: true,
		},
		{
			name:                 "invalid title length - less then 3 symbols",
			title:                "t",
			expectedErrorMessage: constants.PlaylistTitleInvalidLengthMessage,
		},
		{
			name:                 "invalid title length - more then 30 symbols",
			title:                "tttttttttttttttttttttttttttttttttttttttt",
			expectedErrorMessage: constants.PlaylistTitleInvalidLengthMessage,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			isValid, errMsg, err := ValidatePlaylistTitle(currentTest.title)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.Equal(t, currentTest.expected, isValid)
				assert.Equal(t, currentTest.expectedErrorMessage, errMsg)
				assert.NoError(t, err)
			}
		})
	}
}
