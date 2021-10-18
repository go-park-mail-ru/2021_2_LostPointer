package delivery

import (
	"2021_2_LostPointer/pkg/mock"
	"2021_2_LostPointer/pkg/models"
	"errors"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestUserDelivery_Logout(t *testing.T) {
	usecaseMock := &mock.MockUserUseCaseIFace{
		LogoutFunc: func(string) {},
	}

	tests := []struct {
		name 		string
		usecaseMock *mock.MockUserUseCaseIFace
		cookie 		*http.Cookie
		expected    int
	}{
		{
			name: "Handler returned status 200",
			usecaseMock: usecaseMock,
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusOK,
		},
		{
			name: "Handler returned status 401, no cookies was set",
			usecaseMock: usecaseMock,
			cookie: &http.Cookie{ },
			expected: http.StatusUnauthorized,
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

			r := NewUserDelivery(tt.usecaseMock)
			if assert.NoError(t, r.Logout(ctx)) {
				assert.Equal(t, tt.expected, rec.Code)
			}
		})
	}
}

func TestUserDelivery_IsAuthorized(t *testing.T) {
	tests := []struct {
		name 		string
		usecaseMock *mock.MockUserUseCaseIFace
		cookie 		*http.Cookie
		expected 	int
	}{
		{
			name: "Handler returned status 200",
			usecaseMock: &mock.MockUserUseCaseIFace{
				IsAuthorizedFunc: func(s string) (bool, *models.CustomError) {
					return true, nil
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusOK,
		},
		{
			name: "Handler returned status 401, usecase.IsAuthorized returned false",
			usecaseMock: &mock.MockUserUseCaseIFace{
				IsAuthorizedFunc: func(s string) (bool, *models.CustomError) {
					return false, nil
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Handler returned status 401, no cookies was set",
			usecaseMock: &mock.MockUserUseCaseIFace{ },
			cookie: &http.Cookie{ },
			expected: http.StatusUnauthorized,
		},
		{
			name: "Handler returned status 401, usecase.IsAuthorized returned CustomError with ErrorType = 401",
			usecaseMock: &mock.MockUserUseCaseIFace{
				IsAuthorizedFunc: func(s string) (bool, *models.CustomError) {
					return false, &models.CustomError{
						ErrorType: 401,
					}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusUnauthorized,
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

			r := NewUserDelivery(tt.usecaseMock)
			if assert.NoError(t, r.IsAuthorized(ctx)) {
				assert.Equal(t, tt.expected, rec.Code)
			}
		})
	}
}

func TestUserDelivery_Login(t *testing.T) {
	tests := []struct {
		name string
		usecaseMock *mock.MockUserUseCaseIFace
		expected int
	}{
		{
			name: "Handler returned status 200",
			usecaseMock: &mock.MockUserUseCaseIFace{
				LoginFunc: func(auth models.Auth) (string, *models.CustomError) {
					return "some_sesion_token", nil
				},
			},
			expected: http.StatusOK,
		},
		{
			name: "Handler returned status 400, usecase.Login returned CustomError with ErrorType = 400",
			usecaseMock: &mock.MockUserUseCaseIFace{
				LoginFunc: func(auth models.Auth) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 400,
					}
				},
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Handler returned status 500, usecase.Login returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCaseIFace{
				LoginFunc: func(auth models.Auth) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 500,
						OriginalError: errors.New("some_error"),
					}
				},
			},
			expected: http.StatusInternalServerError,
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
			r := NewUserDelivery(tt.usecaseMock)
			if assert.NoError(t, r.Login(ctx)) {
				assert.Equal(t, tt.expected, rec.Code)
			}
		})
	}
}

func TestUserDelivery_Register(t *testing.T) {
	tests := []struct {
		name 		string
		usecaseMock *mock.MockUserUseCaseIFace
		expected 	int
	}{
		{
			name: "Handler returned status 201",
			usecaseMock: &mock.MockUserUseCaseIFace{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
					return "token", nil
				},
			},
			expected: http.StatusCreated,
		},
		{
			name: "Handler returned status 401, usecase.Login returned CustomError with ErrorType = 400",
			usecaseMock: &mock.MockUserUseCaseIFace{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 400,
					}
				},
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Handler returned status 500, usecase.Login returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCaseIFace{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 500,
						OriginalError: errors.New("error"),
					}
				},
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/api/v1/user/signup", strings.NewReader(`{"email": "test.inter@ndeiud.com", "password": "jfdIHD#&n873D"}`))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			r := NewUserDelivery(tt.usecaseMock)
			if assert.NoError(t, r.Register(ctx)) {
				assert.Equal(t, tt.expected, rec.Code)
			}
		})
	}
}

func TestUserDelivery_GetSettings(t *testing.T) {
	tests := []struct {
		name 		string
		usecaseMock *mock.MockUserUseCaseIFace
		cookie 		*http.Cookie
		expected 	int
	}{
		{
			name: "Handler returned status 200",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusOK,
		},
		{
			name: "Handler returned status 401, no cookies was set",
			usecaseMock: &mock.MockUserUseCaseIFace{},
			cookie: &http.Cookie{},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Handler returned status 401, usecase.GetSettings returned CustomError with ErrorType = 401",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{
						ErrorType: 401,
					}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Handler returned status 500, usecase.GetSettings returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{
						ErrorType: 500,
						OriginalError: errors.New("error"),
					}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			req.AddCookie(tt.cookie)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/user/settings")

			r := NewUserDelivery(tt.usecaseMock)
			if assert.NoError(t, r.GetSettings(ctx)) {
				assert.Equal(t, tt.expected, rec.Code)
			}
		})
	}
}

func TestUserDelivery_UpdateSettings(t *testing.T) {
	tests := []struct {
		name 		string
		usecaseMock *mock.MockUserUseCaseIFace
		cookie 		*http.Cookie
		expected 	int
	}{
		{
			name: "Handler returned status 200",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(string, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return nil
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusOK,
		},
		{
			name: "Handler returned status 401, no cookies was set",
			usecaseMock: &mock.MockUserUseCaseIFace{},
			cookie: &http.Cookie{},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Handler returned status 401, usecase.GetSettings returned CustomError with ErrorType = 401",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{ErrorType: 401}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Handler returned status 500, usecase.GetSettings returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{ErrorType: 500}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "Handler returned status 400, usecase.UpdateSettings returned CustomError with ErrorType = 400",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(string, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return &models.CustomError{ErrorType: 400}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Handler returned status 500, usecase.UpdateSettings returned CustomError with ErrorType = 500",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return &models.SettingsGet{}, nil
				},
				UpdateSettingsFunc: func(string, *models.SettingsGet, *models.SettingsUpload) *models.CustomError {
					return &models.CustomError{ErrorType: 500, OriginalError: errors.New("some_error")}
				},
			},
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.PATCH, "/api/v1/user/settings",  strings.NewReader(`{"email": "test.inter@ndeiud.com"}`))
			req.AddCookie(tt.cookie)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			r := NewUserDelivery(tt.usecaseMock)
			if assert.NoError(t, r.UpdateSettings(ctx)) {
				assert.Equal(t, tt.expected, rec.Code)
			}
		})
	}
}
