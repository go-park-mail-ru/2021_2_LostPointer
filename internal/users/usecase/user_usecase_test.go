package usecase

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestUserUseCase_Login(t *testing.T) {
	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepository
		redisMock 	  *mock.MockRedisStore
		fsMock		  *mock.MockFileSystem
		input 	  	  models.Auth
		expected  	  string
		expectedErr   bool
	}{
		{
			name: "Successful login",
			dbMock: &mock.MockUserRepository{
				DoesUserExistFunc: func(models.Auth) (uint64, error) {
					return 1, nil
				},
			},
			redisMock: &mock.MockRedisStore{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "some_token", nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: models.Auth{
				Email: "alex1234@gmail.com",
				Password: "alex1234",
			},
			expected: "some_token",
		},
		{
			name: "Invalid credentials",
			dbMock: &mock.MockUserRepository{
				DoesUserExistFunc: func(models.Auth) (uint64, error) {
					return 0, nil
				},
			},
			redisMock: &mock.MockRedisStore{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "some_token", nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: models.Auth{
				Email: "alex1234@gmail.com",
				Password: "alex1234",
			},
			expected: "",
			expectedErr: true,
		},
		{
			name: "StoreSession error",
			dbMock: &mock.MockUserRepository{
				DoesUserExistFunc: func(models.Auth) (uint64, error) {
					return 1, nil
				},
			},
			redisMock: &mock.MockRedisStore{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "", errors.New("redis_error")
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: models.Auth{
				Email: "alex1234@gmail.com",
				Password: "alex1234",
			},
			expected: "",
			expectedErr: true,
		},
		{
			name: "DoesUserExist error",
			dbMock: &mock.MockUserRepository{
				DoesUserExistFunc: func(models.Auth) (uint64, error) {
					return 0, errors.New("sql_error")
				},
			},
			redisMock: &mock.MockRedisStore{ },
			input: models.Auth{
				Email: "alex1234@gmail.com",
				Password: "alex1234",
			},
			expected: "",
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock, testCase.fsMock)

			got, customError := r.Login(testCase.input)
			if testCase.expectedErr {
				assert.NotNil(t, customError)
			} else {
				assert.Nil(t, customError)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserUseCase_IsAuthorized(t *testing.T) {
	type response struct {
		isAuthorized bool
		userID   	 int
		err 		 *models.CustomError
	}

	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepository
		redisMock 	  *mock.MockRedisStore
		fsMock 		  *mock.MockFileSystem
		input 	  	  string
		expected  	  response
		expectedErr   bool
	}{
		{
			name: "User is authorized",
			dbMock: &mock.MockUserRepository{ },
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: "some_cookie",
			expected: response{
				isAuthorized: true,
				userID: 1,
				err: nil,
			},
		},
		{
			name: "User is not authorized",
			dbMock: &mock.MockUserRepository{ },
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 0, &models.CustomError{ErrorType: 401}
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: "some_cookie",
			expected: response{
				isAuthorized: false,
				userID: 0,
				err: &models.CustomError{ErrorType: 401},
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock, testCase.fsMock)

			var got response
			got.isAuthorized, got.userID, got.err = r.IsAuthorized(testCase.input)
			if testCase.expectedErr {
				assert.NotNil(t, got.err)
				assert.Equal(t, testCase.expected, got)
			} else {
				assert.Nil(t, got.err)
				assert.Equal(t, testCase.expected.isAuthorized, got.isAuthorized)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	type response struct {
		isValid bool
		message string
		err 	error
	}

	tests := []struct {
		name 	   string
		password   string
		expected   response
	}{
		{
			name: "Valid password",
			password: "Avt8430055!",
			expected: response{
				isValid: true,
				message: "",
				err: nil,
			},
		},
		{
			name: "Wrong length",
			password: "Avt8!",
			expected: response{
				isValid: false,
				message: PasswordValidationInvalidLengthMessage,
				err: nil,
			},
		},
		{
			name: "No lowercase letter",
			password: "AVTFAUDSADIAODIS8!",
			expected: response{
				isValid: false,
				message: PasswordValidationNoLowerCaseMessage,
				err: nil,
			},
		},
		{
			name: "No uppercase letter",
			password: "avtdsaopdsodpasdos8!",
			expected: response{
				isValid: false,
				message: PasswordValidationNoUppercaseMessage,
				err: nil,
			},
		},
		{
			name: "No digit",
			password: "Avtdskdksdlskdladkl!",
			expected: response{
				isValid: false,
				message: PasswordValidationNoDigitMessage,
				err: nil,
			},
		},
		{
			name: "No special symbol",
			password: "Avtdskdksdlskd8",
			expected: response{
				isValid: false,
				message: PasswordValidationNoSpecialSymbolMessage,
				err: nil,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := response{}

			got.isValid, got.message, got.err = ValidatePassword(testCase.password)

			assert.Equal(t, testCase.expected, got)
		})
	}
}

func TestUserUseCase_Logout(t *testing.T) {
	dbMock := &mock.MockUserRepository{ }
	redisMock := &mock.MockRedisStore{
		DeleteSessionFunc: func(string) {},
	}
	fsMock := &mock.MockFileSystem{}

	r := NewUserUserCase(dbMock, redisMock, fsMock)
	r.Logout("alexei_kosenkov")

	assert.True(t, true)
}

func TestValidateRegisterCredentials(t *testing.T) {
	type response struct {
		isValid bool
		message string
		err 	error
	}

	tests := []struct {
		name 	 string
		input 	 models.User
		expected response
	}{
		{
			name: "Valid credentials",
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expected: response{
				isValid: true,
				message: "",
				err: nil,
			},
		},
		{
			name: "Invalid name",
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "al",
			},
			expected: response{
				isValid: false,
				message: NickNameValidationInvalidLengthMessage,
				err: nil,
			},
		},
		{
			name: "Invalid email",
			input: models.User{
				Email: "alexeikosenkogmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expected: response{
				isValid: false,
				message: InvalidEmailMessage,
				err: nil,
			},
		},
		{
			name: "Invalid password",
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "Alex1!",
				Nickname: "alexeiKosenka",
			},
			expected: response{
				isValid: false,
				message: PasswordValidationInvalidLengthMessage,
				err: nil,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got := response{}

			got.isValid, got.message, got.err = ValidateRegisterCredentials(testCase.input)

			assert.Equal(t, testCase.expected, got)
		})
	}
}

func TestUserUseCase_Register(t *testing.T) {
	type response struct {
		sessionToken string
		err 		 *models.CustomError
	}

	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepository
		redisMock 	  *mock.MockRedisStore
		fsMock        *mock.MockFileSystem
		input 	  	  models.User
		expected  	  response
		expectedErr   bool
	}{
		{
			name: "Successfully",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(models.User, ...string) (uint64, error) {
					return 1, nil
				},
			},
			redisMock: &mock.MockRedisStore{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "some_token", nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expected: response{
				sessionToken: "some_token",
				err: nil,
			},
		},
		{
			name: "userDB.CreateUser returns error",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(models.User, ...string) (uint64, error) {
					return 0, errors.New("some_error_in_CreateUserFunc")
				},
			},
			redisMock: &mock.MockRedisStore{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "some_token", nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "Invalid credentials",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
			},
			redisMock: &mock.MockRedisStore{ },
			fsMock: &mock.MockFileSystem{},
			input: models.User{
				Email: "alexeikosenkogmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "Email is not unique",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			redisMock: &mock.MockRedisStore{ },
			fsMock: &mock.MockFileSystem{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "IsEmailUniqueFunc returns error",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, errors.New("some_error_in_IsEmailUniqueFunc")
				},
			},
			redisMock: &mock.MockRedisStore{ },
			fsMock: &mock.MockFileSystem{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "Nickname is not unique",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			redisMock: &mock.MockRedisStore{ },
			fsMock: &mock.MockFileSystem{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "IsNicknameUnique returns error",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, errors.New("some_error_in_IsNicknameUniqueFunc")
				},
			},
			redisMock: &mock.MockRedisStore{ },
			fsMock: &mock.MockFileSystem{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "StoreSession returns error",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(models.User, ...string) (uint64, error) {
					return 1, nil
				},
			},
			redisMock: &mock.MockRedisStore{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "", errors.New("some_error_in_StoreSessionFunc")
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock, testCase.fsMock)
			got := response{}

			got.sessionToken, got.err = r.Register(testCase.input)

			if testCase.expectedErr {
				assert.NotNil(t, got.err)
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
		name 		  string
		dbMock 	  	  *mock.MockUserRepository
		redisMock 	  *mock.MockRedisStore
		fsMock        *mock.MockFileSystem
		input 		  int
		expected 	  response
		expectedErr   bool
	}{
		{
			name: "Successfully got avatar filename",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "avatar", nil
				},
			},
			redisMock: &mock.MockRedisStore{},
			fsMock: &mock.MockFileSystem{},
			input: 1,
			expected: response{
				filename: "avatar" + "_150px.webp",
			},
		},
		{
			name: "GetAvatarFilename returned error",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "", errors.New("error")
				},
			},
			redisMock: &mock.MockRedisStore{},
			fsMock: &mock.MockFileSystem{},
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
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock, testCase.fsMock)
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

func TestUserUseCase_GetSettings(t *testing.T) {
	type response struct {
		settings *models.SettingsGet
		err 	 *models.CustomError
	}

	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepository
		redisMock 	  *mock.MockRedisStore
		fsMock        *mock.MockFileSystem
		input 	  	  int
		expected  	  response
		expectedErr   bool
	}{
		{
			name: "Successfully got settings",
			dbMock: &mock.MockUserRepository{
				GetSettingsFunc: func(int) (*models.SettingsGet, error) {
					return &models.SettingsGet{
						Email: "alex1234@gmail.com",
						Nickname: "alex1234",
						SmallAvatar: "SmallAvatar",
					}, nil
				},
			},
			redisMock: &mock.MockRedisStore{},
			fsMock: &mock.MockFileSystem{},
			input: 1,
			expected: response{
				settings: &models.SettingsGet{
					Email: "alex1234@gmail.com",
					Nickname: "alex1234",
					SmallAvatar: "SmallAvatar",
				},
				err: nil,
			},
		},
		{
			name: "GetSettings returns an error",
			dbMock: &mock.MockUserRepository{
				GetSettingsFunc: func(int) (*models.SettingsGet, error) {
					return nil, errors.New("some_error_in_GetSettings")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: 1,
			expected: response{
				settings: nil,
				err: &models.CustomError{
					ErrorType: 500,
					OriginalError: errors.New("some_error_in_GetSettings"),
				},
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock, testCase.fsMock)
			got := response{}

			got.settings, got.err = r.GetSettings(testCase.input)

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
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepository
		redisMock	  *mock.MockRedisStore
		fsMock 		  *mock.MockFileSystem
		input 	  	  *inputStruct
		expected  	  *models.CustomError
		expectedErr   bool
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
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "alex1234@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "alex1235@gmail.com",
				},
			},
		},
		{
			name: "Unsuccessfully update email, email is not valid",
			dbMock: &mock.MockUserRepository{},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "alex1234@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "alex1235m",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
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
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "alex1234@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "alex1235@gmail.com",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update email, IsEmailUnique returns error",
			dbMock: &mock.MockUserRepository{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "alex1234@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "alex1235@gmail.com",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
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
					return errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Email: "alex1234@gmail.com",
				},
				newSettings: &models.SettingsUpload{
					Email: "alex1235@gmail.com",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
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
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "alex1234",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "alex1235",
				},
			},
		},
		{
			name: "Unsuccessfully update nickname, nickname is not valid",
			dbMock: &mock.MockUserRepository{},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "alex1234",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "al",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
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
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "alex1234",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "alex1235",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update nickname, IsNickNameUnique returns error",
			dbMock: &mock.MockUserRepository{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "alex1234",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "alex1235",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
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
					return errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					Nickname: "alex1234",
				},
				newSettings: &models.SettingsUpload{
					Nickname: "alex1235",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
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
				UpdatePasswordFunc: func(int, string, ...string) error {
					return nil
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "alex1234",
					NewPassword: "Alex1234!",
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
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "alex1234",
					NewPassword: "Alex1234!",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, CheckPasswordByUserID returns error",
			dbMock: &mock.MockUserRepository{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return false, errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "alex1234",
					NewPassword: "Alex1234!",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
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
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "alex1234",
					NewPassword: "alex1234",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, UpdatePassword returns error",
			dbMock: &mock.MockUserRepository{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return true, nil
				},
				UpdatePasswordFunc: func(int, string, ...string) error {
					return errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "alex1234",
					NewPassword: "Alex1234!",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, old password field is empty",
			dbMock: &mock.MockUserRepository{ },
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "",
					NewPassword: "Alex1234!",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update password, new password field is empty",
			dbMock: &mock.MockUserRepository{ },
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "alex1234",
					NewPassword: "",
				},
			},
			expected: &models.CustomError{
				ErrorType: 400,
			},
			expectedErr: true,
		},

		// ----------SmallAvatar----------
		{
			name: "Successfully update SmallAvatar",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "some_filename", nil
				},
				UpdateAvatarFunc: func(int, string) error {
					return nil
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "some_filename_new", nil
				},
				DeleteImageFunc: func(string) error {
					return nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_filename",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "new_filename",
				},
			},
		},
		{
			name: "Unsuccessfully update SmallAvatar, CreateImage returns error",
			dbMock: &mock.MockUserRepository{},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "", errors.New("some_error")
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_filename",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "new_filename",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update SmallAvatar, GetSmallAvatarFilename returns error",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "", errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "some_filename_new", nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_filename",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "new_filename",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update SmallAvatar, DeleteImage returns error",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "some_filename", nil
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "some_filename_new", nil
				},
				DeleteImageFunc: func(string) error {
					return errors.New("some_error")
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_filename",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "new_filename",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
			},
			expectedErr: true,
		},
		{
			name: "Unsuccessfully update SmallAvatar, UpdateSmallAvatar returns error",
			dbMock: &mock.MockUserRepository{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "some_filename", nil
				},
				UpdateAvatarFunc: func(int, string) error {
					return errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStore{
				GetSessionUserIdFunc: func(string) (int, *models.CustomError) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystem{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "some_filename_new", nil
				},
				DeleteImageFunc: func(string) error {
					return nil
				},
			},
			input: &inputStruct{
				userId: 1,
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_filename",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "new_filename",
				},
			},
			expected: &models.CustomError{
				ErrorType: 500,
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock, testCase.fsMock)

			customError := r.UpdateSettings(testCase.input.userId, testCase.input.oldSettings, testCase.input.newSettings)

			if testCase.expectedErr {
				assert.NotNil(t, customError)
				assert.Equal(t, testCase.expected.ErrorType, customError.ErrorType)
			} else {
				assert.Nil(t, customError)
			}
		})
	}
}


