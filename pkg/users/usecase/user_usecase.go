package usecase

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/users"
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

type UserUseCase struct {
	userDB	   users.UserRepositoryIFace
	redisStore users.RedisStoreIFace
}

func NewUserUserCase(userDB users.UserRepositoryIFace, redisStore users.RedisStoreIFace) UserUseCase {
	return UserUseCase{
		userDB: userDB,
		redisStore: redisStore,
	}
}

func ValidatePassword(password string) (bool, string, error) {
	patterns := map[string]string {
		`^.{` + passwordRequiredLength + `,}$`: PasswordValidationInvalidLengthMessage,
		`[0-9]`: PasswordValidationNoDigitMessage,
		`[A-Z]`: PasswordValidationNoUppercaseMessage,
		`[a-z]`: PasswordValidationNoLowerCaseMessage,
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

func ValidateRegisterCredentials(userData models.User) (bool, string, error) {
	isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{` + minNicknameLength + `,` + maxNicknameLength + `}$`, userData.Nickname)
	if err != nil {
		return false, "", err
	}
	if !isNicknameValid {
		return false, NickNameValidationInvalidLengthMessage, nil
	}

	isEmailValid, err := regexp.MatchString(`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`, userData.Email)
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

func (userR UserUseCase) Register(userData models.User) (string, string, error) {
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

	isNicknameUnique, err := userR.userDB.IsNicknameUnique(userData.Nickname)
	if err != nil {
		return "", "", err
	}
	if !isNicknameUnique {
		return "", "Nickname is already taken", nil
	}

	userID, err := userR.userDB.CreateUser(userData)
	sessionToken, err := userR.redisStore.StoreSession(userID)
	if err != nil {
		return "", "", err
	}

	return sessionToken, "", nil
}

func (userR UserUseCase) Login(authData models.Auth) (string, error) {
	userID, err := userR.userDB.DoesUserExist(authData)
	if err != nil {
		return "", err
	}
	if userID == 0 {
		return "", nil
	}

	sessionToken, err := userR.redisStore.StoreSession(userID)
	if err != nil {
		return "", err
	}

	return sessionToken, nil
}

func (userR UserUseCase) Logout(cookieValue string) {
	userR.redisStore.DeleteSession(cookieValue)
}

func (userR UserUseCase) IsAuthorized(cookieValue string) (bool, error) {
	id, err := userR.redisStore.GetSessionUserId(cookieValue)
	if err != nil {
		return false, err
	}
	if id == 0 {
		return false, nil
	}
	return true, nil
}

func (userR UserUseCase) GetSettings(cookieValue string) (*models.Settings, error) {
	userID, err := userR.redisStore.GetSessionUserId(cookieValue)
	if err != nil {
		return nil, err
	}

	settings, err := userR.userDB.GetSettings(userID)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
