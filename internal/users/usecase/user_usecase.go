package usecase

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/users"
	"regexp"
)

const passwordRequiredLength = "8"

type UserUseCaseRealisation struct {
	userDB   		users.UserRepository
}

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
		if err != nil {
			return false, "", err
		}
		if !isValid {
			return false, errorMessage, err
		}
	}

	return true, "", nil
}

func ValidateRegisterCredentials(userData models.User) (bool, string, error) {
	isNicknameValid, err := regexp.MatchString(`[a-z0-9_-]{3,15}`, userData.NickName)
	if err != nil {
		return false, "", err
	}
	if !isNicknameValid {
		return false, "The length of the name must be from 2 to 15 characters and must not contain spaces," +
			" special characters and numbers", nil
	}

	isEmailValid, err := regexp.MatchString(`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`, userData.Email)
	if err != nil {
		return false, "", err
	}
	if !isEmailValid {
		return false, "Invalid email", nil
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

func (userR UserUseCaseRealisation) Register(userData models.User) (string, string, error) {
	isValidCredentials, msg, err := ValidateRegisterCredentials(userData)
	if err != nil {
		return "", "", err
	}
	if !isValidCredentials {
		return "", msg, nil
	}

	isEmailUnique, err := userR.userDB.IsEmailUnique(userData.Email)
	if err != nil {
		return "", "", err
	}
	if !isEmailUnique {
		return "", "Email is already taken", nil
	}

	isNicknameUnique, err := userR.userDB.IsNicknameUnique(userData.NickName)
	if err != nil {
		return "", "", err
	}
	if !isNicknameUnique {
		return "", "Nickname is already taken", nil
	}

	sessionToken, err := userR.userDB.CreateUser(userData)
	if err != nil {
		return "", "", err
	}

	return sessionToken, "", nil
}

func NewUserUserCaseRealization(userDB users.UserRepository) UserUseCaseRealisation {
	return UserUseCaseRealisation{
		userDB: userDB,
	}
}