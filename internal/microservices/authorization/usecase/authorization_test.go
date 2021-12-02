package usecase

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	customErrors "2021_2_LostPointer/internal/errors"
	"2021_2_LostPointer/internal/microservices/authorization/mock"
	"2021_2_LostPointer/internal/microservices/authorization/proto"
)

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockAuthStorage
		input       *proto.AuthData
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockAuthStorage{
				GetUserByPasswordFunc: func(*proto.AuthData) (int64, error) {
					return 1, nil
				},
				CreateSessionFunc: func(int64, string) error {
					return nil
				},
			},
			input: &proto.AuthData{Email: "lahaine@gmail.com", Password: "Avt8430066!"},
		},
		{
			name: "Error 500. mock.CreateSession returned error",
			storageMock: &mock.MockAuthStorage{
				GetUserByPasswordFunc: func(*proto.AuthData) (int64, error) {
					return 1, nil
				},
				CreateSessionFunc: func(int64, string) error {
					return errors.New("error")
				},
			},
			input:       &proto.AuthData{Email: "lahaine@gmail.com", Password: "Avt8430066!"},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. mock.GetUserByPassword returned error \"Wrong credentials\"",
			storageMock: &mock.MockAuthStorage{
				GetUserByPasswordFunc: func(*proto.AuthData) (int64, error) {
					return 1, customErrors.ErrWrongCredentials
				},
			},
			input:       &proto.AuthData{Email: "lahaine@gmail.com", Password: "Avt8430066!"},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, customErrors.ErrWrongCredentials.Error()),
		},
		{
			name: "Error 500. mock.GetUserByPassword returned error",
			storageMock: &mock.MockAuthStorage{
				GetUserByPasswordFunc: func(*proto.AuthData) (int64, error) {
					return 1, errors.New("error")
				},
			},
			input:       &proto.AuthData{Email: "lahaine@gmail.com", Password: "Avt8430066!"},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewAuthService(currentTest.storageMock)

			res, err := storage.Login(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
			} else {
				assert.True(t, len(res.Cookies) != 0)
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockAuthStorage
		input       *proto.RegisterData
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockAuthStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(*proto.RegisterData) (int64, error) {
					return 1, nil
				},
				CreateSessionFunc: func(int64, string) error {
					return nil
				},
			},
			input: &proto.RegisterData{Email: "lahaine@gmail.com", Password: "Avt8430066!", Nickname: "LaHaine"},
		},
		{
			name: "Error 500. mock.IsEmailUnique returned error",
			storageMock: &mock.MockAuthStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, errors.New("error")
				},
			},
			input:       &proto.RegisterData{Email: "lahaine@gmail.com", Password: "Avt8430066!", Nickname: "LaHaine"},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. mock.IsEmailUnique returned false",
			storageMock: &mock.MockAuthStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			input:       &proto.RegisterData{Email: "lahaine@gmail.com", Password: "Avt8430066!", Nickname: "LaHaine"},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.EmailNotUniqueMessage),
		},
		{
			name: "Error 500. mock.IsNicknameUnique returned error",
			storageMock: &mock.MockAuthStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, errors.New("error")
				},
			},
			input:       &proto.RegisterData{Email: "lahaine@gmail.com", Password: "Avt8430066!", Nickname: "LaHaine"},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. mock.IsNicknameUnique returned false",
			storageMock: &mock.MockAuthStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			input:       &proto.RegisterData{Email: "lahaine@gmail.com", Password: "Avt8430066!", Nickname: "LaHaine"},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.NicknameNotUniqueMessage),
		},
		{
			name: "Error 400. Credentials are not valid",
			storageMock: &mock.MockAuthStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
			},
			input:       &proto.RegisterData{Email: "lahainegmail.com", Password: "Avt8430066!", Nickname: "LaHaine"},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.EmailInvalidSyntaxMessage),
		},
		{
			name: "Error 500. mock.CreateUser returned error",
			storageMock: &mock.MockAuthStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(*proto.RegisterData) (int64, error) {
					return 0, errors.New("error")
				},
			},
			input:       &proto.RegisterData{Email: "lahaine@gmail.com", Password: "Avt8430066!", Nickname: "LaHaine"},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},

		{
			name: "Error 500. mock.CreateSession returned error",
			storageMock: &mock.MockAuthStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(*proto.RegisterData) (int64, error) {
					return 1, nil
				},
				CreateSessionFunc: func(int64, string) error {
					return errors.New("error")
				},
			},
			input:       &proto.RegisterData{Email: "lahaine@gmail.com", Password: "Avt8430066!", Nickname: "LaHaine"},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewAuthService(currentTest.storageMock)

			res, err := storage.Register(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.True(t, len(res.Cookies) != 0)
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_GetAvatar(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockAuthStorage
		input       *proto.UserID
		expected    *proto.Avatar
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockAuthStorage{
				GetAvatarFunc: func(int64) (string, error) {
					return "lahaine", nil
				},
			},
			input:    &proto.UserID{ID: 1},
			expected: &proto.Avatar{Filename: os.Getenv("USERS_ROOT_PREFIX") + "lahaine" + constants.UserAvatarExtension150px},
		},
		{
			name: "Error 500. mock.GetAvatar returned error",
			storageMock: &mock.MockAuthStorage{
				GetAvatarFunc: func(int64) (string, error) {
					return "", errors.New("error")
				},
			},
			input:       &proto.UserID{ID: 1},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewAuthService(currentTest.storageMock)

			res, err := storage.GetAvatar(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockAuthStorage
		input       *proto.Cookie
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockAuthStorage{
				DeleteSessionFunc: func(string) error {
					return nil
				},
			},
			input: &proto.Cookie{Cookies: "cookie"},
		},
		{
			name: "Error 500. mock.DeleteSession returned error",
			storageMock: &mock.MockAuthStorage{
				DeleteSessionFunc: func(string) error {
					return errors.New("error")
				},
			},
			input:       &proto.Cookie{Cookies: "cookie"},
			expectedErr: true,
			err:         status.Error(codes.Internal, "rpc error: code = Internal desc = error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewAuthService(currentTest.storageMock)

			_, err := storage.Logout(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_GetUserByCookie(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockAuthStorage
		input       *proto.Cookie
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockAuthStorage{
				GetUserByCookieFunc: func(string) (int64, error) {
					return 1, nil
				},
			},
			input: &proto.Cookie{Cookies: "cookie"},
		},
		{
			name: "Error 500. mock.GetUserByCookie returned error",
			storageMock: &mock.MockAuthStorage{
				GetUserByCookieFunc: func(string) (int64, error) {
					return 0, errors.New("error")
				},
			},
			input:       &proto.Cookie{Cookies: "cookie"},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewAuthService(currentTest.storageMock)

			_, err := storage.GetUserByCookie(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
