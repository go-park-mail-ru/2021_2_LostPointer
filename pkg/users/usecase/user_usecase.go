package usecase

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/users"
	"log"
	"mime/multipart"
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
const NotUniqueEmailMessage = "Email is not unique"
const NotUniqueNicknameMessage = "Nickname is not unique"
const WrongCredentialsMessage = "Wrong email or password"

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

func (userR UserUseCase) Register(userData models.User) (string, *models.CustomError) {
	// 1) Проверка email и nickname на уникальность
	isEmailUnique, err := userR.userDB.IsEmailUnique(userData.Email)
	if err != nil {
		return "", &models.CustomError{ErrorType: 500, OriginalError: err}
	}
	if !isEmailUnique {
		return "", &models.CustomError{ErrorType: 400, OriginalError: nil, Message: NotUniqueEmailMessage}
	}
	isNicknameUnique, err := userR.userDB.IsNicknameUnique(userData.Nickname)
	if err != nil {
		return "", &models.CustomError{ErrorType: 500, OriginalError: err}
	}
	if !isNicknameUnique {
		return "", &models.CustomError{ErrorType: 400, OriginalError: nil, Message: NotUniqueNicknameMessage}
	}

	// 2) Валидация данных (email, nickname, password)
	isValidCredentials, msg, err := ValidateRegisterCredentials(userData)
	if err != nil {
		return "", &models.CustomError{ErrorType: 500, OriginalError: err}
	}
	if !isValidCredentials {
		return "", &models.CustomError{ErrorType: 400, OriginalError: nil, Message: msg}
	}

	// 3) Создание пользователя в базе
	userID, err := userR.userDB.CreateUser(userData)
	if err != nil {
		return "", &models.CustomError{ErrorType: 500, OriginalError: err}
	}

	// 4) Создание сессии в Redis
	sessionToken, err := userR.redisStore.StoreSession(userID)
	if err != nil {
		return "", &models.CustomError{ErrorType: 500, OriginalError: err}
	}

	return sessionToken, nil
}

func (userR UserUseCase) Login(authData models.Auth) (string, *models.CustomError) {
	// 1) Проверка что пользователь существует в базе
	userID, err := userR.userDB.DoesUserExist(authData)
	if err != nil {
		return "", &models.CustomError{ErrorType: 500, OriginalError: err}
	}
	if userID == 0 {
		return "", &models.CustomError{ErrorType: 400, OriginalError: nil, Message: WrongCredentialsMessage}
	}

	log.Println("OK")

	// 2) Создание сессии в Redis
	sessionToken, err := userR.redisStore.StoreSession(userID)
	if err != nil {
		return "", &models.CustomError{ErrorType: 500, OriginalError: err}
	}

	log.Println(sessionToken)

	return sessionToken, nil
}

func (userR UserUseCase) IsAuthorized(cookieValue string) (bool, *models.CustomError) {
	// 1) Получаем id пользователя по сессии
	_, err := userR.redisStore.GetSessionUserId(cookieValue)
	if err != nil {
		return false, &models.CustomError{ErrorType: 500, OriginalError: err}
	}

	return true, nil
}

func (userR UserUseCase) Logout(cookieValue string) {
	userR.redisStore.DeleteSession(cookieValue)
}

func (userR UserUseCase) GetSettings(cookieValue string) (*models.SettingsGet, *models.CustomError) {
	// 1) Получаем ID пользователя из redis по значению куки
	userID, err := userR.redisStore.GetSessionUserId(cookieValue)
	if err != nil {
		return nil, &models.CustomError{ErrorType: 500, OriginalError: err}
	}

	// 2) Получаем настройки пользователя из базы по его ID
	settings, err := userR.userDB.GetSettings(userID)
	if err != nil {
		return nil,  &models.CustomError{ErrorType: 500, OriginalError: err}
	}

	return settings, nil
}

func (userR UserUseCase) UploadSettings(cookieValue string, file *multipart.FileHeader, oldSettingsData *models.SettingsGet, settingsData models.SettingsUpload) *models.CustomError {
	// 1) Проверка, что пользователь авторизован
	userID, err := userR.redisStore.GetSessionUserId(cookieValue)
	if err != nil {
		return &models.CustomError{
			ErrorType: 401,
			OriginalError: err,
		}
	}

	// 2) Проверка что введен правильный пароль
	isCorrect, err := userR.userDB.CheckPasswordByUserID(userID, settingsData.OldPassword)
	if err != nil {
		return &models.CustomError{
			ErrorType: 500,
			OriginalError: err,
		}
	}
	if !isCorrect {
		return &models.CustomError{
			ErrorType: 400,
			OriginalError: nil,
			Message: "Wrong password",
		}
	}

	// 3) Проверка email и nickname на уникальность
	if oldSettingsData.Email != settingsData.Email {
		isEmailUnique, err := userR.userDB.IsEmailUnique(settingsData.Email)
		if err != nil {
			return &models.CustomError{
				ErrorType: 500,
				OriginalError: err,
			}
		}
		if !isEmailUnique {
			return &models.CustomError{
				ErrorType: 400,
				OriginalError: nil,
				Message: "Email is already taken",
			}
		}
	}
	if oldSettingsData.Nickname != settingsData.Nickname {
		isNicknameUnique, err := userR.userDB.IsNicknameUnique(settingsData.Nickname)
		if err != nil {
			return &models.CustomError{
				ErrorType: 500,
				OriginalError: err,
			}
		}
		if !isNicknameUnique {
			return &models.CustomError{
				ErrorType: 400,
				OriginalError: nil,
				Message: "Nickname is already taken",
			}
		}
	}

	// 4) Валидация новых данных
	isValidCredentials, msg, err := ValidateRegisterCredentials(
		models.User{Email: settingsData.Email, Password: settingsData.NewPassword, Nickname: settingsData.Nickname})
	if err != nil {
		return &models.CustomError{
			ErrorType: 500,
			OriginalError: err,
		}
	}
	if !isValidCredentials {
		return &models.CustomError{
			ErrorType: 400,
			OriginalError: nil,
			Message: msg,
		}
	}

	err = userR.userDB.UploadSettings(userID, file, settingsData)
	if err != nil {
		return &models.CustomError{
			ErrorType: 500,
			OriginalError: err,
		}
	}

	return nil
}
