package delivery

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/users"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"time"
)

const cookieLifetime = time.Hour * 24 * 30

type UserDelivery struct {
	userLogic users.UserUseCaseIFace
}

func NewUserDelivery(userRealization users.UserUseCaseIFace) UserDelivery {
	return UserDelivery{userLogic: userRealization}
}

func (userD UserDelivery) Register(ctx echo.Context) error {
	var userData models.User

	err := ctx.Bind(&userData)
	if err != nil {
		log.Println(err.Error())
		return ctx.NoContent(http.StatusInternalServerError)
	}

	sessionToken, customError := userD.userLogic.Register(userData)
	if customError != nil {
		if customError.ErrorType == 500 {
			log.Println(customError.OriginalError.Error())
			return ctx.NoContent(http.StatusInternalServerError)
		} else {
			return ctx.JSON(http.StatusBadRequest, &models.Response{Message: customError.Message})
		}
	}

	cookie := &http.Cookie{
		Name:     "Session_cookie",
		Value:    sessionToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(cookieLifetime),
	}
	ctx.SetCookie(cookie)

	return ctx.JSON(http.StatusCreated, &models.Response{Message: "User successfully created"})
}

func (userD UserDelivery) Login(ctx echo.Context) error {
	var authData models.Auth

	err := ctx.Bind(&authData)
	if err != nil {
		log.Println(err.Error())
		return ctx.NoContent(http.StatusInternalServerError)
	}

	sessionToken, customError := userD.userLogic.Login(authData)
	if customError != nil {
		if customError.ErrorType == 500 {
			log.Println(customError.OriginalError.Error())
			return ctx.NoContent(http.StatusInternalServerError)
		} else {
			return ctx.JSON(http.StatusBadRequest, &models.Response{Message: customError.Message})
		}

	}

	cookie := &http.Cookie{
		Name:     "Session_cookie",
		Value:    sessionToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(cookieLifetime),
	}
	ctx.SetCookie(cookie)

	return ctx.JSON(http.StatusOK, &models.Response{Message: "User is authorized"})
}

func (userD UserDelivery) IsAuthorized(ctx echo.Context) error {
	cookie, err := ctx.Cookie("Session_cookie")
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: "User not authorized"})
	}

	isAuthorized, customError := userD.userLogic.IsAuthorized(cookie.Value)
	if customError != nil {
		if customError.ErrorType == 401 {
			return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: customError.Message})
		}
	}
	if !isAuthorized {
		return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: "User not authorized"})
	}

	return ctx.JSON(http.StatusOK, &models.Response{Message: "User is authorized"})
}

func (userD UserDelivery) Logout(ctx echo.Context) error {
	cookie, err := ctx.Cookie("Session_cookie")
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: "User not authorized"})
	}
	userD.userLogic.Logout(cookie.Value)
	cookie.Expires = time.Now().AddDate(0, 0, -1)
	ctx.SetCookie(cookie)

	return ctx.JSON(http.StatusOK, &models.Response{Message: "Logged out"})
}

func (userD UserDelivery) GetSettings(ctx echo.Context) error {
	cookie, err := ctx.Cookie("Session_cookie")
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: "User not authorized"})
	}

	settings, customError := userD.userLogic.GetSettings(cookie.Value)
	if customError != nil {
		if customError.ErrorType == 401 {
			return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: customError.Message})
		} else if customError.ErrorType == 500 {
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}

	return ctx.JSON(http.StatusOK, settings)
}

func (userD UserDelivery) UpdateSettings(ctx echo.Context) error {
	cookie, err := ctx.Cookie("Session_cookie")
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: "User not authorized"})
	}

	oldSettings, customError := userD.userLogic.GetSettings(cookie.Value)
	if customError != nil {
		if customError.ErrorType == 401 {
			return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: customError.Message})
		} else if customError.ErrorType == 500 {
			return ctx.NoContent(http.StatusInternalServerError)
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
		Email: email,
		Nickname: nickname,
		OldPassword: oldPassword,
		NewPassword: newPassword,
		Avatar: file,
		AvatarFileName: fileName,
	}

	customError = userD.userLogic.UpdateSettings(cookie.Value, oldSettings, newSettings)
	if customError != nil {
		if customError.ErrorType == 500 {
			log.Println(customError.OriginalError.Error())
			return ctx.NoContent(http.StatusInternalServerError)
		} else {
			return ctx.JSON(http.StatusBadRequest, &models.Response{Message: customError.Message})
		}
	}

	return ctx.JSON(http.StatusOK, &models.Response{Message: "Settings uploaded successfully"})
}

func (userD UserDelivery) InitHandlers(server *echo.Echo) {
	server.POST("/api/v1/user/signup", userD.Register)
	server.POST("/api/v1/user/signin", userD.Login)
	server.POST("/api/v1/user/logout", userD.Logout)
	server.GET("/api/v1/auth", userD.IsAuthorized)
	server.GET("/api/v1/user/settings", userD.GetSettings)
	server.PATCH("/api/v1/user/settings", userD.UpdateSettings)
}
