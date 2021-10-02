package utils

import (
	"2021_2_LostPointer/models"
	"regexp"
)

func validatePassword(password string) (bool, string, error) {
	isLong := len(password) > 8
	if !isLong {
		return false, "Password must contain at least 8 characters", nil
	}

	containDigit, err := regexp.MatchString(`[0-9]`, password)
	if err != nil {
		return false, "", err
	}
	if !containDigit {
		return false, "Password must contain at least one digit", nil
	}

	containUpper, err := regexp.MatchString(`[A-Z]`, password)
	if err != nil {
		return false, "", err
	}
	if !containUpper {
		return false, "Password must contain at least one uppercase letter", nil
	}

	containLower, err := regexp.MatchString(`[a-z]`, password)
	if err != nil {
		return false, "", err
	}
	if !containLower {
		return false, "Password must contain at least one lowercase letter", nil
	}

	containSpecial, err := regexp.MatchString(`[\!\@\#\$\%\^\&\*]`, password)
	if err != nil {
		return false, "", err
	}
	if !containSpecial {
		return false, "Password must contain as least one special symbol", nil
	}

	return true, "", nil
}

func ValidateSignUp(user *models.User) (bool, string, error) {
	nameValid, err := regexp.MatchString(`^([a-zA-Z]{2,15})$`, user.Name)
	if err != nil {
		return false, "", err
	}
	if !nameValid {
		return false, "The length of the name must be from 2 to 15 characters and must not contain spaces, special characters and numbers", nil
	}

	usernameValid, err := regexp.MatchString(`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`, user.Username)
	if err != nil {
		return false, "", err
	}
	if !usernameValid {
		return false, `Invalid email`, nil
	}

	passwordValid, message, err := validatePassword(user.Password)
	if err != nil {
		return false, "", err
	}
	if !passwordValid {
		return false, message, nil
	}

	return true, "", nil
}
