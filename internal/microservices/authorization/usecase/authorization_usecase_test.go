package usecase

import (
	"2021_2_LostPointer/internal/constants"
	session "2021_2_LostPointer/internal/microservices/authorization/delivery"
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestAuthorizationUseCase_GetUserBySession(t *testing.T) {
	type response struct {
		userID *session.UserID
		error error
	}
	tests := []struct {
		name 		string
		mockDB 		*mock.MockUserRepository
		mockSession *mock.MockSessionRepository
		input 		*session.SessionData
		expected 	response
		expectedErr bool
	}{
		{
			name: "Successfully returned user id",
			mockDB: &mock.MockUserRepository{},
			mockSession: &mock.MockSessionRepository{
				GetUserIdByCookieFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			input: &session.SessionData{
				Cookies: "cookie_value",
			},
			expected: response{
				userID: &session.UserID{
					UserID: 1,
				},
				error: nil,
			},
		},
		{
			name: "Session not found, userID = 0",
			mockDB: &mock.MockUserRepository{},
			mockSession: &mock.MockSessionRepository{
				GetUserIdByCookieFunc: func(string) (int, error) {
					return 0, nil
				},
			},
			input: &session.SessionData{
				Cookies: "cookie_value",
			},
			expected: response{
				userID: &session.UserID{
					UserID: -1,
				},
				error: nil,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewAuthorizationUseCase(testCase.mockDB, testCase.mockSession)

			got := response{}

			got.userID, got.error = r.GetUserBySession(context.Background(), testCase.input)
			if testCase.expectedErr {
				assert.NotNil(t, got.error)
			} else {
				assert.Nil(t, got.error)
				assert.Equal(t, testCase.expected.userID.UserID, got.userID.UserID)
			}
		})
	}
}

func TestAuthorizationUseCase_DeleteSession(t *testing.T) {
	tests := []struct {
		name 			string
		mockDB 			*mock.MockUserRepository
		mockSession 	*mock.MockSessionRepository
		input 			*session.SessionData
		expectedErr  	bool
	}{
		{
			name: "Successfully deleted session",
			mockDB: &mock.MockUserRepository{},
			mockSession: &mock.MockSessionRepository{
				DeleteSessionFunc: func(string) error {
					return nil
				},
			},
			input: &session.SessionData{
				Cookies: "cookie_value",
			},

		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewAuthorizationUseCase(testCase.mockDB, testCase.mockSession)


			_, err := r.DeleteSession(context.Background(), testCase.input)
			if testCase.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthorizationUseCase_SignIn(t *testing.T) {
	type response struct {
		sessionData *session.SessionData
		error 		error
	}

	tests := []struct {
		name 			string
		mockDB 			*mock.MockUserRepository
		mockSession 	*mock.MockSessionRepository
		input 			*session.Auth
		expected 		response
		expectedErr  	bool
	}{
		{
			name: "Successfully signed in",
			mockDB: &mock.MockUserRepository {
				DoesUserExistFunc: func(*models.Auth) (int, error) {
					return 1, nil
				},
			},
			mockSession: &mock.MockSessionRepository {
				CreateSessionFunc: func(int, string) error {
					return nil
				},
			},
			input: &session.Auth {
				Login: "LaHaine@gmail.com",
				Password: "JesusLovesMe",
			},
		},
		{
			name: "db.DoesUserExist returned error",
			mockDB: &mock.MockUserRepository {
				DoesUserExistFunc: func(*models.Auth) (int, error) {
					return 0, errors.New("error")
				},
			},
			mockSession: &mock.MockSessionRepository {},
			input: &session.Auth {
				Login: "LaHaine@gmail.com",
				Password: "JesusLovesMe",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Internal, "error"),
			},
			expectedErr: true,
		},
		{
			name: "db.DoesUserExist returned userID = 0",
			mockDB: &mock.MockUserRepository {
				DoesUserExistFunc: func(*models.Auth) (int, error) {
					return 0, nil
				},
			},
			mockSession: &mock.MockSessionRepository {},
			input: &session.Auth {
				Login: "LaHaine@gmail.com",
				Password: "JesusLovesMe",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Aborted, constants.WrongCredentials),
			},
			expectedErr: true,
		},
		{
			name: "sessions.CreateSession returned error",
			mockDB: &mock.MockUserRepository {
				DoesUserExistFunc: func(*models.Auth) (int, error) {
					return 1, nil
				},
			},
			mockSession: &mock.MockSessionRepository {
				CreateSessionFunc: func(int, string) error {
					return errors.New("error")
				},
			},
			input: &session.Auth {
				Login: "LaHaine@gmail.com",
				Password: "JesusLovesMe",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Internal, "error"),
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewAuthorizationUseCase(testCase.mockDB, testCase.mockSession)

			got := response{}

			got.sessionData, got.error = r.SignIn(context.Background(), testCase.input)
			if testCase.expectedErr {
				assert.NotNil(t, got.error)
				assert.Equal(t, testCase.expected, got)
			} else {
				assert.Nil(t, got.error)
				assert.True(t, len(got.sessionData.Cookies) > 0)
			}
		})
	}
}

func TestAuthorizationUseCase_Signup(t *testing.T) {
	type response struct {
		sessionData *session.SessionData
		error 		error
	}

	tests := []struct {
		name 			string
		mockDB 			*mock.MockUserRepository
		mockSession 	*mock.MockSessionRepository
		input 			*session.SignUpData
		expected 		response
		expectedErr  	bool
	}{
		{
			name: "Successfully signed up",
			mockDB: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(*models.User) (int, error) {
					return 1, nil
				},
			},
			mockSession: &mock.MockSessionRepository{
				CreateSessionFunc: func(int, string) error {
					return nil
				},
			},
			input: &session.SignUpData{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337!",
				Nickname: "LaHaine",
			},
			expected: response{
				sessionData: &session.SessionData{
					Cookies: "some_cookies",
				},
				error: nil,
			},
		},
		{
			name: "db.IsEmailUnique returned error",
			mockDB: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, errors.New("error")
				},
			},
			mockSession: &mock.MockSessionRepository{},
			input: &session.SignUpData{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337!",
				Nickname: "LaHaine",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Internal, "error"),
			},
			expectedErr: true,
		},
		{
			name: "db.IsEmailUnique returned false",
			mockDB: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			mockSession: &mock.MockSessionRepository{},
			input: &session.SignUpData{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337!",
				Nickname: "LaHaine",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Aborted, constants.NotUniqueEmailMessage),
			},
			expectedErr: true,
		},
		{
			name: "db.IsNicknameUnique returned error",
			mockDB: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, errors.New("error")
				},
			},
			mockSession: &mock.MockSessionRepository{},
			input: &session.SignUpData{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337!",
				Nickname: "LaHaine",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Internal, "error"),
			},
			expectedErr: true,
		},
		{
			name: "db.IsNicknameUnique returned false",
			mockDB: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			mockSession: &mock.MockSessionRepository{},
			input: &session.SignUpData{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337!",
				Nickname: "LaHaine",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Aborted, constants.NotUniqueNicknameMessage),
			},
			expectedErr: true,
		},
		{
			name: "Not valid credentials",
			mockDB: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
			},
			mockSession: &mock.MockSessionRepository{},
			input: &session.SignUpData{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337",
				Nickname: "LaHaine",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Aborted, constants.PasswordValidationNoSpecialSymbolMessage),
			},
			expectedErr: true,
		},
		{
			name: "db.CreateUser returned error",
			mockDB: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(*models.User) (int, error) {
					return 0, errors.New("error")
				},
			},
			mockSession: &mock.MockSessionRepository{},
			input: &session.SignUpData{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337!",
				Nickname: "LaHaine",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Internal, "error"),
			},
			expectedErr: true,
		},
		{
			name: "sessions.CreateSession returned error",
			mockDB: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(*models.User) (int, error) {
					return 1, nil
				},
			},
			mockSession: &mock.MockSessionRepository{
				CreateSessionFunc: func(int, string) error {
					return errors.New("error")
				},
			},
			input: &session.SignUpData{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337!",
				Nickname: "LaHaine",
			},
			expected: response{
				sessionData: nil,
				error: status.Error(codes.Internal, "error"),
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewAuthorizationUseCase(testCase.mockDB, testCase.mockSession)

			got := response{}

			got.sessionData, got.error = r.Signup(context.Background(), testCase.input)
			if testCase.expectedErr {
				assert.NotNil(t, got.error)
				assert.Equal(t, testCase.expected, got)
			} else {
				assert.Nil(t, got.error)
				assert.True(t, len(got.sessionData.Cookies) > 0)
			}
		})
	}
}
