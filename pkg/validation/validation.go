package validation

import (
	"regexp"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
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

func ValidateRegisterCredentials(registerData *models.RegisterData) (bool, string, error) {
	isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{`+constants.MinNicknameLength+`,`+constants.MaxNicknameLength+`}$`, registerData.Nickname)
	if err != nil {
		return false, "", err
	}
	if !isNicknameValid {
		return false, constants.InvalidNicknameMessage, nil
	}

	isEmailValid, err := regexp.MatchString(constants.EmailRegexPattern, registerData.Email)
	if err != nil {
		return false, "", err
	}
	if !isEmailValid {
		return false, constants.InvalidEmailMessage, nil
	}

	passwordValid, message, err := ValidatePassword(registerData.Password)
	if err != nil {
		return false, "", err
	}
	if !passwordValid {
		return false, message, nil
	}

	return true, "", nil
}
