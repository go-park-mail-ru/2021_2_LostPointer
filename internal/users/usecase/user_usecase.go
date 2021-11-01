package usecase

import (
	session "2021_2_LostPointer/internal/microservices/authorization/delivery"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/users"
	"2021_2_LostPointer/internal/utils/constants"
	"2021_2_LostPointer/internal/utils/images"
	"2021_2_LostPointer/internal/utils/validation"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"regexp"
)


type UserUseCase struct {
	userDB	   	   users.UserRepository
	sessionChecker session.SessionCheckerClient
}

func NewUserUserCase(userDB users.UserRepository, sessionChecker session.SessionCheckerClient) UserUseCase {
	return UserUseCase{
		userDB: userDB,
		sessionChecker: sessionChecker,
	}
}

func (userR UserUseCase) Register(userData *models.User) (string, *models.CustomError) {
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

func (userR UserUseCase) Login(authData *models.Auth) (string, *models.CustomError) {
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
	if newSettings.Email != oldSettings.Email && len(newSettings.Email) != 0 {
		isEmailValid, err := regexp.MatchString(constants.EmailRegexPattern, newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isEmailValid {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.InvalidEmailMessage}
		}

		isEmailUnique, err := userR.userDB.IsEmailUnique(newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isEmailUnique {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.NotUniqueEmailMessage}
		}

		err = userR.userDB.UpdateEmail(userID, newSettings.Email)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
	}

	if newSettings.Nickname != oldSettings.Nickname && len(newSettings.Nickname) != 0 {
		isNicknameValid, err := regexp.MatchString(`^[a-zA-Z0-9_-]{` + constants.MinNicknameLength + `,` + constants.MaxNicknameLength + `}$`, newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isNicknameValid {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.NickNameValidationInvalidLengthMessage}
		}

		isNicknameUnique, err := userR.userDB.IsNicknameUnique(newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isNicknameUnique {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.NotUniqueNicknameMessage}
		}

		err = userR.userDB.UpdateNickname(userID, newSettings.Nickname)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
	}

	if len(newSettings.OldPassword) != 0 && len(newSettings.NewPassword) != 0 {
		isOldPasswordCorrect, err := userR.userDB.CheckPasswordByUserID(userID, newSettings.OldPassword)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isOldPasswordCorrect {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.WrongPasswordMessage}
		}

		isNewPasswordValid, msg, err := validation.ValidatePassword(newSettings.NewPassword)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		if !isNewPasswordValid {
			return &models.CustomError{ErrorType: http.StatusBadRequest, Message: msg}
		}

		err = userR.userDB.UpdatePassword(userID, newSettings.NewPassword)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
	} else if len(newSettings.OldPassword) == 0 && len(newSettings.NewPassword) != 0 {
		return &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.OldPasswordFieldIsEmptyMessage}
	} else if len(newSettings.OldPassword) != 0 && len(newSettings.NewPassword) == 0 {
		return &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.NewPasswordFieldIsEmptyMessage}
	}

	if len(newSettings.AvatarFileName) != 0 {
		createdAvatarFilename, err := images.CreateImage(newSettings.Avatar)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}

		oldAvatarFilename, err := userR.userDB.GetAvatarFilename(userID)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}
		err = images.DeleteImage(oldAvatarFilename)
		if err != nil {
			return &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: err}
		}

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
	return os.Getenv("ROOT_PATH_PREFIX") + filename + constants.LittleAvatarPostfix, nil
}
