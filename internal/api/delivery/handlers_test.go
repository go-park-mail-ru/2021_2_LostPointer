package delivery

import (
	"2021_2_LostPointer/internal/constants"
	authorizationMock "2021_2_LostPointer/internal/microservices/authorization/mock"
	authMicroservice "2021_2_LostPointer/internal/microservices/authorization/proto"
	authorizationProto "2021_2_LostPointer/internal/microservices/authorization/proto"
	musicMock "2021_2_LostPointer/internal/microservices/music/mock"
	musicMicroservice "2021_2_LostPointer/internal/microservices/music/proto"
	playlistsMicroservice "2021_2_LostPointer/internal/microservices/playlists/proto"
	profileMock "2021_2_LostPointer/internal/microservices/profile/mock"
	profileMicroservice "2021_2_LostPointer/internal/microservices/profile/proto"
	profileProto "2021_2_LostPointer/internal/microservices/profile/proto"
	"2021_2_LostPointer/pkg/image"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestAPIMicroservices_Login(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		email             string
		password          string
		mock              func(*gomock.Controller) *authorizationMock.MockAuthorizationClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
	}{
		{
			name:     "Handler returned status 200",
			email:    "testEmail",
			password: "testPassword",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Login(gomock.Any(), &authorizationProto.AuthData{
					Email:    "testEmail",
					Password: "testPassword",
				}).Return(&authorizationProto.Cookie{Cookies: "cookie"}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":200,\"message\":\"User is authorized\"}\n",
		},
		{
			name:     "Handler returned status 400, Login error",
			email:    "testEmail",
			password: "testPassword",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Login(gomock.Any(), &authorizationProto.AuthData{
					Email:    "testEmail",
					Password: "testPassword",
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name:     "No RequestID",
			email:    "testEmail",
			password: "testPassword",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Login(gomock.Any(), &authorizationProto.AuthData{
					Email:    "testEmail",
					Password: "testPassword",
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/user/signin",
				strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, currentTest.email, currentTest.password)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			authManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManagerMock, profileManager, musicManager, playlistsManager)
			if assert.NoError(t, r.Login(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_Register(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		email             string
		password          string
		nickname          string
		mock              func(*gomock.Controller) *authorizationMock.MockAuthorizationClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
	}{
		{
			name:     "Handler returned status 201",
			email:    "",
			password: "",
			nickname: "",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Register(gomock.Any(), &authorizationProto.RegisterData{
					Email:    "",
					Password: "",
					Nickname: "",
				}).Return(&authorizationProto.Cookie{Cookies: "cookie"}, nil)
				return moq
			},
			expectedStatus: http.StatusCreated,
			expectedJSON:   "{\"status\":201,\"message\":\"User was created successfully\"}\n",
		},
		{
			name:     "Handler returned status 400, Login error",
			email:    "",
			password: "",
			nickname: "",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Register(gomock.Any(), &authorizationProto.RegisterData{
					Email:    "",
					Password: "",
					Nickname: "",
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name:     "No RequestID",
			email:    "testEmail",
			password: "testPassword",
			nickname: "testNickName",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Register(gomock.Any(), &authorizationProto.RegisterData{
					Email:    "testEmail",
					Password: "testPassword",
					Nickname: "testNickname",
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/user/signup",
				strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "%s", "nickname": "%s"}`, currentTest.email, currentTest.password, currentTest.nickname)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			authManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManagerMock, profileManager, musicManager, playlistsManager)
			if assert.NoError(t, r.Register(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_GetUserAvatar(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	var ID int64

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *authorizationMock.MockAuthorizationClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
		doNotSetUserID    bool
		userID            int
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().GetAvatar(gomock.Any(), &authorizationProto.UserID{ID: ID}).Return(&authorizationProto.Avatar{Filename: "testFilename"}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":200,\"avatar\":\"testFilename\"}\n",
		},
		{
			name: "Handler returned status 400, Login error",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().GetAvatar(gomock.Any(), &authorizationProto.UserID{ID: ID}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().GetAvatar(gomock.Any(), &authorizationProto.UserID{ID: ID}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().GetAvatar(gomock.Any(), &authorizationProto.UserID{ID: 1}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   "",
			doNotSetUserID: true,
		},
		{
			name: "Not authorized -> UserID = -1",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().GetAvatar(gomock.Any(), &authorizationProto.UserID{ID: 1}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":401,\"message\":\"User is not authorized\"}\n",
			userID:         -1,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/auth", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSetUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			authManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManagerMock, profileManager, musicManager, playlistsManager)
			if assert.NoError(t, r.GetUserAvatar(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

//func (api *APIMicroservices) Logout(ctx echo.Context) error {
//	requestID, ok := ctx.Get("REQUEST_ID").(string)
//	if !ok {
//		api.logger.Error(
//			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
//			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//	cookie, err := ctx.Cookie("Session_cookie")
//	if err != nil {
//		api.logger.Info(
//			zap.String("ID", requestID),
//			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
//			zap.Int("ANSWER STATUS", http.StatusUnauthorized))
//		return ctx.JSON(http.StatusOK, &models.Response{
//			Status:  http.StatusUnauthorized,
//			Message: constants.UserIsNotAuthorizedMessage,
//		})
//	}
//
//	_, err = api.authMicroservice.Logout(context.Background(), &authorization.Cookie{Cookies: cookie.Value})
//	if err != nil {
//		api.logger.Info(
//			zap.String("ID", requestID),
//			zap.String("MESSAGE", err.Error()),
//			zap.Int("ANSWER STATUS", http.StatusConflict))
//		return ctx.NoContent(http.StatusConflict)
//	}
//	cookie = &http.Cookie{
//		Name:     "Session_cookie",
//		Value:    "",
//		Path:     "/",
//		Secure:   true,
//		HttpOnly: true,
//		SameSite: http.SameSiteNoneMode,
//		Expires:  time.Now().AddDate(0, 0, -1),
//	}
//	ctx.SetCookie(cookie)
//
//	api.logger.Info(
//		zap.String("ID", requestID),
//		zap.Int("ANSWER STATUS", http.StatusOK),
//	)
//
//	return ctx.JSON(http.StatusOK, &models.Response{
//		Status:  http.StatusOK,
//		Message: constants.LoggedOutMessage,
//	})
//}

func TestAPIMicroservices_Logout(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *authorizationMock.MockAuthorizationClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
		doNotSetCookie    bool
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Logout(gomock.Any(), &authorizationProto.Cookie{Cookies: "testCookie"}).Return(&authorizationProto.Empty{}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":200,\"message\":\"Logged out\"}\n",
		},
		{
			name: "Handler returned status 400, Login error",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Logout(gomock.Any(), &authorizationProto.Cookie{Cookies: "testCookie"}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusConflict,
			expectedJSON:   "",
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Logout(gomock.Any(), &authorizationProto.Cookie{Cookies: "testCookie"}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
		{
			name: "No Cookie",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Logout(gomock.Any(), &authorizationProto.Cookie{Cookies: "testCookie"}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":401,\"message\":\"User is not authorized\"}\n",
			doNotSetCookie: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/user/signin",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			if !currentTest.doNotSetCookie {
				req.AddCookie(&http.Cookie{
					Name:       "Session_cookie",
					Value:      "testCookie",
					Path:       "",
					Domain:     "",
					Expires:    time.Time{},
					RawExpires: "",
					MaxAge:     0,
					Secure:     false,
					HttpOnly:   false,
					SameSite:   0,
					Raw:        "",
					Unparsed:   nil,
				})
			}

			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			authManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManagerMock, profileManager, musicManager, playlistsManager)
			if assert.NoError(t, r.Logout(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_GetSettings(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	authConn, _ := grpc.Dial(
		os.Getenv("AUTH_HOST"),
		grpc.WithInsecure(),
	)
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	var ID int64

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *profileMock.MockProfileClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
		doNotSetUserID    bool
		userID            int
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *profileMock.MockProfileClient {
				moq := profileMock.NewMockProfileClient(controller)
				moq.EXPECT().GetSettings(gomock.Any(), &profileProto.GetSettingsOptions{ID: ID}).Return(&profileProto.UserSettings{}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{}\n",
		},
		{
			name: "Handler returned status 400, Login error",
			mock: func(controller *gomock.Controller) *profileMock.MockProfileClient {
				moq := profileMock.NewMockProfileClient(controller)
				moq.EXPECT().GetSettings(gomock.Any(), &profileProto.GetSettingsOptions{ID: ID}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *profileMock.MockProfileClient {
				moq := profileMock.NewMockProfileClient(controller)
				moq.EXPECT().GetSettings(gomock.Any(), &profileProto.GetSettingsOptions{ID: ID}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *profileMock.MockProfileClient {
				moq := profileMock.NewMockProfileClient(controller)
				moq.EXPECT().GetSettings(gomock.Any(), &profileProto.GetSettingsOptions{ID: ID}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   "",
			doNotSetUserID: true,
		},
		{
			name: "Not authorized -> UserID = -1",
			mock: func(controller *gomock.Controller) *profileMock.MockProfileClient {
				moq := profileMock.NewMockProfileClient(controller)
				moq.EXPECT().GetSettings(gomock.Any(), &profileProto.GetSettingsOptions{ID: ID}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":401,\"message\":\"User is not authorized\"}\n",
			userID:         -1,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/auth", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSetUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			authManager := authMicroservice.NewAuthorizationClient(authConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			profileManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManagerMock, musicManager, playlistsManager)
			if assert.NoError(t, r.GetSettings(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

//func (api *APIMicroservices) GenerateCSRF(ctx echo.Context) error {
//	requestID, ok := ctx.Get("REQUEST_ID").(string)
//	if !ok {
//		api.logger.Error(
//			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
//			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//	userID, ok := ctx.Get("USER_ID").(int)
//	if !ok {
//		api.logger.Error(
//			zap.String("ID", requestID),
//			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
//			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//	if userID == -1 {
//		api.logger.Info(
//			zap.String("ID", requestID),
//			zap.String("MESSAGE", constants.UserIsNotAuthorizedMessage),
//			zap.Int("ANSWER STATUS", http.StatusUnauthorized))
//		return ctx.JSON(http.StatusOK, &models.Response{
//			Status:  http.StatusUnauthorized,
//			Message: constants.UserIsNotAuthorizedMessage,
//		})
//	}
//
//	cookie, _ := ctx.Cookie("Session_cookie")
//	token, _ := csrf.Tokens.Create(cookie.Value, time.Now().Unix()+constants.CSRFTokenLifetime)
//	return ctx.JSON(http.StatusOK, &models.Response{
//		Status:  http.StatusOK,
//		Message: token,
//	})
//}

func TestAPIMicroservices_GenerateCSRF(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)
	authConn, _ := grpc.Dial(
		os.Getenv("AUTH_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		password          string
		mock              func(*gomock.Controller) *authorizationMock.MockAuthorizationClient
		expectedStatus    int
		doNotSetRequestID bool
		doNotSetUserID    bool
		doNotSetCookie    bool
		userID            int
	}{
		{
			name:           "Handler returned status 200",
			password:       "testPassword",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Handler returned status 400, Unauthorized error",
			password:       "testPassword",
			expectedStatus: http.StatusOK,
			userID:         -1,
		},
		{
			name:              "No RequestID",
			password:          "testPassword",
			expectedStatus:    http.StatusInternalServerError,
			doNotSetRequestID: true,
		},
		{
			name:           "No UserID",
			password:       "testPassword",
			expectedStatus: http.StatusInternalServerError,
			doNotSetUserID: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/csrf",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSetUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}
			if !currentTest.doNotSetCookie {
				req.AddCookie(&http.Cookie{
					Name:  "Session_cookie",
					Value: "Session_cookie",
				})
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			imageServices := image.NewImagesService()

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManager, playlistsManager)
			if assert.NoError(t, r.GenerateCSRF(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
			}
		})
	}
}

//func (api *APIMicroservices) GetHomeTracks(ctx echo.Context) error {
//	requestID, ok := ctx.Get("REQUEST_ID").(string)
//	if !ok {
//		api.logger.Error(
//			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
//			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//	userID, ok := ctx.Get("USER_ID").(int)
//	if !ok {
//		api.logger.Error(
//			zap.String("ID", requestID),
//			zap.String("ERROR", constants.UserIDTypeAssertionFailed),
//			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//	var isAuthorized bool
//	if userID != -1 {
//		isAuthorized = true
//	}
//
//	tracksListProto, err := api.musicMicroservice.RandomTracks(context.Background(),
//		&music.RandomTracksOptions{Amount: constants.HomePageTracksSelectionAmount, IsAuthorized: isAuthorized})
//	if err != nil {
//		return api.ParseErrorByCode(ctx, requestID, err)
//	}
//
//	tracks := make([]models.Track, 0, constants.HomePageTracksSelectionAmount)
//	for _, current := range tracksListProto.Tracks {
//		var track models.Track
//		track.BindProto(current)
//		tracks = append(tracks, track)
//	}
//	api.logger.Info(
//		zap.String("ID", requestID),
//		zap.Int("ANSWER STATUS", http.StatusOK),
//	)
//
//	return ctx.JSON(http.StatusOK, tracks)
//}

func TestAPIMicroservices_GetHomeTracks(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	authConn, _ := grpc.Dial(
		os.Getenv("AUTH_HOST"),
		grpc.WithInsecure(),
	)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *musicMock.MockMusicClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
		doNotSetUserID    bool
		userID            int
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomTracks(gomock.Any(), &musicMicroservice.RandomTracksOptions{
					Amount:       constants.HomePageTracksSelectionAmount,
					IsAuthorized: true,
				}).Return(&musicMicroservice.Tracks{Tracks: []*musicMicroservice.Track{&musicMicroservice.Track{
					ID:          0,
					Title:       "",
					Explicit:    false,
					Genre:       "",
					Number:      0,
					File:        "",
					ListenCount: 0,
					Duration:    0,
					Lossless:    false,
					Album:       &musicMicroservice.Album{},
					Artist:      &musicMicroservice.Artist{},
				}}}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "[{\"album\":{},\"artist\":{\"name\":\"\"}}]\n",
		},
		{
			name: "Handler returned status 400, Login error",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomTracks(gomock.Any(), &musicMicroservice.RandomTracksOptions{
					Amount:       constants.HomePageTracksSelectionAmount,
					IsAuthorized: true,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomTracks(gomock.Any(), &musicMicroservice.RandomTracksOptions{}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomTracks(gomock.Any(), &musicMicroservice.RandomTracksOptions{}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   "",
			doNotSetUserID: true,
		},
		{
			name: "User in unauthorized",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomTracks(gomock.Any(), &musicMicroservice.RandomTracksOptions{
					Amount:       constants.HomePageTracksSelectionAmount,
					IsAuthorized: false,
				}).Return(&musicMicroservice.Tracks{}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "[]\n",
			userID:         -1,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/home/tracks",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSetUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.GetHomeTracks(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_GetHomeAlbums(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	authConn, _ := grpc.Dial(
		os.Getenv("AUTH_HOST"),
		grpc.WithInsecure(),
	)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *musicMock.MockMusicClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomAlbums(gomock.Any(), &musicMicroservice.RandomAlbumsOptions{Amount: constants.HomePageAlbumsSelectionAmount}).
					Return(&musicMicroservice.Albums{Albums: []*musicMicroservice.Album{&musicMicroservice.Album{}}}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "[{}]\n",
		},
		{
			name: "Handler returned status 400, Login error",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomAlbums(gomock.Any(), &musicMicroservice.RandomAlbumsOptions{Amount: constants.HomePageAlbumsSelectionAmount}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomAlbums(gomock.Any(), &musicMicroservice.RandomAlbumsOptions{Amount: constants.HomePageAlbumsSelectionAmount}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/home/albums",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.GetHomeAlbums(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

//func (api *APIMicroservices) GetHomeArtists(ctx echo.Context) error {
//	requestID, ok := ctx.Get("REQUEST_ID").(string)
//	if !ok {
//		api.logger.Error(
//			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
//			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//
//	artistsListProto, err := api.musicMicroservice.RandomArtists(context.Background(), &music.RandomArtistsOptions{Amount: constants.HomePageArtistsSelectionAmount})
//	if err != nil {
//		return api.ParseErrorByCode(ctx, requestID, err)
//	}
//
//	artists := make([]models.Artist, 0, constants.HomePageArtistsSelectionAmount)
//	for _, current := range artistsListProto.Artists {
//		var artist models.Artist
//		artist.BindProto(current)
//		artists = append(artists, artist)
//	}
//	api.logger.Info(
//		zap.String("ID", requestID),
//		zap.Int("ANSWER STATUS", http.StatusOK),
//	)
//
//	return ctx.JSON(http.StatusOK, artists)
//}

func TestAPIMicroservices_GetHomeArtists(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	authConn, _ := grpc.Dial(
		os.Getenv("AUTH_HOST"),
		grpc.WithInsecure(),
	)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *musicMock.MockMusicClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomArtists(gomock.Any(), &musicMicroservice.RandomArtistsOptions{Amount: constants.HomePageArtistsSelectionAmount}).
					Return(&musicMicroservice.Artists{Artists: []*musicMicroservice.Artist{&musicMicroservice.Artist{
						ID:     0,
						Name:   "",
						Avatar: "",
						Video:  "",
						Tracks: []*musicMicroservice.Track{},
						Albums: []*musicMicroservice.Album{},
					}}}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "[{\"name\":\"\"}]\n",
		},
		{
			name: "Handler returned status 400, Login error",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomArtists(gomock.Any(), &musicMicroservice.RandomArtistsOptions{Amount: constants.HomePageArtistsSelectionAmount}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomArtists(gomock.Any(), &musicMicroservice.RandomArtistsOptions{Amount: constants.HomePageArtistsSelectionAmount}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/home/artists",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.GetHomeArtists(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_GetArtistProfile(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	authConn, _ := grpc.Dial(
		os.Getenv("AUTH_HOST"),
		grpc.WithInsecure(),
	)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *musicMock.MockMusicClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
		doNotSetUserID    bool
		doNotSetIDParam   bool
		userID            int
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: true,
				}).
					Return(&musicMicroservice.Artist{
						ID:     0,
						Name:   "",
						Avatar: "",
						Video:  "",
						Tracks: []*musicMicroservice.Track{},
						Albums: []*musicMicroservice.Album{},
					}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"name\":\"\"}\n",
		},
		{
			name: "Handler returned status 400, Login error",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: true,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: true,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
		{
			name: "User is unauthorized",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: false,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
			userID:         -1,
		},
		{
			name: "No userID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: false,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   "",
			doNotSetUserID: true,
		},
		{
			name: "No id param",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: false,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   "",
			doNotSetIDParam: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/home/artist/:id",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetIDParam {
				ctx.SetParamNames("id")
				ctx.SetParamValues("1")
			}

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSetUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.GetArtistProfile(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

//func (api *APIMicroservices) IncrementListenCount(ctx echo.Context) error {
//	requestID, ok := ctx.Get("REQUEST_ID").(string)
//	if !ok {
//		api.logger.Error(
//			zap.String("ERROR", constants.RequestIDTypeAssertionFailed),
//			zap.Int("ANSWER STATUS", http.StatusInternalServerError))
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//
//	var trackID models.TrackID
//	err := ctx.Bind(&trackID)
//	if err != nil {
//		api.logger.Error(
//			zap.String("ID", requestID),
//			zap.String("ERROR", err.Error()),
//			zap.Int("ANSWER STATUS", http.StatusInternalServerError),
//		)
//		return ctx.NoContent(http.StatusInternalServerError)
//	}
//
//	_, err = api.musicMicroservice.IncrementListenCount(context.Background(), &music.IncrementListenCountOptions{ID: trackID.ID})
//	if err != nil {
//		return api.ParseErrorByCode(ctx, requestID, err)
//	}
//
//	api.logger.Info(
//		zap.String("ID", requestID),
//		zap.Int("ANSWER STATUS", http.StatusOK),
//	)
//
//	return ctx.JSON(http.StatusOK, &models.Response{
//		Status:  http.StatusOK,
//		Message: "Incremented track listen count",
//	})
//}

func TestAPIMicroservices_IncrementListenCount(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer func(prLogger *zap.Logger) {
		_ = prLogger.Sync()
	}(prLogger)
	authConn, _ := grpc.Dial(
		os.Getenv("AUTH_HOST"),
		grpc.WithInsecure(),
	)
	profileConn, _ := grpc.Dial(
		os.Getenv("PROFILE_HOST"),
		grpc.WithInsecure(),
	)
	playlistsConn, _ := grpc.Dial(
		os.Getenv("PLAYLISTS_HOST"),
		grpc.WithInsecure(),
	)

	var ID int64 = 1

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *musicMock.MockMusicClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
		trackID int
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.IncrementListenCountOptions{ID: ID}).
					Return(&musicMicroservice.IncrementListenCountEmpty{}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"name\":\"\"}\n",
			trackID: 1,
		},
		{
			name: "Handler returned status 400, Login error",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: true,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: true,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON:      "",
			doNotSetRequestID: true,
		},
		{
			name: "User is unauthorized",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: false,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
		},
		{
			name: "No userID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: false,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   "",
		},
		{
			name: "No id param",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
					IsAuthorized: false,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   "",
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/home/artist/:id",
				strings.NewReader(fmt.Sprintf(`{"trackID": "%d"}`, currentTest.trackID)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.IncrementListenCount(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}
