package middleware

import (
	"2021_2_LostPointer/pkg/users"
	"github.com/labstack/echo"
)

type Middleware struct {
	UserUseCase users.UserUseCaseIFace
}

func NewMiddlewareHandler(userUseCase users.UserUseCaseIFace) Middleware {
	return Middleware{UserUseCase: userUseCase}
}

func (middleware Middleware) InitMiddlewareHandlers(server *echo.Echo) {
	server.Use(middleware.CheckAuthorization)
}

func (middleware Middleware) CheckAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var isAuthorized bool
		var userID int

		cookie, err := ctx.Cookie("Session_cookie")

		if err == nil && cookie.String() != "" {
			isAuthorized, userID = middleware.UserUseCase.IsAuthorized(cookie.Value)
		}

		ctx.Set("is_authorized", isAuthorized)
		ctx.Set("user_id", userID)

		return next(ctx)
	}

}
