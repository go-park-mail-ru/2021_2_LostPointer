package middleware

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/users"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"time"
)

type Middleware struct {
	logger *zap.SugaredLogger
	UserUseCase users.UserUseCaseIFace
}

func NewMiddlewareHandler(logger *zap.SugaredLogger, userUseCase users.UserUseCaseIFace) Middleware {
	return Middleware{
		UserUseCase: userUseCase,
		logger: logger,
	}
}

func (middleware Middleware) InitMiddlewareHandlers(server *echo.Echo) {
	server.Use(middleware.CheckAuthorization)
	server.Use(middleware.AccessLog)
}

func (middleware Middleware) CheckAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var isAuthorized bool
		var userID int
		var customError *models.CustomError
		var authorizationErrorMessage string

		cookie, err := ctx.Cookie("Session_cookie")

		if err == nil && cookie.String() != "" {
			isAuthorized, userID, customError = middleware.UserUseCase.IsAuthorized(cookie.Value)
		}
		if customError != nil {
			if customError.ErrorType == 500 {
				authorizationErrorMessage = customError.OriginalError.Error()
			}
		}

		ctx.Set("IS_AUTHORIZED", isAuthorized)
		ctx.Set("USER_ID", userID)
		ctx.Set("AUTHORIZATION_ERROR", authorizationErrorMessage)

		return next(ctx)
	}
}

func (middleware Middleware) AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		uniqueID := uuid.NewV4()
		start := time.Now()
		ctx.Set("REQUEST_ID", uniqueID.String())

		middleware.logger.Info(
			zap.String("ID", uniqueID.String()),
			zap.String("URL", ctx.Request().URL.Path),
			zap.String("METHOD", ctx.Request().Method),
		)

		err := next(ctx)

		respTime := time.Since(start)
		middleware.logger.Info(
			zap.String("ID", uniqueID.String()),
			zap.Duration("TIME FOR ANSWER", respTime),
		)

		return err
	}
}
