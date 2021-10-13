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

	sessionToken, msg, err := userD.userLogic.Register(userData)
	if err != nil {
		log.Println(err.Error())
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if len(sessionToken) == 0 {
		return ctx.JSON(http.StatusBadRequest, &models.Response{Message: msg})
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

	sessionToken, err := userD.userLogic.Login(authData)
	if err != nil {
		log.Println(err.Error())
		return ctx.NoContent(http.StatusInternalServerError)
	}
	if len(sessionToken) == 0 {
		return ctx.JSON(http.StatusBadRequest, &models.Response{Message: "Wrong username or password"})
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

	isAuthorized, err := userD.userLogic.IsAuthorized(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: "User not authorized"})
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

	settings, err := userD.userLogic.GetSettings(cookie.Value)
	if err != nil {
		log.Println(err.Error())
		return ctx.JSON(http.StatusUnauthorized, &models.Response{Message: "User not authorized"})
	}

	return ctx.JSON(http.StatusOK, settings)
}

func (userD UserDelivery) InitHandlers(server *echo.Echo) {
	server.POST("/api/v1/user/signup", userD.Register)
	server.POST("/api/v1/user/signin", userD.Login)
	server.POST("/api/v1/user/logout", userD.Logout)
	server.GET("/api/v1/auth", userD.IsAuthorized)
	server.GET("/api/v1/user/settings", userD.GetSettings)
}
