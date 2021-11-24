package delivery

import (
	"2021_2_LostPointer/internal/constants"
	authorizationMock "2021_2_LostPointer/internal/microservices/authorization/mock"
	authMicroservice "2021_2_LostPointer/internal/microservices/authorization/proto"
	authorizationProto "2021_2_LostPointer/internal/microservices/authorization/proto"
	musicMock "2021_2_LostPointer/internal/microservices/music/mock"
	musicMicroservice "2021_2_LostPointer/internal/microservices/music/proto"
	playlistsMock "2021_2_LostPointer/internal/microservices/playlists/mock"
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
	"strconv"
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
			name:     "Handler returned status 400",
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
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Register(gomock.Any(), &authorizationProto.RegisterData{
				}).Return(&authorizationProto.Cookie{Cookies: "cookie"}, nil)
				return moq
			},
			expectedStatus: http.StatusCreated,
			expectedJSON:   "{\"status\":201,\"message\":\"User was created successfully\"}\n",
		},
		{
			name:     "Handler returned status 400",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Register(gomock.Any(), &authorizationProto.RegisterData{
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
			name: "Handler returned status 400",
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
			req := httptest.NewRequest(echo.GET, "/api/v1/auth", strings.NewReader(""))
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
			name: "Handler returned status 409",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Logout(gomock.Any(), &authorizationProto.Cookie{Cookies: "testCookie"}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *authorizationMock.MockAuthorizationClient {
				moq := authorizationMock.NewMockAuthorizationClient(controller)
				moq.EXPECT().Logout(gomock.Any(), &authorizationProto.Cookie{Cookies: "testCookie"}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
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
			req := httptest.NewRequest(echo.POST, "/api/v1/user/logout",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			if !currentTest.doNotSetCookie {
				req.AddCookie(&http.Cookie{
					Name:       "Session_cookie",
					Value:      "testCookie",
					Expires:    time.Time{},
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
			name: "Handler returned status 400",
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
			req := httptest.NewRequest(echo.GET, "/api/v1/user/settings", strings.NewReader(""))
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
			name:           "Handler returned status 401, Unauthorized user",
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
			req := httptest.NewRequest(echo.GET, "/api/v1/csrf",
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
					Album:       &musicMicroservice.Album{},
					Artist:      &musicMicroservice.Artist{},
				}}}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "[{\"album\":{},\"artist\":{\"name\":\"\"}}]\n",
		},
		{
			name: "Handler returned status 400",
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
			doNotSetUserID: true,
		},
		{
			name: "User in unauthorized",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().RandomTracks(gomock.Any(), &musicMicroservice.RandomTracksOptions{
					Amount:       constants.HomePageTracksSelectionAmount,
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
			req := httptest.NewRequest(echo.GET, "/api/v1/home/tracks",
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
			name: "Handler returned status 400",
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
			doNotSetRequestID: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/api/v1/home/albums",
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
						Tracks: []*musicMicroservice.Track{},
						Albums: []*musicMicroservice.Album{},
					}}}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "[{\"name\":\"\"}]\n",
		},
		{
			name: "Handler returned status 400",
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
			doNotSetRequestID: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/api/v1/home/artists",
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
						Tracks: []*musicMicroservice.Track{},
						Albums: []*musicMicroservice.Album{},
					}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"name\":\"\"}\n",
		},
		{
			name: "Handler returned status 400",
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
			doNotSetRequestID: true,
		},
		{
			name: "User is unauthorized",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
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
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			doNotSetUserID: true,
		},
		{
			name: "No id param",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().ArtistProfile(gomock.Any(), &musicMicroservice.ArtistProfileOptions{
					ArtistID:     1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:  http.StatusInternalServerError,
			doNotSetIDParam: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/api/v1/home/artist/:id",
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
		trackID           int64
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().IncrementListenCount(gomock.Any(), &musicMicroservice.IncrementListenCountOptions{ID: ID}).
					Return(&musicMicroservice.IncrementListenCountEmpty{}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":200,\"message\":\"Incremented track listen count\"}\n",
			trackID:        ID,
		},
		{
			name: "Handler returned status 400",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().IncrementListenCount(gomock.Any(), &musicMicroservice.IncrementListenCountOptions{ID: ID}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
			trackID:        ID,
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().IncrementListenCount(gomock.Any(), &musicMicroservice.IncrementListenCountOptions{ID: ID}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			doNotSetRequestID: true,
			trackID:           ID,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/inc_listencount",
				strings.NewReader(fmt.Sprintf(`{"id": %d}`, currentTest.trackID)))
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

func TestAPIMicroservices_GetAlbumPage(t *testing.T) {
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

	//var ID int64 = 1

	tests := []struct {
		name                 string
		mock                 func(*gomock.Controller) *musicMock.MockMusicClient
		expectedStatus       int
		expectedJSON         string
		doNotSetRequestID    bool
		doNotSerUserID       bool
		userID               int
		albumID              int64
		wrongTypeOfParameter bool
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().AlbumPage(gomock.Any(), &musicMicroservice.AlbumPageOptions{
					IsAuthorized: true,
				}).
					Return(&musicMicroservice.AlbumPageResponse{
						Artist:         &musicMicroservice.Artist{},
						Tracks:         []*musicMicroservice.AlbumTrack{},
					}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"artist\":{\"name\":\"\"}}\n",
		},
		{
			name: "Handler returned status 400",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().AlbumPage(gomock.Any(), &musicMicroservice.AlbumPageOptions{
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
				moq.EXPECT().AlbumPage(gomock.Any(), &musicMicroservice.AlbumPageOptions{}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().AlbumPage(gomock.Any(), &musicMicroservice.AlbumPageOptions{}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			doNotSerUserID: true,
		},
		{
			name: "Wrong type of parameter",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().AlbumPage(gomock.Any(), &musicMicroservice.AlbumPageOptions{}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:       http.StatusInternalServerError,
			wrongTypeOfParameter: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/api/v1/home/album/:id",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			ctx.SetParamNames("id")
			if currentTest.wrongTypeOfParameter {
				ctx.SetParamValues("qwe!123scd")
			} else {
				ctx.SetParamValues(strconv.FormatInt(currentTest.albumID, 10))
			}

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSerUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.GetAlbumPage(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_SearchMusic(t *testing.T) {
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
		doNotSerUserID    bool
		userID            int
		formValue         string
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().Find(gomock.Any(), &musicMicroservice.FindOptions{
					Text:         "testText",
					IsAuthorized: true,
				}).
					Return(&musicMicroservice.FindResponse{
						Tracks:  []*musicMicroservice.Track{},
						Albums:  []*musicMicroservice.Album{},
						Artists: []*musicMicroservice.Artist{},
					}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{}\n",
			formValue:      "testText",
		},
		{
			name: "Handler returned status 400",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().Find(gomock.Any(), &musicMicroservice.FindOptions{
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
				moq.EXPECT().AlbumPage(gomock.Any(), &musicMicroservice.FindOptions{
					IsAuthorized: true,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().AlbumPage(gomock.Any(), &musicMicroservice.FindOptions{
					Text:         "testText",
					IsAuthorized: true,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			doNotSerUserID: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/api/v1/music/search",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			req.PostFormValue("text")
			req.Form.Add("text", currentTest.formValue)
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSerUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.SearchMusic(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_AddTrack(t *testing.T) {
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
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *playlistsMock.MockPlaylistsClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
		doNotSerUserID    bool
		userID            int
		playlistID        int64
		trackID           int64
	}{
		{
			name: "Handler returned status 201",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().AddTrack(gomock.Any(), &playlistsMicroservice.AddTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).
					Return(&playlistsMicroservice.AddTrackResponse{}, nil)
				return moq
			},
			expectedStatus: http.StatusCreated,
			expectedJSON:   "{\"status\":201,\"message\":\"Track was successfully added to playlist\"}\n",
			trackID:        2,
			playlistID:     1,
			userID:         3,
		},
		{
			name: "Handler returned status 400",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().AddTrack(gomock.Any(), &playlistsMicroservice.AddTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
			trackID:        2,
			playlistID:     1,
			userID:         3,
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().AddTrack(gomock.Any(), &playlistsMicroservice.AddTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().AddTrack(gomock.Any(), &playlistsMicroservice.AddTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			doNotSerUserID: true,
		},
		{
			name: "User unauthorized -> UserID = -1",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().AddTrack(gomock.Any(), &playlistsMicroservice.AddTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
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
			req := httptest.NewRequest(echo.POST, "/api/v1/playlist/track",
				strings.NewReader(fmt.Sprintf(`{"track_id": %d, "playlist_id": %d}`, currentTest.trackID, currentTest.playlistID)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSerUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			playlistsManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManager, playlistsManagerMock)
			if assert.NoError(t, r.AddTrack(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_DeletePlaylist(t *testing.T) {
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
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name              string
		mock              func(*gomock.Controller) *playlistsMock.MockPlaylistsClient
		expectedStatus    int
		expectedJSON      string
		doNotSetRequestID bool
		doNotSerUserID    bool
		userID            int
		playlistID        int64
		trackID           int64
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().DeleteTrack(gomock.Any(), &playlistsMicroservice.DeleteTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).
					Return(&playlistsMicroservice.DeleteTrackResponse{}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":200,\"message\":\"Track was successfully deleted from playlist\"}\n",
			trackID:        2,
			playlistID:     1,
			userID:         3,
		},
		{
			name: "Handler returned status 400",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().DeleteTrack(gomock.Any(), &playlistsMicroservice.DeleteTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
			trackID:        2,
			playlistID:     1,
			userID:         3,
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().AddTrack(gomock.Any(), &playlistsMicroservice.AddTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().DeleteTrack(gomock.Any(), &playlistsMicroservice.DeleteTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			doNotSerUserID: true,
		},
		{
			name: "User unauthorized -> UserID = -1",
			mock: func(controller *gomock.Controller) *playlistsMock.MockPlaylistsClient {
				moq := playlistsMock.NewMockPlaylistsClient(controller)
				moq.EXPECT().DeleteTrack(gomock.Any(), &playlistsMicroservice.DeleteTrackOptions{
					PlaylistID: 1,
					TrackID:    2,
					UserID:     3,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
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
			req := httptest.NewRequest(echo.DELETE, "/api/v1/playlist/track",
				strings.NewReader(fmt.Sprintf(`{"track_id": %d, "playlist_id": %d}`, currentTest.trackID, currentTest.playlistID)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSerUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			playlistsManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManager, playlistsManagerMock)
			if assert.NoError(t, r.DeleteTrack(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_GetUserPlaylists(t *testing.T) {
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
		doNotSerUserID    bool
		userID            int
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().UserPlaylists(gomock.Any(), &musicMicroservice.UserPlaylistsOptions{
					UserID: 1,
				}).
					Return(&musicMicroservice.PlaylistsData{
						Playlists: []*musicMicroservice.PlaylistData{},
					}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{}\n",
			userID:         1,
		},
		{
			name: "Handler returned status 400",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().UserPlaylists(gomock.Any(), &musicMicroservice.UserPlaylistsOptions{
					UserID: 1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
			userID:         1,
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().UserPlaylists(gomock.Any(), &musicMicroservice.UserPlaylistsOptions{
					UserID: 1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().UserPlaylists(gomock.Any(), &musicMicroservice.UserPlaylistsOptions{
					UserID: 1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			doNotSerUserID: true,
		},
		{
			name: "User is unauthorized => userId = -1",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().UserPlaylists(gomock.Any(), &musicMicroservice.UserPlaylistsOptions{
					UserID: 1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
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
			req := httptest.NewRequest(echo.GET, "/api/v1/playlists",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSerUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.GetUserPlaylists(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_GetPlaylistPage(t *testing.T) {
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
		name                 string
		mock                 func(*gomock.Controller) *musicMock.MockMusicClient
		expectedStatus       int
		expectedJSON         string
		doNotSetRequestID    bool
		doNotSerUserID       bool
		userID               int
		playlistId           int64
		wrongTypeOfParameter bool
	}{
		{
			name: "Handler returned status 200",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().PlaylistPage(gomock.Any(), &musicMicroservice.PlaylistPageOptions{
					PlaylistID: 2,
					UserID:     1,
				}).
					Return(&musicMicroservice.PlaylistPageResponse{
						PlaylistID: 2,
						Title:      "testTitle",
					}, nil)
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"id\":2,\"title\":\"testTitle\"}\n",
			userID:         1,
			playlistId:     2,
		},
		{
			name: "Handler returned status 400",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().PlaylistPage(gomock.Any(), &musicMicroservice.PlaylistPageOptions{
					PlaylistID: 2,
					UserID:     1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error()))
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
			userID:         1,
			playlistId:     2,
		},
		{
			name: "No RequestID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().PlaylistPage(gomock.Any(), &musicMicroservice.PlaylistPageOptions{
					PlaylistID: 2,
					UserID:     1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:    http.StatusInternalServerError,
			doNotSetRequestID: true,
		},
		{
			name: "No UserID",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().PlaylistPage(gomock.Any(), &musicMicroservice.PlaylistPageOptions{
					PlaylistID: 2,
					UserID:     1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusInternalServerError,
			doNotSerUserID: true,
		},
		{
			name: "User is unauthorized => userId = -1",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().PlaylistPage(gomock.Any(), &musicMicroservice.PlaylistPageOptions{
					PlaylistID: 2,
					UserID:     1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":401,\"message\":\"User is not authorized\"}\n",
			userID:         -1,
		},
		{
			name: "Wrong type of param in query",
			mock: func(controller *gomock.Controller) *musicMock.MockMusicClient {
				moq := musicMock.NewMockMusicClient(controller)
				moq.EXPECT().PlaylistPage(gomock.Any(), &musicMicroservice.PlaylistPageOptions{
					PlaylistID: 2,
					UserID:     1,
				}).Return(nil, status.Error(codes.InvalidArgument, errors.New("error").Error())).AnyTimes()
				return moq
			},
			expectedStatus:       http.StatusInternalServerError,
			wrongTypeOfParameter: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/api/v1/playlists/:id",
				strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			ctx.SetParamNames("id")
			if currentTest.wrongTypeOfParameter {
				ctx.SetParamValues("qwe!123scd")
			} else {
				ctx.SetParamValues(strconv.FormatInt(currentTest.playlistId, 10))
			}

			if !currentTest.doNotSetRequestID {
				ctx.Set("REQUEST_ID", "1")
			}
			if !currentTest.doNotSerUserID {
				ctx.Set("USER_ID", currentTest.userID)
			}

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			imageServices := image.NewImagesService()

			controller := gomock.NewController(t)
			musicManagerMock := currentTest.mock(controller)

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManagerMock, playlistsManager)
			if assert.NoError(t, r.GetPlaylistPage(ctx)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestAPIMicroservices_ParseErrorByCode(t *testing.T) {
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
	musicConn, _ := grpc.Dial(
		os.Getenv("MUSIC_HOST"),
		grpc.WithInsecure(),
	)

	tests := []struct {
		name           string
		expectedStatus int
		expectedJSON   string
		requestID      string
		error          error
	}{
		{
			name:           "Parse 500 error",
			expectedStatus: http.StatusInternalServerError,
			requestID: "1",
			error: status.Error(codes.Internal, errors.New("error").Error()),
		},
		{
			name:           "Parse 400 error",
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":400,\"message\":\"error\"}\n",
			requestID: "1",
			error: status.Error(codes.InvalidArgument, errors.New("error").Error()),
		},
		{
			name:              "Parse 403 error",
			expectedStatus:    http.StatusOK,
			expectedJSON:      "{\"status\":403,\"message\":\"error\"}\n",
			requestID: "1",
			error: status.Error(codes.PermissionDenied, errors.New("error").Error()),
		},
		{
			name:           "Parse 404 error",
			expectedStatus: http.StatusOK,
			expectedJSON:   "{\"status\":404,\"message\":\"error\"}\n",
			requestID: "1",
			error: status.Error(codes.NotFound, errors.New("error").Error()),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/api/v1/home/album/:id",
				strings.NewReader(""))
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)

			profileManager := profileMicroservice.NewProfileClient(profileConn)
			authManager := authMicroservice.NewAuthorizationClient(authConn)
			playlistsManager := playlistsMicroservice.NewPlaylistsClient(playlistsConn)
			musicManager := musicMicroservice.NewMusicClient(musicConn)
			imageServices := image.NewImagesService()

			r := NewAPIMicroservices(logger, imageServices, authManager, profileManager, musicManager, playlistsManager)
			if assert.NoError(t, r.ParseErrorByCode(ctx, currentTest.requestID, currentTest.error)) {
				assert.Equal(t, currentTest.expectedStatus, rec.Code)
				assert.Equal(t, currentTest.expectedJSON, rec.Body.String())
			}
		})
	}
}
