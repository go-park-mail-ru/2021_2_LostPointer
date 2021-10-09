package delivery

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/users"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"time"
)

const cookieLifetime = time.Hour * 24 * 30

type UserDelivery struct {
	userLogic users.UserUseCase
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

func NewUserDelivery(userRealization users.UserUseCase) UserDelivery {
	return UserDelivery{userLogic: userRealization}
}

func (userD UserDelivery) InitHandlers(server *echo.Echo) {
	server.POST("/api/v1/user/signup", userD.Register)
}
