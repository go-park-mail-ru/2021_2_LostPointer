package delivery

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
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
	usecaseMock := &mock.MockUserUseCase{
		LogoutFunc: func(string) {},
	}
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
			name: "Handler returned status 200",
			usecaseMock: usecaseMock,
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":200,\"message\":\"Logged out\"}\n",
		},
		{
			name: "Handler returned status 401, no cookies was set",
			usecaseMock: usecaseMock,
			cookie: &http.Cookie{ },
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",

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

func TestUserDelivery_IsAuthorized(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	tests := []struct {
		name 			string
		usecaseMock 	*mock.MockUserUseCase
		cookie 			*http.Cookie
		expectedStatus 	int
		expectedJSON	string
	}{
		{
			name: "Handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				IsAuthorizedFunc: func(s string) (bool, int, *models.CustomError) {
					return true, 1, nil
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":200,\"message\":\"User is authorized\"}\n",
		},
		{
			name: "Handler returned status 401, usecase.IsAuthorized returned false",
			usecaseMock: &mock.MockUserUseCase{
				IsAuthorizedFunc: func(s string) (bool, int,  *models.CustomError) {
					return false, 0, &models.CustomError{ErrorType: 401}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
		{
			name: "Handler returned status 401, no cookies was set",
			usecaseMock: &mock.MockUserUseCase{ },
			cookie: &http.Cookie{ },
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.AddCookie(tt.cookie)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/auth")
			ctx.Set("REQUEST_ID", "1")

			r := NewUserDelivery(logger, tt.usecaseMock)
			if assert.NoError(t, r.IsAuthorized(ctx)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedJSON, rec.Body.String())
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
			name: "Handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				LoginFunc: func(auth models.Auth) (string, *models.CustomError) {
					return "some_sesion_token", nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":200,\"message\":\"User is authorized\"}\n",
		},
		{
			name: "Handler returned status 400, usecase.Login returned CustomError with ErrorType = 400",
			usecaseMock: &mock.MockUserUseCase{
				LoginFunc: func(auth models.Auth) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 400,
						Message: "BadRequest",
					}
				},
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":400,\"message\":\"BadRequest\"}\n",
		},
		{
			name: "Handler returned status 500, usecase.Login returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCase{
				LoginFunc: func(auth models.Auth) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 500,
						OriginalError: errors.New("error"),
					}
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: "{\"status\":500,\"message\":\"error\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/user/signin",
				strings.NewReader(`{"email": "test.inter@ndeiud.com", "password": "jfdIHD#&n873D"}`))
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
			name: "Handler returned status 201",
			usecaseMock: &mock.MockUserUseCase{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
					return "token", nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedJSON: "{\"status\":201,\"message\":\"User was created successfully\"}\n",
		},
		{
			name: "Handler returned status 401, usecase.Login returned CustomError with ErrorType = 400",
			usecaseMock: &mock.MockUserUseCase{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 400,
						Message: "BadRequest",
					}
				},
			},
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":400,\"message\":\"BadRequest\"}\n",
		},
		{
			name: "Handler returned status 500, usecase.Login returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCase{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 500,
						OriginalError: errors.New("error"),
					}
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: "{\"status\":500,\"message\":\"error\"}\n",
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
		input 			int
		expectedStatus 	int
		expectedJSON	string
	}{
		{
			name: "Handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{Email: "alex1234@gmail.com"}, nil
				},
			},
			input: 1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"email\":\"alex1234@gmail.com\"}\n",
		},
		{
			name: "Handler returned status 401, user was not authorized",
			usecaseMock: &mock.MockUserUseCase{},
			input: 0,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
		{
			name: "Handler returned status 500, usecase.GetSettings returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{
						ErrorType: 500,
						OriginalError: errors.New("error"),
					}
				},
			},
			input: 1,
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: "{\"status\":500,\"message\":\"error\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/user/settings")
			ctx.Set("USER_ID", tt.input)
			ctx.Set("REQUEST_ID", "1")
			ctx.Set("AUTHORIZATION_ERROR", "1")

			r := NewUserDelivery(logger, tt.usecaseMock)
			if assert.NoError(t, r.GetSettings(ctx)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedJSON, rec.Body.String())
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
		input			int
		expectedStatus 	int
		expectedJSON	string
	}{
		{
			name: "Handler returned status 200",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return nil
				},
			},
			input: 1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":200,\"message\":\"Settings were uploaded successfully\"}\n",
		},
		{
			name: "Handler returned status 401, user was not authorized",
			usecaseMock: &mock.MockUserUseCase{},
			input: 0,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":401,\"message\":\"User is not authorized\"}\n",
		},
		{
			name: "Handler returned status 500, usecase.GetSettings returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{ErrorType: 500, OriginalError: errors.New("error")}
				},
			},
			input: 1,
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: "{\"status\":500,\"message\":\"error\"}\n",
		},
		{
			name: "Handler returned status 400, usecase.UpdateSettings returned CustomError with ErrorType = 400",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return &models.CustomError{
						ErrorType: 400,
						Message: "BadRequest",
					}
				},
			},
			input: 1,
			expectedStatus: http.StatusOK,
			expectedJSON: "{\"status\":400,\"message\":\"BadRequest\"}\n",
		},
		{
			name: "Handler returned status 500, usecase.UpdateSettings returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCase{
				GetSettingsFunc: func(int) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return &models.CustomError{ErrorType: 500, OriginalError: errors.New("error")}
				},
			},
			input: 1,
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: "{\"status\":500,\"message\":\"error\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.PATCH, "/api/v1/user/settings",  strings.NewReader(`{"email": "test.inter@ndeiud.com"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.Set("USER_ID", tt.input)
			ctx.Set("REQUEST_ID", "1")
			ctx.Set("AUTHORIZATION_ERROR", "1")

			r := NewUserDelivery(logger, tt.usecaseMock)
			if assert.NoError(t, r.UpdateSettings(ctx)) {
				assert.Equal(t, tt.expectedStatus, rec.Code)
				assert.Equal(t, tt.expectedJSON, rec.Body.String())
			}
		})
	}
}
