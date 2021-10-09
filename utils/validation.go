package utils

import (
	"2021_2_LostPointer/models"
	"regexp"
)

const passwordRequiredLength = "8"
const minNicknameLength = "3"
const maxNicknameLength = "15"

func ValidatePassword(password string) (bool, string, error) {
	patterns := map[string]string {
		`^.{` + passwordRequiredLength + `,}$`: "Password must contain at least" + passwordRequiredLength + "characters",
		`[0-9]`: "Password must contain at least one digit",
		`[A-Z]`: "Password must contain at least one uppercase letter",
		`[a-z]`: "Password must contain at least one lowercase letter",
		`[\@\ \!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\?\[\\\]\^\_]`: "Password must contain as least one special character",

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

func ValidateSignUp(user *models.User) (bool, string, error) {
	nickNameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{` + minNicknameLength + `,` + maxNicknameLength + `}$`, user.Nickname)
	if err != nil {
		return false, "", err
	}
	if !nickNameValid {
		return false, "The length of the name must be from " + minNicknameLength + " to " + maxNicknameLength + " characters", nil
	}

	usernameValid, err := regexp.MatchString(`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`, user.Email)
	if err != nil {
		return false, "", err
	}
	if !usernameValid {
		return false, "Invalid email", nil
	}

	passwordValid, message, err := ValidatePassword(user.Password)
	if err != nil {
		return false, "", err
	}
	if !passwordValid {
		return false, message, nil
	}

	return true, "", nil
}
