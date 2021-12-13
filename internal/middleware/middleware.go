package middleware

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/csrf"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/monitoring/delivery"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"

	authorization "2021_2_LostPointer/internal/microservices/authorization/proto"
)

type Middleware struct {
	logger           *zap.SugaredLogger
	authMicroservice authorization.AuthorizationClient
	metrics          *delivery.PrometheusMetrics
}

func NewMiddlewareHandler(authMicroservice authorization.AuthorizationClient, logger *zap.SugaredLogger, monitoring *delivery.PrometheusMetrics) Middleware {
	return Middleware{
		logger:           logger,
		authMicroservice: authMicroservice,
		metrics:          monitoring,
	}
}

func (middleware *Middleware) InitMiddlewareHandlers(server *echo.Echo) {
	server.Use(middleware.AccessLog)
	server.Use(middleware.CheckAuthorization)
	server.Use(middleware.CSRF)
}

func (middleware *Middleware) CheckAuthorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		//cookie, err := ctx.Cookie("Session_cookie")
		//userID := &authorization.UserID{
		//	ID: -1,
		//}
		//if err == nil {
		//	userID, err = middleware.authMicroservice.GetUserByCookie(context.Background(), &authorization.Cookie{
		//		Cookies: cookie.Value,
		//	})
		//	if err != nil {
		//		cookie = &http.Cookie{Expires: time.Now().AddDate(0, 0, -1)}
		//		ctx.SetCookie(cookie)
		//		ctx.Set("USER_ID", -1)
		//		return next(ctx)
		//	}
		//}
		//if err != nil {
		//	cookie = &http.Cookie{Expires: time.Now().AddDate(0, 0, -1)}
		//	ctx.SetCookie(cookie)
		//}
		//
		//ctx.Set("USER_ID", int(userID.ID))

		ctx.Set("USER_ID", 268)

		return next(ctx)
	}
}

func (middleware *Middleware) PanicRecovering(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		//nolint:errcheck
		defer func() error {
			if err := recover(); err != nil {
				requestID, ok := ctx.Get("REQUEST_ID").(string)
				if !ok {
					return ctx.NoContent(http.StatusInternalServerError)
				}
				middleware.logger.Info(
					zap.String("ID", requestID),
					zap.String("ERROR", err.(error).Error()),
					zap.Int("ANSWER STATUS", http.StatusInternalServerError),
				)

				log.Println("panic")

				status := strconv.Itoa(ctx.Response().Status)
				path := ctx.Request().URL.Path
				method := ctx.Request().Method

				middleware.metrics.Hits.WithLabelValues(status, path, method).Inc()
				middleware.metrics.Duration.WithLabelValues(status, path, method).Observe(0)
				return ctx.JSON(http.StatusInternalServerError, &models.Response{
					Status:  http.StatusInternalServerError,
					Message: constants.PanicRecover,
				})
			}
			return nil
		}()

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

		status := strconv.Itoa(ctx.Response().Status)
		path := ctx.Request().URL.Path
		method := ctx.Request().Method

		middleware.metrics.Hits.WithLabelValues(status, path, method).Inc()
		middleware.metrics.Duration.WithLabelValues(status, path, method).Observe(respTime.Seconds())

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
