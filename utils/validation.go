package utils

import (
	"2021_2_LostPointer/models"
	"regexp"
)

func ValidateSignUp(user *models.User) (bool, error) {
	nameValid, err := regexp.MatchString(`^([a-zA-Z]{2,15})$`, user.Name)
	/*log.Println("Name:", nameValid)*/
	if err != nil {
		return false, err
	}
	usernameValid, err := regexp.MatchString(`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`, user.Username)
	/*log.Println("Username", usernameValid)*/

	return nameValid && usernameValid, nil
}
