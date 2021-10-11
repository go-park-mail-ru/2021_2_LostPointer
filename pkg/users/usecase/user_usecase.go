package usecase

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/users"
	"math/rand"
	"regexp"
	"time"
)


const SessionTokenLength = 40
const passwordRequiredLength = "8"
const minNicknameLength = "3"
const maxNicknameLength = "15"

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

func GetRandomString(l int) string {
	validCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = validCharacters[RandInt(0, len(validCharacters) - 1)]
	}
	return string(bytes)
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

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

func ValidateRegisterCredentials(userData models.User) (bool, string, error) {
	isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{` + minNicknameLength + `,` + maxNicknameLength + `}$`, userData.Nickname)
	if err != nil {
		return false, "", err
	}
	if !isNicknameValid {
		return false, "The length of nickname must be from " + minNicknameLength + " to " + maxNicknameLength + " characters", nil
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
	sessionToken := GetRandomString(SessionTokenLength)
	sessionToken, err = userR.redisStore.StoreSession(userID)
	if err != nil {
		return "", "", err
	}

	return sessionToken, "", nil
}

func (userR UserUseCase) Login(authData models.Auth) (string, error) {
	userID, err := userR.userDB.UserExits(authData)
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