package validation

import (
	"regexp"

	"github.com/asaskevich/govalidator"

	"2021_2_LostPointer/internal/constants"
)

func ValidatePassword(password string) (bool, string, error) {
	patterns := map[string]string{
		`^.{` + constants.PasswordRequiredLength + `,}$`: constants.PasswordValidationInvalidLengthMessage,
		`[0-9]`: constants.PasswordValidationNoDigitMessage,
		`[A-Z]`: constants.PasswordValidationNoUppercaseMessage,
		`[a-z]`: constants.PasswordValidationNoLowerCaseMessage,
		`[\@\ \!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\?\[\\\]\^\_]`: constants.PasswordValidationNoSpecialSymbolMessage,
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

func ValidateRegisterCredentials(email string, password string, nickname string) (bool, string, error) {
	isEmailValid := govalidator.IsEmail(email)
	if !isEmailValid {
		return false, constants.InvalidEmailMessage, nil
	}
	isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{`+constants.MinNicknameLength+`,`+constants.MaxNicknameLength+`}$`, nickname)
	if err != nil {
		return false, "", err
	}
	if !isNicknameValid {
		return false, constants.InvalidNicknameMessage, nil
	}

	passwordValid, message, err := ValidatePassword(password)
	if err != nil {
		return false, "", err
	}
	if !passwordValid {
		return false, message, nil
	}

	return true, "", nil
}
