package utils

import (
	"2021_2_LostPointer/models"
	"log"
	"regexp"
)

const passwordRequiredLength = "8"

func ValidatePassword(password string) (bool, string, error) {
	patterns := map[string]string {
		`^.{` + passwordRequiredLength + `,}$`: "Password must contain at least" + passwordRequiredLength + "characters",
		`[0-9]`: "Password must contain at least one digit",
		`[A-Z]`: "Password must contain at least one uppercase letter",
		`[a-z]`: "Password must contain at least one lowercase letter",
		`[\@\ \!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\?\[\\\]\^\_]`: "Password must contain as least one special symbol",

	}

	for pattern, errorMessage := range patterns {
		isValid, err := regexp.MatchString(pattern, password)
		log.Println(pattern, isValid)
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
	nameValid, err := regexp.MatchString(`^([a-zA-Z]{2,15})$`, user.Name)
	if err != nil {
		return false, "", err
	}
	if !nameValid {
		return false, "The length of the name must be from 2 to 15 characters and must not contain spaces," +
			" special characters and numbers", nil
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
