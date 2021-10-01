package utils

import (
	"2021_2_LostPointer/models"
	"regexp"
)

func passwordValidation(password string) (bool, error) {
	isLong := len(password) > 8
	containDigit, err := regexp.MatchString(`[0-9]`, password)
	if err != nil {
		return false, err
	}
	containLower, err := regexp.MatchString(`[a-z]`, password)
	if err != nil {
		return false, err
	}
	containUpper, err := regexp.MatchString(`[A-Z]`, password)
	if err != nil {
		return false, err
	}
	containSpecial, err := regexp.MatchString(`[\!\@\#\$\%\^\&\*]`, password)
	if err != nil {
		return false, err
	}
	return isLong && containDigit && containLower && containUpper && containSpecial, nil
}

func ValidateSignUp(user *models.User) (bool, error) {
	nameValid, err := regexp.MatchString(`^([a-zA-Z]{2,15})$`, user.Name)
	if err != nil {
		return false, err
	}
	usernameValid, err := regexp.MatchString(`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`, user.Username)
	if err != nil {
		return false, err
	}
	passwordValid, err := passwordValidation(user.Password)
	if err != nil {
		return false, err
	}

	return nameValid && usernameValid && passwordValid, nil
}
