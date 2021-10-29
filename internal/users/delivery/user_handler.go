package delivery

import (
	"2021_2_LostPointer/internal/csrf"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/users"
	"github.com/labstack/echo"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const cookieLifetime = time.Hour * 24 * 30

const UserIsNotAuthorizedMessage = "User is not authorized"
const LoggedOutMessage = "Logged out"
const SettingsUploadedMessage = "Settings were uploaded successfully"
const UserCreatedMessage = "User was created successfully"
const UserAuthorizedMessage = "User is authorized"

type UserDelivery struct {
	userLogic users.UserUseCase
	logger    *zap.SugaredLogger
}

func NewUserDelivery(logger *zap.SugaredLogger, userRealization users.UserUseCase) UserDelivery {
	return UserDelivery{userLogic: userRealization, logger: logger}
}

func (userD UserDelivery) Register(ctx echo.Context) error {
	var userData models.User
	requestID := ctx.Get("REQUEST_ID").(string)

	err := ctx.Bind(&userData)
	if err != nil {
		userD.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)
		return ctx.JSON(http.StatusInternalServerError, &models.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	cookieValue, customError := userD.userLogic.Register(userData)
	if customError != nil {
		if customError.ErrorType == http.StatusInternalServerError {
			userD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			return ctx.JSON(http.StatusInternalServerError, &models.Response{
				Status:  http.StatusInternalServerError,
				Message: customError.OriginalError.Error(),
			})
		} else if customError.ErrorType == http.StatusBadRequest {
			userD.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.Message),
				zap.Int("ANSWER STATUS", http.StatusBadRequest),
			)

			return ctx.JSON(http.StatusOK, &models.Response{
				Status:  http.StatusBadRequest,
				Message: customError.Message,
			})
		}
	}

	cookie := &http.Cookie{
		Name:     "Session_cookie",
		Value:    cookieValue,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(cookieLifetime),
	}
	ctx.SetCookie(cookie)

	userD.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)

	return ctx.JSON(http.StatusCreated, &models.Response{
		Status:  http.StatusCreated,
		Message: UserCreatedMessage,
	})
}

func (userD UserDelivery) Login(ctx echo.Context) error {
	var authData models.Auth
	requestID := ctx.Get("REQUEST_ID").(string)

	err := ctx.Bind(&authData)
	if err != nil {
		userD.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
		)

		return ctx.JSON(http.StatusInternalServerError, &models.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	cookieValue, customError := userD.userLogic.Login(authData)
	if customError != nil {
		if customError.ErrorType == http.StatusInternalServerError {
			userD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			return ctx.JSON(http.StatusInternalServerError, &models.Response{
				Status:  http.StatusInternalServerError,
				Message: customError.OriginalError.Error(),
			})
		} else if customError.ErrorType == http.StatusBadRequest {
			userD.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.Message),
				zap.Int("ANSWER STATUS", http.StatusBadRequest),
			)

			return ctx.JSON(http.StatusOK, &models.Response{
				Status:  http.StatusBadRequest,
				Message: customError.Message,
			})
		}
	}

	cookie := &http.Cookie{
		Name:     "Session_cookie",
		Value:    cookieValue,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(cookieLifetime),
	}
	ctx.SetCookie(cookie)

	userD.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: UserAuthorizedMessage,
	})
}

func (userD UserDelivery) GetAvatarForMainPage(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	userID := ctx.Get("USER_ID").(int)

	avatarFilename, customError := userD.userLogic.GetAvatarFilename(userID)
	if customError != nil {
		if customError.ErrorType == http.StatusInternalServerError {
			userD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			return ctx.JSON(http.StatusInternalServerError, &models.Response{
				Status:  http.StatusInternalServerError,
				Message: customError.OriginalError.Error(),
			})
		}
	}

	userD.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK,
		struct {
			Status int    `json:"status"`
			Avatar string `json:"avatar"`
		}{http.StatusOK, avatarFilename})
}

func (userD UserDelivery) Logout(ctx echo.Context) error {
	cookie, err := ctx.Cookie("Session_cookie")
	requestID := ctx.Get("REQUEST_ID").(string)
	if err != nil {
		userD.logger.Info(
			zap.String("ID", requestID),
			zap.String("ERROR", UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized),
		)

		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: UserIsNotAuthorizedMessage,
		})
	}

	err = userD.userLogic.Logout(cookie.Value)
	if err != nil {
		userD.logger.Debug(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusConflict),
		)

		return ctx.NoContent(http.StatusConflict)
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	ctx.SetCookie(cookie)

	userD.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: LoggedOutMessage,
	})
}

func (userD UserDelivery) GetSettings(ctx echo.Context) error {
	requestID := ctx.Get("REQUEST_ID").(string)
	userID := ctx.Get("USER_ID").(int)
	if userID == -1 {
		userD.logger.Info(
			zap.String("ID", requestID),
			zap.String("ERROR", UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized),
		)

		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: UserIsNotAuthorizedMessage,
		})
	}

	settings, customError := userD.userLogic.GetSettings(userID)
	if customError != nil {
		if customError.ErrorType == http.StatusInternalServerError {
			userD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			return ctx.JSON(http.StatusInternalServerError, &models.Response{
				Status:  http.StatusInternalServerError,
				Message: customError.OriginalError.Error(),
			})
		}
	}

	userD.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, settings)
}

func (userD UserDelivery) UpdateSettings(ctx echo.Context) error {
	userID := ctx.Get("USER_ID").(int)
	requestID := ctx.Get("REQUEST_ID").(string)
	if userID == -1 {
		userD.logger.Info(
			zap.String("ID", requestID),
			zap.String("ERROR", UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized),
		)

		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: UserIsNotAuthorizedMessage,
		})
	}

	oldSettings, customError := userD.userLogic.GetSettings(userID)
	if customError != nil {
		if customError.ErrorType == 500 {
			userD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			return ctx.JSON(http.StatusInternalServerError, &models.Response{
				Status:  http.StatusInternalServerError,
				Message: customError.OriginalError.Error(),
			})
		}
	}

	var fileName string

	email := ctx.FormValue("email")
	nickname := ctx.FormValue("nickname")
	oldPassword := ctx.FormValue("old_password")
	newPassword := ctx.FormValue("new_password")
	file, err := ctx.FormFile("avatar")
	if err != nil {
		fileName = ""
	} else {
		fileName = file.Filename
	}

	newSettings := &models.SettingsUpload{
		Email:          email,
		Nickname:       nickname,
		OldPassword:    oldPassword,
		NewPassword:    newPassword,
		Avatar:         file,
		AvatarFileName: fileName,
	}

	customError = userD.userLogic.UpdateSettings(userID, oldSettings, newSettings)
	if customError != nil {
		if customError.ErrorType == 500 {
			userD.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.OriginalError.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError),
			)

			return ctx.JSON(http.StatusInternalServerError, &models.Response{
				Status:  http.StatusInternalServerError,
				Message: customError.OriginalError.Error(),
			})
		} else if customError.ErrorType == 400 {
			userD.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", customError.Message),
				zap.Int("ANSWER STATUS", http.StatusBadRequest),
			)

			return ctx.JSON(http.StatusOK, &models.Response{
				Status:  http.StatusBadRequest,
				Message: customError.Message,
			})
		}
	}

	userD.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: SettingsUploadedMessage,
	})
}

func (userD UserDelivery) GetCsrf(ctx echo.Context) error {
	userID := ctx.Get("USER_ID").(int)
	requestID := ctx.Get("REQUEST_ID").(string)
	if userID == -1  {
		userD.logger.Info(
			zap.String("ID", requestID),
			zap.String("ERROR", UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized),
		)

		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: UserIsNotAuthorizedMessage,
		})
	}
	cookie, _ := ctx.Cookie("Session_cookie")
	token, _ := csrf.Tokens.Create(cookie.Value, 900+time.Now().Unix())

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: token,
	})
}

func (userD UserDelivery) InitHandlers(server *echo.Echo) {
	server.POST("/api/v1/user/signup", userD.Register)
	server.POST("/api/v1/user/signin", userD.Login)
	server.POST("/api/v1/user/logout", userD.Logout)
	server.GET("/api/v1/auth", userD.GetAvatarForMainPage)
	server.GET("/api/v1/user/settings", userD.GetSettings)
	server.GET("api/v1/csrf", userD.GetCsrf)
	server.PATCH("/api/v1/user/settings", userD.UpdateSettings)
}
