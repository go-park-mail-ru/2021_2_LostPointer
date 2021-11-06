package delivery

import (
	"2021_2_LostPointer/internal/constants"
	authorization "2021_2_LostPointer/internal/microservices/authorization/proto"
	"2021_2_LostPointer/internal/models"
	"context"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"
	"time"
)

type ApiMicroservices struct {
	authMicroservice authorization.AuthorizationClient
}

func NewApiMicroservices(auth authorization.AuthorizationClient) ApiMicroservices {
	return ApiMicroservices{
		authMicroservice: auth,
	}
}

func (api *ApiMicroservices) Login(ctx echo.Context) error {
	var authData models.AuthData

	err := ctx.Bind(&authData)
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	cookies, err := api.authMicroservice.Login(context.Background(),
		&authorization.AuthData{Login: authData.Email, Password: authData.Password})
	if err != nil {
		if e, temp := status.FromError(err); temp {
			switch e.Code() {
			case codes.Aborted:
				log.Println(err)
				return ctx.NoContent(http.StatusBadRequest)
			case codes.Internal:
				log.Println(err)
				return ctx.NoContent(http.StatusInternalServerError)
			}
		}
	}
	cookie := &http.Cookie{
		Name:     "Session_cookie",
		Value:    cookies.Cookies,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Now().Add(constants.CookieLifetime),
	}
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}

func (api *ApiMicroservices) Init(server *echo.Echo) {
	// Authorization
	server.POST("/api/v1/user/signin", api.Login)
}
