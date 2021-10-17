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
	server.Use(middleware.CheckAuthentication)
}

func (middleware Middleware) CheckAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var isAuthorized bool

		cookie, err := ctx.Cookie("Session_cookie")

		if err == nil && cookie.String() != "" {
			isAuthorized, _ = middleware.UserUseCase.IsAuthorized(cookie.Value)
		}

		ctx.Set("is_authorized", isAuthorized)

		return next(ctx)
	}

}
