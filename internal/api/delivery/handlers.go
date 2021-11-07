package delivery

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	authorization "2021_2_LostPointer/internal/microservices/authorization/proto"
	"2021_2_LostPointer/internal/models"
)

type APIMicroservices struct {
	logger           *zap.SugaredLogger
	authMicroservice authorization.AuthorizationClient
}

func NewAPIMicroservices(logger *zap.SugaredLogger, auth authorization.AuthorizationClient) APIMicroservices {
	return APIMicroservices{
		logger:           logger,
		authMicroservice: auth,
	}
}

type MyStringType string
type MyIntType int

func (api *APIMicroservices) ParseErrorByCode(ctx echo.Context, requestID string, err error) error {
	if e, temp := status.FromError(err); temp {
		if e.Code() == codes.Internal {
			api.logger.Error(
				zap.String("ID", requestID),
				zap.String("ERROR", err.Error()),
				zap.Int("ANSWER STATUS", http.StatusInternalServerError))
			return ctx.NoContent(http.StatusInternalServerError)
		}
		if e.Code() == codes.Aborted {
			api.logger.Info(
				zap.String("ID", requestID),
				zap.String("ERROR", e.Message()),
				zap.Int("ANSWER STATUS", http.StatusBadRequest))
			return ctx.JSON(http.StatusOK, &models.Response{
				Status:  http.StatusBadRequest,
				Message: e.Message(),
			})
		}
	}
	return nil
}

func (api *APIMicroservices) Login(ctx echo.Context) error {
	requestID, _ := ctx.Get("REQUEST_ID").(string)
	var authData models.AuthData

	if err := ctx.Bind(&authData); err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	cookies, err := api.authMicroservice.Login(context.Background(),
		&authorization.AuthData{Login: authData.Email, Password: authData.Password})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
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
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: constants.UserAuthorizedMessage,
	})
}

func (api *APIMicroservices) Register(ctx echo.Context) error {
	requestID, _ := ctx.Get("REQUEST_ID").(string)
	var registerData models.RegisterData

	if err := ctx.Bind(&registerData); err != nil {
		api.logger.Error(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
		return ctx.NoContent(http.StatusInternalServerError)
	}

	cookies, err := api.authMicroservice.Register(context.Background(),
		&authorization.RegisterData{Login: registerData.Email, Password: registerData.Password, Nickname: registerData.Nickname})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
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
	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusCreated),
	)

	return ctx.JSON(http.StatusCreated, &models.Response{
		Status:  http.StatusCreated,
		Message: constants.UserCreatedMessage,
	})
}

func (api *APIMicroservices) GetUserAvatar(ctx echo.Context) error {
	requestID, _ := ctx.Get("REQUEST_ID").(string)
	userID, _ := ctx.Get("USER_ID").(int)
	if userID == -1 {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
	}

	avatar, err := api.authMicroservice.GetAvatar(context.Background(), &authorization.UserID{ID: int64(userID)})
	if err != nil {
		return api.ParseErrorByCode(ctx, requestID, err)
	}

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)
	return ctx.JSON(http.StatusOK,
		struct {
			Status int    `json:"status"`
			Avatar string `json:"avatar"`
		}{http.StatusOK, avatar.Filename})
}

func (api *APIMicroservices) Logout(ctx echo.Context) error {
	requestID, _ := ctx.Get("REQUEST_ID").(string)
	cookie, err := ctx.Cookie("Session_cookie")
	if err != nil {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("ERROR", constants.UserIsNotAuthorizedMessage),
			zap.Int("ANSWER STATUS", http.StatusUnauthorized))
		return ctx.JSON(http.StatusOK, &models.Response{
			Status:  http.StatusUnauthorized,
			Message: constants.UserIsNotAuthorizedMessage,
		})
	}

	_, err = api.authMicroservice.Logout(context.Background(), &authorization.Cookie{Cookies: cookie.Value})
	if err != nil {
		api.logger.Info(
			zap.String("ID", requestID),
			zap.String("ERROR", err.Error()),
			zap.Int("ANSWER STATUS", http.StatusConflict))
		return ctx.NoContent(http.StatusConflict)
	}
	cookie.Expires = time.Now().AddDate(0, 0, -1)

	api.logger.Info(
		zap.String("ID", requestID),
		zap.Int("ANSWER STATUS", http.StatusOK),
	)

	return ctx.JSON(http.StatusOK, &models.Response{
		Status:  http.StatusOK,
		Message: constants.LoggedOutMessage,
	})
}

func (api *APIMicroservices) Init(server *echo.Echo) {
	// Authorization
	server.POST("/api/v1/user/signin", api.Login)
	server.POST("/api/v1/user/signup", api.Register)
	server.GET("/api/v1/auth", api.GetUserAvatar)
	server.POST("/api/v1/user/logout", api.Logout)
}
