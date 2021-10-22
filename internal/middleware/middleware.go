package middleware

import (
	"2021_2_LostPointer/internal/csrf"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/users"
	"errors"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

type Middleware struct {
	logger *zap.SugaredLogger
	UserUseCase users.UserUseCase
}

func NewMiddlewareHandler(logger *zap.SugaredLogger, userUseCase users.UserUseCase) Middleware {
	return Middleware{
		UserUseCase: userUseCase,
		logger: logger,
	}
}

func (middleware Middleware) InitMiddlewareHandlers(server *echo.Echo) {
	server.Use(middleware.CheckAuthorization)
	server.Use(middleware.AccessLog)
	server.Use(middleware.CSRF)
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

func (middleware Middleware) CSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(rwContext echo.Context) error {
		if rwContext.Request().RequestURI == "/api/v1/user/settings" || rwContext.Request().Method == "PATCH" {
			cookie, err := rwContext.Cookie("Session_cookie")
			if err != nil {
				middleware.logger.Debug(
					zap.String("COOKIE", errors.New("cookie expired").Error()),
					zap.Int("ANSWER STATUS", http.StatusUnauthorized),
				)

				return rwContext.JSON(http.StatusUnauthorized, &models.Response{
					Status: http.StatusUnauthorized,
					Message: "Cookie expired",
				})
			}

			tokenReq := rwContext.Request().Header.Get("X-CSRF-Token")
			log.Println(tokenReq)

			isValidCsrf, err := csrf.Tokens.Check(cookie.Value, tokenReq)

			if err != nil {
				return rwContext.JSON(http.StatusForbidden, &models.Response{
					Status: http.StatusForbidden,
					Message: "Cookie expired",
				})
			}

			if !isValidCsrf  {
				return rwContext.JSON(http.StatusForbidden, &models.Response{
					Status: http.StatusForbidden,
					Message: "Cookie expired",
				})
			}
		}
		return next(rwContext)
	}
}
