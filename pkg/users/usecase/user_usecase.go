package usecase

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/users"
	"log"
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
	fileSystem users.FileSystemIFace
}

func NewUserUserCase(userDB users.UserRepositoryIFace, redisStore users.RedisStoreIFace, fileSystem users.FileSystemIFace) UserUseCase {
	return UserUseCase{
		userDB: userDB,
		redisStore: redisStore,
		fileSystem: fileSystem,
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

func (userR UserUseCase) IsAuthorized(cookieValue string) (bool, int, *models.CustomError) {
	// 1) Получаем id пользователя по сессии
	id, customError := userR.redisStore.GetSessionUserId(cookieValue)
	if customError != nil {
		return false, id, customError
	}
	return true, id, nil
}

func (userR UserUseCase) Logout(cookieValue string) {
	userR.redisStore.DeleteSession(cookieValue)
}

func (userR UserUseCase) GetSettings(userID int) (*models.SettingsGet, *models.CustomError) {
	// 1) Получаем настройки пользователя из базы по его ID
	settings, err := userR.userDB.GetSettings(userID)
	if err != nil {
		return nil,  &models.CustomError{ErrorType: 500, OriginalError: err}
	}

	return settings, nil
}

func (userR UserUseCase) UpdateSettings(userID int, oldSettings *models.SettingsGet, newSettings *models.SettingsUpload) *models.CustomError {
	// 1) Проверяем, что изменился email
	if newSettings.Email != oldSettings.Email && len(newSettings.Email) != 0 {
		// 1.1) Валидация нового email
		isEmailValid, err := regexp.MatchString(`[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`, newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
		if !isEmailValid {
			return &models.CustomError{ErrorType: 400, Message: InvalidEmailMessage}
		}

		// 1.2) Проверка, что новый email уникален
		isEmailUnique, err := userR.userDB.IsEmailUnique(newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
		if !isEmailUnique {
			return &models.CustomError{ErrorType: 400, Message: NotUniqueEmailMessage}
		}

		// 1.3) Обновляем email в базе
		err = userR.userDB.UpdateEmail(userID, newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
	}

	// 2) Проверяем, что изменился nickname
	if newSettings.Nickname != oldSettings.Nickname && len(newSettings.Nickname) != 0 {
		// 2.1) Валидация нового nickname
		isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{` + minNicknameLength + `,` + maxNicknameLength + `}$`, newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
		if !isNicknameValid {
			return &models.CustomError{ErrorType: 400, Message: NickNameValidationInvalidLengthMessage}
		}

		// 2.2) Проврека, что новый nickname уникален
		isNicknameUnique, err := userR.userDB.IsNicknameUnique(newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
		if !isNicknameUnique {
			return &models.CustomError{ErrorType: 400, Message: NotUniqueNicknameMessage}
		}

		// 2.3) Обновляем nickname в базе
		err = userR.userDB.UpdateNickname(userID, newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
	}

	log.Println(len(newSettings.OldPassword), len(newSettings.NewPassword))

	// 3) Проверяем, что изменили пароль
	if len(newSettings.OldPassword) != 0 && len(newSettings.NewPassword) != 0 {
		// 3.1) Проверка, что старый пароль введен правильно
		isOldPasswordCorrect, err := userR.userDB.CheckPasswordByUserID(userID, newSettings.OldPassword)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
		if !isOldPasswordCorrect {
			return &models.CustomError{ErrorType: 400, Message: "Wrong password"}
		}

		// 3.2) Валидация нового пароля
		isNewPasswordValid, msg, err := ValidatePassword(newSettings.NewPassword)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
		if !isNewPasswordValid {
			return &models.CustomError{ErrorType: 400, Message: msg}
		}

		// 3.3) Обновляем пароль в базе
		err = userR.userDB.UpdatePassword(userID, newSettings.NewPassword)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
	} else if len(newSettings.OldPassword) == 0 && len(newSettings.NewPassword) != 0 {
		return &models.CustomError{ErrorType: 400, Message: "Old password field is empty"}
	} else if len(newSettings.OldPassword) != 0 && len(newSettings.NewPassword) == 0 {
		return &models.CustomError{ErrorType: 400, Message: "New password field is empty"}
	}

	// 4) Проверяем, что изменили аватарку
	if len(newSettings.AvatarFileName) != 0 {
		// 4.1) Создаем файл, получаем его название
		createdAvatarFilename, err := userR.fileSystem.CreateImage(newSettings.Avatar)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}

		// 4.2) Удаляем старый файл
		oldAvatarFilename, err := userR.userDB.GetAvatarFilename(userID)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
		err = userR.fileSystem.DeleteImage(oldAvatarFilename)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}

		// 4.3) Обновляем аватарку в базе
		err = userR.userDB.UpdateAvatar(userID, createdAvatarFilename)
		if err != nil {
			return &models.CustomError{ErrorType: 500, OriginalError: err}
		}
	}

	return nil
}

