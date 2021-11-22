package usecase

import (
	"2021_2_LostPointer/internal/constants"
	customErrors "2021_2_LostPointer/internal/errors"
	"2021_2_LostPointer/internal/microservices/profile/mock"
	"2021_2_LostPointer/internal/microservices/profile/proto"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestProfileService_GetSettings(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockUserSettingsStorage
		input       *proto.GetSettingsOptions
		expected    *proto.UserSettings
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockUserSettingsStorage{
				GetSettingsFunc: func(int64) (*proto.UserSettings, error) {
					return &proto.UserSettings{}, nil
				},
			},
			input:    &proto.GetSettingsOptions{ID: 1},
			expected: &proto.UserSettings{},
		},
		{
			name: "Error 404. User not found",
			storageMock: &mock.MockUserSettingsStorage{
				GetSettingsFunc: func(int64) (*proto.UserSettings, error) {
					return nil, customErrors.ErrUserNotFound
				},
			},
			input:       &proto.GetSettingsOptions{ID: 1},
			expectedErr: true,
			err:         status.Error(codes.NotFound, customErrors.ErrUserNotFound.Error()),
		},
		{
			name: "Error 500. mock.GetSettings returned error",
			storageMock: &mock.MockUserSettingsStorage{
				GetSettingsFunc: func(int64) (*proto.UserSettings, error) {
					return nil, errors.New("error")
				},
			},
			input:       &proto.GetSettingsOptions{ID: 1},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewProfileService(test.storageMock)

			res, err := storage.GetSettings(context.Background(), test.input)
			if test.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, test.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, res)
			}
		})
	}
}

func TestProfileService_UpdateSettings(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockUserSettingsStorage
		input       *proto.UpdateSettingsOptions
		expected    *proto.EmptyProfile
		expectedErr bool
		err         error
	}{
		//---------EMAIL---------
		{
			name: "Successfully updated email",
			storageMock: &mock.MockUserSettingsStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateEmailFunc: func(int64, string) error {
					return nil
				},
			},
			input: &proto.UpdateSettingsOptions{
				Email: "lahaine@gmail.com",
				OldSettings: &proto.UserSettings{
					Email: "lannister@mercy.com",
				},
			},
			expected: &proto.EmptyProfile{},
		},
		{
			name:        "Error 400. New email is not valid",
			storageMock: &mock.MockUserSettingsStorage{},
			input: &proto.UpdateSettingsOptions{
				Email: "lahainegmail.com",
				OldSettings: &proto.UserSettings{
					Email: "lannister@mercy.com",
				},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.EmailInvalidSyntaxMessage),
		},
		{
			name: "Error 400. New email is not unique",
			storageMock: &mock.MockUserSettingsStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			input: &proto.UpdateSettingsOptions{
				Email: "lahaine@gmail.com",
				OldSettings: &proto.UserSettings{
					Email: "lannister@mercy.com",
				},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.EmailNotUniqueMessage),
		},
		{
			name: "Error 500. mock.IsEmailUnique returned error",
			storageMock: &mock.MockUserSettingsStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.UpdateSettingsOptions{
				Email: "lahaine@gmail.com",
				OldSettings: &proto.UserSettings{
					Email: "lannister@mercy.com",
				},
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.UpdateEmail returned error",
			storageMock: &mock.MockUserSettingsStorage{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateEmailFunc: func(int64, string) error {
					return errors.New("error")
				},
			},
			input: &proto.UpdateSettingsOptions{
				Email: "lahaine@gmail.com",
				OldSettings: &proto.UserSettings{
					Email: "lannister@mercy.com",
				},
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},

		//---------NICKNAME---------
		{
			name: "Successfully updated nickname",
			storageMock: &mock.MockUserSettingsStorage{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateNicknameFunc: func(int64, string) error {
					return nil
				},
			},
			input: &proto.UpdateSettingsOptions{
				Nickname: "LaHaine",
				OldSettings: &proto.UserSettings{
					Email: "AragorN",
				},
			},
			expected: &proto.EmptyProfile{},
		},
		{
			name:        "Error 400. New nickname is not valid",
			storageMock: &mock.MockUserSettingsStorage{},
			input: &proto.UpdateSettingsOptions{
				Nickname: "LaHaine!",
				OldSettings: &proto.UserSettings{
					Email: "AragorN",
				},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.NicknameInvalidSyntaxMessage),
		},
		{
			name: "Error 400. New nickname is not unique",
			storageMock: &mock.MockUserSettingsStorage{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			input: &proto.UpdateSettingsOptions{
				Nickname: "LaHaine",
				OldSettings: &proto.UserSettings{
					Email: "AragorN",
				},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.NicknameNotUniqueMessage),
		},
		{
			name: "Error 500. mock.IsNicknameUnique returned error",
			storageMock: &mock.MockUserSettingsStorage{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.UpdateSettingsOptions{
				Nickname: "LaHaine",
				OldSettings: &proto.UserSettings{
					Email: "AragorN",
				},
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.UpdateNickname returned error",
			storageMock: &mock.MockUserSettingsStorage{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateNicknameFunc: func(int64, string) error {
					return errors.New("error")
				},
			},
			input: &proto.UpdateSettingsOptions{
				Nickname: "LaHaine",
				OldSettings: &proto.UserSettings{
					Email: "AragorN",
				},
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},

		//---------AVATAR---------
		{
			name: "Successfully updated avatar",
			storageMock: &mock.MockUserSettingsStorage{
				UpdateAvatarFunc: func(int64, string) error {
					return nil
				},
			},
			input: &proto.UpdateSettingsOptions{
				AvatarFilename: "lahaine_mercy.webp",
				OldSettings:    &proto.UserSettings{},
			},
			expected: &proto.EmptyProfile{},
		},
		{
			name: "Error 500. mock.UpdateAvatar returned error",
			storageMock: &mock.MockUserSettingsStorage{
				UpdateAvatarFunc: func(int64, string) error {
					return errors.New("error")
				},
			},
			input: &proto.UpdateSettingsOptions{
				AvatarFilename: "lahaine_mercy.webp",
				OldSettings:    &proto.UserSettings{},
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},

		//---------PASSWORD---------
		{
			name: "Successfully updated password",
			storageMock: &mock.MockUserSettingsStorage{
				CheckPasswordByUserIDFunc: func(int64, string) (bool, error) {
					return true, nil
				},
				UpdatePasswordFunc: func(int64, string) error {
					return nil
				},
			},
			input: &proto.UpdateSettingsOptions{
				NewPassword: "Avt8430056!",
				OldPassword: "Avt8430066!",
				OldSettings: &proto.UserSettings{},
			},
			expected: &proto.EmptyProfile{},
		},
		{
			name:        "Error 400. OldPassword field is empty",
			storageMock: &mock.MockUserSettingsStorage{},
			input: &proto.UpdateSettingsOptions{
				NewPassword: "Avt8430056!",
				OldPassword: "",
				OldSettings: &proto.UserSettings{},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.OldPasswordFieldIsEmptyMessage),
		},
		{
			name:        "Error 400. NewPassword field is empty",
			storageMock: &mock.MockUserSettingsStorage{},
			input: &proto.UpdateSettingsOptions{
				NewPassword: "",
				OldPassword: "Avt8430055!",
				OldSettings: &proto.UserSettings{},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.NewPasswordFieldIsEmptyMessage),
		},
		{
			name: "Error 400. mock.CheckPasswordByUserID returned error \"Wrong credentials\"",
			storageMock: &mock.MockUserSettingsStorage{
				CheckPasswordByUserIDFunc: func(int64, string) (bool, error) {
					return false, customErrors.ErrWrongCredentials
				},
			},
			input: &proto.UpdateSettingsOptions{
				NewPassword: "Avt8430056",
				OldPassword: "Avt8430055!",
				OldSettings: &proto.UserSettings{},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.WrongPasswordMessage),
		},
		{
			name: "Error 500. mock.CheckPasswordByUserID returned error",
			storageMock: &mock.MockUserSettingsStorage{
				CheckPasswordByUserIDFunc: func(int64, string) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.UpdateSettingsOptions{
				NewPassword: "Avt8430056",
				OldPassword: "Avt8430055!",
				OldSettings: &proto.UserSettings{},
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. OldPassword is wrong",
			storageMock: &mock.MockUserSettingsStorage{
				CheckPasswordByUserIDFunc: func(int64, string) (bool, error) {
					return false, nil
				},
			},
			input: &proto.UpdateSettingsOptions{
				NewPassword: "Avt8430056",
				OldPassword: "Avt8430055!",
				OldSettings: &proto.UserSettings{},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.WrongPasswordMessage),
		},
		{
			name: "Error 400. New password is not valid",
			storageMock: &mock.MockUserSettingsStorage{
				CheckPasswordByUserIDFunc: func(int64, string) (bool, error) {
					return true, nil
				},
			},
			input: &proto.UpdateSettingsOptions{
				NewPassword: "Av",
				OldPassword: "Avt8430055!",
				OldSettings: &proto.UserSettings{},
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, "Password must contain at least 8 characters"),
		},
		{
			name: "Error 500. mock.UpdatePassword returned error",
			storageMock: &mock.MockUserSettingsStorage{
				CheckPasswordByUserIDFunc: func(int64, string) (bool, error) {
					return true, nil
				},
				UpdatePasswordFunc: func(int64, string) error {
					return errors.New("error")
				},
			},
			input: &proto.UpdateSettingsOptions{
				NewPassword: "Avt8430056!",
				OldPassword: "Avt8430055!",
				OldSettings: &proto.UserSettings{},
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := NewProfileService(test.storageMock)

			res, err := storage.UpdateSettings(context.Background(), test.input)
			if test.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, test.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, res)
			}
		})
	}
}
