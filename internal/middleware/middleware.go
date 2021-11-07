package middleware

import (
	"2021_2_LostPointer/internal/csrf"
	"2021_2_LostPointer/internal/models"
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	authorization "2021_2_LostPointer/internal/microservices/authorization/proto"
)

type Middleware struct {
	logger           *zap.SugaredLogger
	authMicroservice authorization.AuthorizationClient
}

func NewMiddlewareHandler(authMicroservice authorization.AuthorizationClient, logger *zap.SugaredLogger) Middleware {
	return Middleware{
		logger:           logger,
		authMicroservice: authMicroservice,
	}
}

func (middleware *Middleware) InitMiddlewareHandlers(server *echo.Echo) {
	server.Use(middleware.CheckAuthorization)
	server.Use(middleware.AccessLog)
	server.Use(middleware.CORS)
	server.Use(middleware.CSRF)
}

func (middleware *Middleware) CheckAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		cookie, err := ctx.Cookie("Session_cookie")
		userID := &authorization.UserID{
			ID: -1,
		}

		if err == nil {
			userID, err = middleware.authMicroservice.GetUserByCookie(context.Background(), &authorization.Cookie{
				Cookies: cookie.Value,
			})
		}
		if err != nil {
			cookie = &http.Cookie{Expires: time.Now().AddDate(0, 0, -1)}
			ctx.SetCookie(cookie)
		}

		ctx.Set("USER_ID", int(userID.ID))

		return next(ctx)
	}
}

func (middleware *Middleware) AccessLog(next echo.HandlerFunc) echo.HandlerFunc {
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
		if rwContext.Request().Method == "PATCH" {
			cookie, err := rwContext.Cookie("Session_cookie")
			if err != nil {
				middleware.logger.Debug(
					zap.String("COOKIE", errors.New("cookie expired").Error()),
					zap.Int("ANSWER STATUS", http.StatusUnauthorized),
				)

				return rwContext.JSON(http.StatusUnauthorized, &models.Response{
					Status:  http.StatusUnauthorized,
					Message: "Cookie expired",
				})
			}

			tokenReq := rwContext.Request().Header.Get("X-CSRF-Token")

			isValidCsrf, err := csrf.Tokens.Check(cookie.Value, tokenReq)

			if err != nil {
				return rwContext.JSON(http.StatusForbidden, &models.Response{
					Status:  http.StatusForbidden,
					Message: "Cookie expired",
				})
			}

			if !isValidCsrf {
				return rwContext.JSON(http.StatusForbidden, &models.Response{
					Status:  http.StatusForbidden,
					Message: "Cookie expired",
				})
			}
		}
		return next(rwContext)
	}
}

func (middleware *Middleware) CORS(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		c.Response().Header().Set("Access-Control-Allow-Origin", os.Getenv("CORS_ORIGIN"))
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, PUT, DELETE, POST, PATCH")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, X-Login, Set-Cookie, Content-Type, Content-Length, Accept-Encoding, X-Csrf-Token, csrf-token, Authorization")
		c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
		c.Response().Header().Set("Vary", "Cookie")

		if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusOK)
		}

		return next(c)
	}
}
