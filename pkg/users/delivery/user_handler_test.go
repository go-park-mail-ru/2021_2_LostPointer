package delivery

import (
	"2021_2_LostPointer/pkg/mock"
	"2021_2_LostPointer/pkg/models"
	"errors"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
		cookie 		http.Cookie
		expected    int
	}{
		{
			name: "Successfully logged out",
			usecaseMock: usecaseMock,
			cookie: http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusOK,
		},
		{
			name: "User was not authorized, no cookies was set",
			usecaseMock: usecaseMock,
			cookie: http.Cookie{ },
			expected: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.POST, "/", nil)
			req.AddCookie(&tt.cookie)
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
		cookie 		http.Cookie
		expected 	int
	}{
		{
			name: "User is authorized",
			usecaseMock: &mock.MockUserUseCaseIFace{
				IsAuthorizedFunc: func(s string) (bool, error) {
					return true, nil
				},
			},
			cookie: http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusOK,
		},
		{
			name: "User is not authorized",
			usecaseMock: &mock.MockUserUseCaseIFace{
				IsAuthorizedFunc: func(s string) (bool, error) {
					return false, nil
				},
			},
			cookie: http.Cookie{
				Name:     "Session_cookie",
				Value:    "Cookie_value",
				Expires:  time.Now().Add(cookieLifetime),
			},
			expected: http.StatusUnauthorized,
		},
		{
			name: "User is not authorized, no cookies set",
			usecaseMock: &mock.MockUserUseCaseIFace{ },
			cookie: http.Cookie{ },
			expected: http.StatusUnauthorized,
		},
		{
			name: "User is not authorized, no session in redis",
			usecaseMock: &mock.MockUserUseCaseIFace{
				IsAuthorizedFunc: func(s string) (bool, error) {
					return false, errors.New("no_session_in_redis")
				},
			},
			cookie: http.Cookie{
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
			req.AddCookie(&tt.cookie)
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
				LoginFunc: func(auth models.Auth) (string, error) {
					return "some_sesion_token", nil
				},
			},
			expected: http.StatusOK,
		},
		{
			name: "Unsuccessful log in, usecase Login returns empty token",
			usecaseMock: &mock.MockUserUseCaseIFace{
				LoginFunc: func(auth models.Auth) (string, error) {
					return "", nil
				},
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Unsuccessful log in, usecase Login returns error",
			usecaseMock: &mock.MockUserUseCaseIFace{
				LoginFunc: func(auth models.Auth) (string, error) {
					return "", errors.New("some_error_in_login")
				},
			},
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/user/signin")
			ctx.Set("email", "alex1234@gmail.com")
			ctx.Set("password", "Alexey123456!")

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
				RegisterFunc: func(user models.User) (string, string, error) {
					return "some_session_token", "", nil
				},
			},
			expected: http.StatusCreated,
		},
		{
			name: "Unsuccessful register, usecase Register returns error",
			usecaseMock: &mock.MockUserUseCaseIFace{
				RegisterFunc: func(user models.User) (string, string, error) {
					return "", "", errors.New("some_error_in_register")
				},
			},
			expected: http.StatusInternalServerError,
		},
		{
			name: "Unsuccessful register, usecase Register returns empty token",
			usecaseMock: &mock.MockUserUseCaseIFace{
				RegisterFunc: func(user models.User) (string, string, error) {
					return "", "", nil
				},
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := echo.New()
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			ctx := server.NewContext(req, rec)
			ctx.SetPath("/api/v1/user/signup")
			ctx.Set("email", "alex1234@gmail.com")
			ctx.Set("password", "Alexey123456!")
			ctx.Set("nickname", "Alexey_Kosenko")

			r := NewUserDelivery(tt.usecaseMock)
			if assert.NoError(t, r.Register(ctx)) {
				assert.Equal(t, tt.expected, rec.Code)
			}
		})
	}
}
