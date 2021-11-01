package delivery

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/utils/constants"
	"errors"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestUserDelivery_Logout(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	tests := []struct {
		name 			  string
		usecaseMock 	  *mock.MockUserUseCase
		cookie 			  *http.Cookie
		expectedStatus    int
		expectedJSON	  string
	}{
		{
			name: "Successfully logged out, handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				LogoutFunc: func(string) error {
					return nil
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(constants.CookieLifetime),
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":200,\"message\":\"Logged out\"}\n",
		},
		{
			name: "No cookies was set",
			usecaseMock: &mock.MockUserUseCase{},
			cookie: &http.Cookie{ },
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
		{
			name: "usecase.Logout returned error",
			usecaseMock: &mock.MockUserUseCase{
				LogoutFunc: func(string) error {
					return errors.New("error")
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(constants.CookieLifetime),
			},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/", nil)
			req.AddCookie(tt.cookie)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/user/logout")
			ctx.Set("REQUEST_ID", "1")

			r := NewUserDelivery(logger, tt.usecaseMock)
			if assert.NoError(t, r.Logout(ctx)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestUserDelivery_GetAvatarForMainPage(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	tests := []struct {
		name 			string
		usecaseMock 	*mock.MockUserUseCase
		cookie 			*http.Cookie
		userID			int
		expectedStatus 	int
		expectedJSON	string
	}{
		{
			name: "Successfully returned avatar, handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				GetAvatarFilenameFunc: func(int) (string, *models.CustomError) {
					return "avatar", nil
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(constants.CookieLifetime),
			},
			userID: 1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":200,\"avatar\":\"avatar\"}\n",
		},
		{
			name: "usecase.GetAvatarFilename returned status 500",
			usecaseMock: &mock.MockUserUseCase{
				GetAvatarFilenameFunc: func(int) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: http.StatusInternalServerError,
						OriginalError: errors.New("error"),
					}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(constants.CookieLifetime),
			},
			userID: 1,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "User is not authorized - userID = -1",
			usecaseMock: &mock.MockUserUseCase{},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(constants.CookieLifetime),
			},
			userID: -1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.AddCookie(testCase.cookie)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/auth")
			ctx.Set("REQUEST_ID", "1")
			ctx.Set("USER_ID", testCase.userID)

			r := NewUserDelivery(logger, testCase.usecaseMock)
			if assert.NoError(t, r.GetAvatarForMainPage(ctx)) {
				assert.Equal(t, testCase.expectedStatus, rec.Code)
				assert.Equal(t, testCase.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestUserDelivery_Login(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	tests := []struct {
		name 		   string
		usecaseMock    *mock.MockUserUseCase
		expectedStatus int
		expectedJSON   string
	}{
		{
			name: "Successfully signed in, Handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				LoginFunc: func(auth *models.Auth) (string, *models.CustomError) {
					return "cookieValue", nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":200,\"message\":\"User is authorized\"}\n",
		},
		{
			name: "usecase.Login returned 500 error",
			usecaseMock: &mock.MockUserUseCase{
				LoginFunc: func(auth *models.Auth) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: http.StatusInternalServerError,
						OriginalError: errors.New("error"),
					}
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "usecase.Login returned 500 error",
			usecaseMock: &mock.MockUserUseCase{
				LoginFunc: func(auth *models.Auth) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: http.StatusBadRequest,
						Message: "Bad request",
					}
				},
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":400,\"message\":\"Bad request\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/user/signin",
				strings.NewReader(`{"email": "LaHaine@gmail.com", "password": "JesusLovesMe1337"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.Set("REQUEST_ID", "1")

			r := NewUserDelivery(logger, tt.usecaseMock)
			if assert.NoError(t, r.Login(ctx)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestUserDelivery_Register(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	tests := []struct {
		name 			string
		usecaseMock 	*mock.MockUserUseCase
		expectedStatus 	int
		expectedJSON	string
	}{
		{
			name: "User signed up successfully, handler returned status 201",
			usecaseMock: &mock.MockUserUseCase{
				RegisterFunc: func(user *models.User) (string, *models.CustomError) {
					return "cookieValue", nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedJSON: "{\"status\":201,\"message\":\"User was created successfully\"}\n",
		},
		{
			name: "usecase.Register returned 500 error",
			usecaseMock: &mock.MockUserUseCase{
				RegisterFunc: func(user *models.User) (string, *models.CustomError) {
					return "", &models.CustomError {
						ErrorType: http.StatusInternalServerError,
						OriginalError: errors.New("error"),
					}
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "usecase.Register returned 400 error",
			usecaseMock: &mock.MockUserUseCase{
				RegisterFunc: func(user *models.User) (string, *models.CustomError) {
					return "", &models.CustomError {
						ErrorType: http.StatusBadRequest,
						Message: "Bad request",
					}
				},
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":400,\"message\":\"Bad request\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/user/signup", strings.NewReader(`{"email": "test.inter@ndeiud.com", "password": "jfdIHD#&n873D"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.Set("REQUEST_ID", "1")

			r := NewUserDelivery(logger, tt.usecaseMock)
			if assert.NoError(t, r.Register(ctx)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestUserDelivery_GetSettings(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	tests := []struct {
		name 			string
		usecaseMock 	*mock.MockUserUseCase
		userID 			int
		expectedStatus 	int
		expectedJSON	string
	}{
		{
			name: "Successfully got settings, handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{Email: "LaHaine@gmail.com"}, nil
				},
			},
			userID: 1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"email\":\"LaHaine@gmail.com\"}\n",
		},
		{
			name: "usecase.GetSettings returned status 500",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{
						ErrorType: http.StatusInternalServerError,
						OriginalError: errors.New("error"),
					}
				},
			},
			userID: 1,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "User is not authorized - userID = -1",
			usecaseMock: &mock.MockUserUseCase{},
			userID: -1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/user/settings")
			ctx.Set("REQUEST_ID", "1")
			ctx.Set("USER_ID", testCase.userID)

			r := NewUserDelivery(logger, testCase.usecaseMock)
			if assert.NoError(t, r.GetSettings(ctx)) {
				assert.Equal(t, testCase.expectedStatus, rec.Code)
				assert.Equal(t, testCase.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestUserDelivery_UpdateSettings(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	tests := []struct {
		name 			string
		usecaseMock 	*mock.MockUserUseCase
		userID			int
		expectedStatus 	int
		expectedJSON	string
	}{
		{
			name: "Successfully updated settings, handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return nil
				},
			},
			userID: 1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":200,\"message\":\"Settings were uploaded successfully\"}\n",
		},
		{
			name: "usecase.GetSettings returned 500 error",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{
						ErrorType: http.StatusInternalServerError,
						OriginalError: errors.New("error"),
					}
				},
			},
			userID: 1,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "usecase.UpdateSettings returned 500 error",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return &models.CustomError{
						ErrorType: http.StatusInternalServerError,
						OriginalError: errors.New("error"),
					}
				},
			},
			userID: 1,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "usecase.UpdateSettings returned 400 error",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return &models.CustomError{
						ErrorType: http.StatusBadRequest,
						Message: "Bad request",
					}
				},
			},
			userID: 1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":400,\"message\":\"Bad request\"}\n",
		},
		{
			name: "User is not authorized - userID = -1",
			usecaseMock: &mock.MockUserUseCase{},
			userID: -1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.PATCH, "/api/v1/user/settings",
				strings.NewReader(`{"email": "LaHaine@gmail.com"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.Set("USER_ID", testCase.userID)
			ctx.Set("REQUEST_ID", "1")


			r := NewUserDelivery(logger, testCase.usecaseMock)
			if assert.NoError(t, r.UpdateSettings(ctx)) {
				assert.Equal(t, testCase.expectedStatus, rec.Code)
				assert.Equal(t, testCase.expectedJSON, rec.Body.String())
			}
		})
	}
}

func TestUserDelivery_GetCsrf(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	tests := []struct {
		name 			string
		usecaseMock 	*mock.MockUserUseCase
		userID 			int
		expectedStatus 	int
		expectedJSON	string
	}{
		{
			name: "User is not authorized - userID = -1",
			usecaseMock: &mock.MockUserUseCase{},
			userID: -1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/csrf")
			ctx.Set("REQUEST_ID", "1")
			ctx.Set("USER_ID", testCase.userID)

			r := NewUserDelivery(logger, testCase.usecaseMock)
			if assert.NoError(t, r.GetCsrf(ctx)) {
				assert.Equal(t, testCase.expectedStatus, rec.Code)
				assert.Equal(t, testCase.expectedJSON, rec.Body.String())
			}
		})
	}
}
