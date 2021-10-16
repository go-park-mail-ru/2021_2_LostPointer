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
			name: "Successfully logged out",
			usecaseMock: usecaseMock,
			cookie: &http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusOK,
		},
		{
			name: "User was not authorized, no cookies was set",
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
			name: "User is authorized",
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
			name: "User is not authorized",
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
			name: "User is not authorized, no cookies set",
			usecaseMock: &mock.MockUserUseCaseIFace{ },
			cookie: &http.Cookie{ },
			expected: http.StatusUnauthorized,
		},
		{
			name: "User is not authorized, no session in redis",
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
			name: "Successfully logged in",
			usecaseMock: &mock.MockUserUseCaseIFace{
				LoginFunc: func(auth models.Auth) (string, *models.CustomError) {
					return "some_sesion_token", nil
				},
			},
			expected: http.StatusOK,
		},
		{
			name: "Unsuccessful log in, BadRequest",
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
			name: "Unsuccessful log in, InternalServerError",
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
			req := httptest.NewRequest(echo.POST, "/api/v1/user/signin",  strings.NewReader(`{"email": "test.inter@ndeiud.com", "password": "jfdIHD#&n873D"}`))
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
			name: "Successful register",
			usecaseMock: &mock.MockUserUseCaseIFace{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
					return "some_session_token", nil
				},
			},
			expected: http.StatusCreated,
		},
		{
			name: "Unsuccessful register, BadRequest",
			usecaseMock: &mock.MockUserUseCaseIFace{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
					return "", &models.CustomError{
						ErrorType: 400,
						OriginalError: nil,
					}
				},
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Unsuccessful register, InternalServerError",
			usecaseMock: &mock.MockUserUseCaseIFace{
				RegisterFunc: func(user models.User) (string, *models.CustomError) {
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
			name: "Successfully returns settings",
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
			name: "Unsuccessfully returns settings, no session_cookie set",
			usecaseMock: &mock.MockUserUseCaseIFace{},
			cookie: &http.Cookie{},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Unsuccessfully returns settings, user is unauthorized",
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
			name: "Unsuccessfully returns settings, internal server error",
			usecaseMock: &mock.MockUserUseCaseIFace{
				GetSettingsFunc: func(string) (*models.SettingsGet, *models.CustomError) {
					return nil, &models.CustomError{
						ErrorType: 500,
						OriginalError: errors.New("some_error"),
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
			name: "Successfully updated settings",
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
			name: "Unsuccessfully updated settings, no cookies set",
			usecaseMock: &mock.MockUserUseCaseIFace{},
			cookie: &http.Cookie{},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Unsuccessfully updated settings, user not authorized",
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
			name: "Unsuccessfully updated settings, internal server error in GetSettings",
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
			name: "Unsuccessfully updated settings, bad request in UpdateSettings",
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
			name: "Unsuccessfully updated settings, internal server error in UpdateSettings",
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
