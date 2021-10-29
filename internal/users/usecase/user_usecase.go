package usecase

import (
	session "2021_2_LostPointer/internal/microservices/authorization/delivery"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/users"
	"2021_2_LostPointer/internal/utils/validation"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
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
	const OldPasswordFieldIsEmptyMessage = "Old password field is empty"
const NewPasswordFieldIsEmptyMessage = "New password field is empty"
const WrongPasswordMessage = "Wrong password"

const EmailRegexPattern = `[a-zA-Z0-9]+@[a-zA-Z0-9]+\.[a-zA-Z0-9]+`

type UserUseCase struct {
	userDB	   	   users.UserRepository
	redisStore 	   users.RedisStore
	fileSystem 	   users.FileSystem
	sessionChecker session.SessionCheckerClient
}

func NewUserUserCase(userDB users.UserRepository, redisStore users.RedisStore, fileSystem users.FileSystem,
	sessionChecker session.SessionCheckerClient) UserUseCase {
	return UserUseCase{
		userDB: userDB,
		redisStore: redisStore,
		fileSystem: fileSystem,
		sessionChecker: sessionChecker,
	}
}

func (userR UserUseCase) Register(userData models.User) (string, *models.CustomError) {
	cookie, err := userR.sessionChecker.Signup(context.Background(), &session.SignUpData{
		Email: userData.Email,
		Password: userData.Password,
		Nickname: userData.Nickname,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Aborted:
				return "", &models.CustomError{ErrorType: http.StatusBadRequest, Message: e.Message()}
			case codes.Internal:
				return "", &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
			}
		}
	}

	return cookie.Cookies, nil
}

func (userR UserUseCase) Login(authData models.Auth) (string, *models.CustomError) {
	cookie, err := userR.sessionChecker.SignIn(context.Background(), &session.Auth{
		Login: authData.Email,
		Password: authData.Password,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Aborted:
				return "", &models.CustomError{ErrorType: http.StatusBadRequest, Message: e.Message()}
			case codes.Internal:
				return "", &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
			}
		}
	}

	return cookie.Cookies, nil
}

func (userR UserUseCase) Logout(cookieValue string) error {
	_, err := userR.sessionChecker.DeleteSession(context.Background(), &session.SessionData{
		Cookies: cookieValue,
	})

	return err
}

func (userR UserUseCase) GetSettings(userID int) (*models.SettingsGet, *models.CustomError) {
	settings, err := userR.userDB.GetSettings(userID)
	if err != nil {
		return nil,  &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
	}

	return settings, nil
}

func (userR UserUseCase) UpdateSettings(userID int, oldSettings *models.SettingsGet, newSettings *models.SettingsUpload) *models.CustomError {
	// 1) Проверяем, что изменился email
	if newSettings.Email != oldSettings.Email && len(newSettings.Email) != 0 {
		// 1.1) Валидация нового email
		isEmailValid, err := regexp.MatchString(EmailRegexPattern, newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isEmailValid {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: InvalidEmailMessage}
		}

		// 1.2) Проверка, что новый email уникален
		isEmailUnique, err := userR.userDB.IsEmailUnique(newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isEmailUnique {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: NotUniqueEmailMessage}
		}

		// 1.3) Обновляем email в базе
		err = userR.userDB.UpdateEmail(userID, newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
	}

	// 2) Проверяем, что изменился nickname
	if newSettings.Nickname != oldSettings.Nickname && len(newSettings.Nickname) != 0 {
		// 2.1) Валидация нового nickname
		isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{` + minNicknameLength + `,` + maxNicknameLength + `}$`, newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isNicknameValid {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: NickNameValidationInvalidLengthMessage}
		}

		// 2.2) Проврека, что новый nickname уникален
		isNicknameUnique, err := userR.userDB.IsNicknameUnique(newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isNicknameUnique {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: NotUniqueNicknameMessage}
		}

		// 2.3) Обновляем nickname в базе
		err = userR.userDB.UpdateNickname(userID, newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
	}

	// 3) Проверяем, что изменили пароль
	if len(newSettings.OldPassword) != 0 && len(newSettings.NewPassword) != 0 {
		// 3.1) Проверка, что старый пароль введен правильно
		isOldPasswordCorrect, err := userR.userDB.CheckPasswordByUserID(userID, newSettings.OldPassword)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isOldPasswordCorrect {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: WrongPasswordMessage}
		}

		// 3.2) Валидация нового пароля
		isNewPasswordValid, msg, err := validation.ValidatePassword(newSettings.NewPassword)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isNewPasswordValid {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: msg}
		}

		// 3.3) Обновляем пароль в базе
		err = userR.userDB.UpdatePassword(userID, newSettings.NewPassword)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
	} else if len(newSettings.OldPassword) == 0 && len(newSettings.NewPassword) != 0 {
		return &models.CustomError{ErrorType: http.StatusBadRequest, Message: OldPasswordFieldIsEmptyMessage}
	} else if len(newSettings.OldPassword) != 0 && len(newSettings.NewPassword) == 0 {
		return &models.CustomError{ErrorType: http.StatusBadRequest, Message: NewPasswordFieldIsEmptyMessage}
	}

	// 4) Проверяем, что изменили аватарку
	if len(newSettings.AvatarFileName) != 0 {
		// 4.1) Создаем файл, получаем его название
		createdAvatarFilename, err := userR.fileSystem.CreateImage(newSettings.Avatar)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}

		// 4.2) Удаляем старый файл
		oldAvatarFilename, err := userR.userDB.GetAvatarFilename(userID)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		err = userR.fileSystem.DeleteImage(oldAvatarFilename)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}

		// 4.3) Обновляем аватарку в базе
		err = userR.userDB.UpdateAvatar(userID, createdAvatarFilename)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
	}

	return nil
}

func (userR UserUseCase) GetAvatarFilename(userID int) (string, *models.CustomError) {
	filename, err := userR.userDB.GetAvatarFilename(userID)
	if err != nil {
		return "", &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
	}
	return os.Getenv("ROOT_PATH_PREFIX") + filename + "_150px.webp", nil
}
