package usecase

import (
	session "2021_2_LostPointer/internal/microservices/authorization/delivery"
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/utils/constants"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestUserUseCase_Login(t *testing.T) {
	type response struct {
		cookieValue string
		err 		*models.CustomError
	}

	tests := []struct {
		name 	  	  		 string
		dbMock 	  	  		 *mock.MockUserRepository
		sessionCheckerMock   *mock.MockSessionCheckerClient
		imageMock			 *mock.MockAvatarRepositoryIFace
		input 	  	  		 *models.Auth
		expected  	  		 response
		expectedErr   		 bool
	}{
		{
			name: "Successful login",
			sessionCheckerMock: &mock.MockSessionCheckerClient{
				SignInFunc: func(context.Context, *session.Auth, ...grpc.CallOption) (*session.SessionData, error) {
					return &session.SessionData{
						Cookies: "cookieValue",
					}, nil
				},
			},
			input: &models.Auth{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe",
			},
			expected: response{
				cookieValue: "cookieValue",
				err: nil,
			},
		},
		{
			name: "Invalid credentials",
			sessionCheckerMock: &mock.MockSessionCheckerClient{
				SignInFunc: func(context.Context, *session.Auth, ...grpc.CallOption) (*session.SessionData, error) {
					return nil, status.Error(codes.Aborted, constants.WrongCredentials)
				},
			},
			input: &models.Auth{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe",
			},
			expected: response{
				cookieValue: "",
				err: &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.WrongCredentials},
			},
			expectedErr: true,
		},
		{
			name: "Internal server error in microservice",
			sessionCheckerMock: &mock.MockSessionCheckerClient{
				SignInFunc: func(context.Context, *session.Auth, ...grpc.CallOption) (*session.SessionData, error) {
					return nil, status.Error(codes.Internal, "error")
				},
			},
			input: &models.Auth{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe",
			},
			expected: response{
				cookieValue: "",
				err: &models.CustomError{
					ErrorType: http.StatusInternalServerError,
					OriginalError: status.Error(codes.Internal, "error"),
				},
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.sessionCheckerMock, testCase.imageMock)
			got := response{}

			got.cookieValue, got.err = r.Login(testCase.input)
			if testCase.expectedErr {
				assert.NotNil(t, got.err)
				assert.Equal(t, testCase.expected, got)
			} else {
				assert.Nil(t, got.err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserUseCase_Logout(t *testing.T) {
	dbMock := &mock.MockUserRepository{}
	sessionCheckerMock := &mock.MockSessionCheckerClient{
		DeleteSessionFunc: func(context.Context, *session.SessionData, ...grpc.CallOption) (*session.Empty, error) {
			return &session.Empty{}, nil
		},
	}
	imageMock := &mock.MockAvatarRepositoryIFace{}

	r := NewUserUserCase(dbMock, sessionCheckerMock, imageMock)
	err := r.Logout("alexei_kosenkov")
	assert.Nil(t, err)
}

func TestUserUseCase_Register(t *testing.T) {
	type response struct {
		cookieValue  string
		err 		 *models.CustomError
	}

	tests := []struct {
		name 	  	  		 string
		dbMock 	  	  		 *mock.MockUserRepository
		sessionCheckerMock   *mock.MockSessionCheckerClient
		imagesMock			 *mock.MockAvatarRepositoryIFace
		input 	  	  		 *models.User
		expected  	  		 response
		expectedErr   		 bool
	}{
		{
			name: "Successfully signed up",
			sessionCheckerMock: &mock.MockSessionCheckerClient{
				SignupFunc: func(context.Context, *session.SignUpData, ...grpc.CallOption) (*session.SessionData, error) {
					return &session.SessionData{
						Cookies: "cookieValue",
					}, nil
				},
			},
			input: &models.User{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337",
				Nickname: "LaHaine",
			},
			expected: response{
				cookieValue: "cookieValue",
				err: nil,
			},
		},
		{
			name: "Invalid credentials",
			sessionCheckerMock: &mock.MockSessionCheckerClient{
				SignupFunc: func(context.Context, *session.SignUpData, ...grpc.CallOption) (*session.SessionData, error) {
					return nil, status.Error(codes.Aborted, constants.NotUniqueNicknameMessage)
				},
			},
			input: &models.User{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337",
				Nickname: "LaHaine",
			},
			expected: response{
				cookieValue: "",
				err: &models.CustomError{ErrorType: http.StatusBadRequest, Message: constants.NotUniqueNicknameMessage},
			},
			expectedErr: true,
		},
		{
			name: "Internal server error in microservice",
			sessionCheckerMock: &mock.MockSessionCheckerClient{
				SignupFunc: func(context.Context, *session.SignUpData, ...grpc.CallOption) (*session.SessionData, error) {
					return nil, status.Error(codes.Internal, "error")
				},
			},
			input: &models.User{
				Email: "LaHaine@gmail.com",
				Password: "JesusLovesMe1337",
				Nickname: "LaHaine",
			},
			expected: response{
				cookieValue: "",
				err: &models.CustomError{ErrorType: http.StatusInternalServerError, OriginalError: status.Error(codes.Internal, "error")},
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.sessionCheckerMock, testCase.imagesMock)
			got := response{}

			got.cookieValue, got.err = r.Register(testCase.input)

			if testCase.expectedErr {
				assert.NotNil(t, got.err)
				assert.Equal(t, testCase.expected, got)
			} else {
				assert.Nil(t, got.err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserUseCase_GetSettings(t *testing.T) {
	type response struct {
		settings *models.SettingsGet
		err 	 *models.CustomError
	}

	tests := []struct {
		name 	  	  		 string
		dbMock 	  	  		 *mock.MockUserRepository
		sessionCheckerMock   *mock.MockSessionCheckerClient
		imagesMock			 *mock.MockAvatarRepositoryIFace
		input 	  	  		 int
		expected  	  		 response
		expectedErr   		 bool
	}{
		{
			name: "Successfully got settings",
			dbMock: &mock.MockUserRepository{
				GetSettingsFunc: func(int) (*models.SettingsGet, error) {
					return &models.SettingsGet{
						Email: "LaHaine@gmail.com",
						Nickname: "LaHaine",
						SmallAvatar: "small",
					}, nil
				},
			},
			input: 1,
			expected: response{
				settings: &models.SettingsGet{
					Email: "LaHaine@gmail.com",
					Nickname: "LaHaine",
					SmallAvatar: "small",
				},
				err: nil,
			},
		},
		{
			name: "GetSettings returns an error",
			dbMock: &mock.MockUserRepository{
				GetSettingsFunc: func(int) (*models.SettingsGet, error) {
					return nil, errors.New("error")
				},
			},
			input: 1,
			expected: response{
				settings: nil,
				err: &models.CustomError{
					ErrorType: 500,
					OriginalError: errors.New("error"),
				},
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.sessionCheckerMock, testCase.imagesMock)
			got := response{}

			got.settings, got.err = r.GetSettings(testCase.input)

			if testCase.expectedErr {
				assert.NotNil(t, got.err)
				assert.Equal(t, testCase.expected, got)
			} else {
				assert.Nil(t, got.err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserUseCase_GetAvatarFilename(t *testing.T) {
	type response struct {
		filename string
		err 	 *models.CustomError
	}

	tests := []struct {
		name 		  		 string
		dbMock 	  	  		 *mock.MockUserRepository
		sessionCheckerMock   *mock.MockSessionCheckerClient
		imagesMock			 *mock.MockAvatarRepositoryIFace
		input 		  		 int
		expected 	  		 response
		expectedErr   		 bool
	}{
		{
			name: "Successfully got avatar filename",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "avatar", nil
				},
			},
			input: 1,
			expected: response{
				filename: "avatar" + constants.LittleAvatarPostfix,
			},
		},
		{
			name: "GetAvatarFilename returned error",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "", errors.New("error")
				},
			},
			input: 1,
			expected: response{
				err: &models.CustomError{
					ErrorType: http.StatusInternalServerError,
					OriginalError: errors.New("error"),
				},
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.sessionCheckerMock, testCase.imagesMock)
			got := response{}

			got.filename, got.err = r.GetAvatarFilename(testCase.input)

			if testCase.expectedErr {
				assert.NotNil(t, got.err)
				assert.Equal(t, testCase.expected.err.ErrorType, got.err.ErrorType)
			} else {
				assert.Nil(t, got.err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserUseCase_UpdateSettings(t *testing.T) {
	type inputStruct struct {
		userId int
		oldSettings *models.SettingsGet
		newSettings *models.SettingsUpload
	}

	tests := []struct {
		name 	  	  		 string
		dbMock 	  	  		 *mock.MockUserRepository
		sessionCheckerMock   *mock.MockSessionCheckerClient
		imagesMock			 *mock.MockAvatarRepositoryIFace
		input 	  	  		 *inputStruct
		expected  	  		 *models.CustomError
		expectedErr   		 bool
	}{
		// ----------EMAIL----------
		{
			name: "Successfully update email",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateEmailFunc: func(int, string) error {
					return nil
				},
			},

			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "LaHaine@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "LaHaineI@gmail.com",
				},
			},
		},
		{
			name: "Unsuccessfully update email, email is not valid",
			dbMock: &mock.MockUserRepository{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "LaHaine@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "LaHaine",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
				Message: constants.InvalidEmailMessage,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update email, email is not unique",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "LaHaine@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "LaHaineI@gmail.com",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
				Message: constants.NotUniqueEmailMessage,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update email, IsEmailUnique returns error",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "LaHaine@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "LaHaineI@gmail.com",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update email, UpdateEmail returns error",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateEmailFunc: func(int, string) error {
					return errors.New("error")
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "LaHaine@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "LaHaineI@gmail.com",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},

		// ----------NICKNAME----------
		{
			name: "Successfully update nickname",
			dbMock: &mock.MockUserRepository{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateNicknameFunc: func(int, string) error {
					return nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "LaHaine",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "LaHaineI",
				},
			},
		},
		{
			name: "Unsuccessfully update nickname, nickname is not valid",
			dbMock: &mock.MockUserRepository{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "LaHaine",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "La",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
				Message: constants.InvalidNicknameMessage,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update nickname, nickname is not unique",
			dbMock: &mock.MockUserRepository{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "LaHaine",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "LaHaineI",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
				Message: constants.NotUniqueNicknameMessage,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update nickname, IsNickNameUnique returns error",
			dbMock: &mock.MockUserRepository{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "LaHaine",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "LaHaineI",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update nickname, UpdateNickname returns error",
			dbMock: &mock.MockUserRepository{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateNicknameFunc: func(int, string) error {
					return errors.New("error")
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "LaHaine",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "LaHaineI",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},

		// ----------PASSWORD----------
		{
			name: "Successfully update password",
			dbMock: &mock.MockUserRepository{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return true, nil
				},
				UpdatePasswordFunc: func(int, string) error {
					return nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "JesusLovesMe",
					NewPassword: "JesusLovesMe1337!",
				},
			},
		},
		{
			name: "Unsuccessfully update password, CheckPasswordByUserID returns false",
			dbMock: &mock.MockUserRepository{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return false, nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "JesusLovesMe",
					NewPassword: "JesusLovesMe1337!",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
				Message: constants.WrongPasswordMessage,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, CheckPasswordByUserID returns error",
			dbMock: &mock.MockUserRepository{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "JesusLovesMe",
					NewPassword: "JesusLovesMe1337!",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, new password is invalid",
			dbMock: &mock.MockUserRepository{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return true, nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "JesusLovesMe",
					NewPassword: "JesusLovesMe1337",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
				Message: constants.PasswordValidationNoSpecialSymbolMessage,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, UpdatePassword returns error",
			dbMock: &mock.MockUserRepository{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return true, nil
				},
				UpdatePasswordFunc: func(int, string) error {
					return errors.New("error")
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "JesusLovesMe",
					NewPassword: "JesusLovesMe1337!",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, old password field is empty",
			dbMock: &mock.MockUserRepository{ },
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "",
					NewPassword: "JesusLovesMe1337!",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
				Message: constants.OldPasswordFieldIsEmptyMessage,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, new password field is empty",
			dbMock: &mock.MockUserRepository{ },
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "JesusLovesMe",
					NewPassword: "",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
				Message: constants.NewPasswordFieldIsEmptyMessage,
			},
			expectedErr: true,
		},
		{
			name: "Successfully updated avatar",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "avatarimg", nil
				},
				UpdateAvatarFunc: func(int, string) error {
					return nil
				},
			},
			imagesMock: &mock.MockAvatarRepositoryIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "avatar", nil
				},
				DeleteImageFunc: func(string) error {
					return nil
				},
			},
			input: &inputStruct {
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_avatar",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "avatar",
				},
			},
		},
		{
			name: "CreateImage returned error",
			dbMock: &mock.MockUserRepository{},
			imagesMock: &mock.MockAvatarRepositoryIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "", errors.New("error")
				},
			},
			input: &inputStruct {
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_avatar",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "avatar",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},
		{
			name: "db.GetAvatarFilename returned error",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "", errors.New("error")
				},
			},
			imagesMock: &mock.MockAvatarRepositoryIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "avatar", nil
				},
			},
			input: &inputStruct {
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_avatar",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "avatar",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},
		{
			name: "DeleteImage returned error",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "avatar", nil
				},
			},
			imagesMock: &mock.MockAvatarRepositoryIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "avatar", nil
				},
				DeleteImageFunc: func(string) error {
					return errors.New("error")
				},
			},
			input: &inputStruct {
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_avatar",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "avatar",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},
		{
			name: "db.UpdateAvatar returned error",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "avatar", nil
				},
				UpdateAvatarFunc: func(int,  string) error {
					return errors.New("error")
				},
			},
			imagesMock: &mock.MockAvatarRepositoryIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "avatar", nil
				},
				DeleteImageFunc: func(string) error {
					return nil
				},
			},
			input: &inputStruct {
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_avatar",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "avatar",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
				OriginalError: errors.New("error"),
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.sessionCheckerMock, testCase.imagesMock)

			customError := r.UpdateSettings(testCase.input.userId, testCase.input.oldSettings, testCase.input.newSettings)

			if testCase.expectedErr {
				assert.NotNil(t, customError)
				assert.Equal(t, testCase.expected, customError)
			} else {
				assert.Nil(t, customError)
			}
		})
	}
}


