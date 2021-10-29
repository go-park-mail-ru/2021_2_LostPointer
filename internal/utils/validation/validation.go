package validation

import (
	"2021_2_LostPointer/internal/models"
	"regexp"
)

const passwordRequiredLength = "8"
const minNicknameLength = "3"
const maxNicknameLength = "15"

const PasswordValidationInvalidLengthMessage = "Password must contain at least " + passwordRequiredLength + " characters"
const PasswordValidationNoDigitMessage = "Password must contain at least one digit"
const PasswordValidationNoUppercaseMessage = "Password must contain at least one uppercase letter"
const PasswordValidationNoLowerCaseMessage = "Password must contain at least one lowercase letter"
const PasswordValidationNoSpecialSymbolMessage = "Password must contain as least one special character"
const NickNameValidationInvalidLengthMessage = "The length of nickname must be from " + minNicknameLength + " to " + maxNicknameLength + " characters"
const InvalidEmailMessage = "Invalid email"

const EmailRegexPattern = `[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`

func ValidatePassword(password string) (bool, string, error) {
	patterns := map[string]string {
		`^.{` + passwordRequiredLength + `,}$`: PasswordValidationInvalidLengthMessage,
		`[0-9]`:                                PasswordValidationNoDigitMessage,
		`[A-Z]`:                                PasswordValidationNoUppercaseMessage,
		`[a-z]`:                                PasswordValidationNoLowerCaseMessage,
		`[\@\ \!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\?\[\\\]\^\_]`: PasswordValidationNoSpecialSymbolMessage,

	}

	for pattern, errorMessage := range patterns {
		isValid, err := regexp.MatchString(pattern, password)
		if err != nil {
			return false, "", err
		}
		if !isValid {
			return false, errorMessage, err
		}
	}

	return true, "", nil
}

func ValidateRegisterCredentials(userData *models.User) (bool, string, error) {
	isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{` +minNicknameLength+ `,` +maxNicknameLength+ `}$`, userData.Nickname)
	if err != nil {
		return false, "", err
	}
	if !isNicknameValid {
		return false, NickNameValidationInvalidLengthMessage, nil
	}

	isEmailValid, err := regexp.MatchString(EmailRegexPattern, userData.Email)
	if err != nil {
		return false, "", err
	}
	if !isEmailValid {
		return false, InvalidEmailMessage, nil
	}

	passwordValid, message, err := ValidatePassword(userData.Password)
	if err != nil {
		return false, "", err
	}
	if !passwordValid {
		return false, message, nil
	}

	return true, "", nil
}
