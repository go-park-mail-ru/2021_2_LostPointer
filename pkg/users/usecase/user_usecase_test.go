package usecase

import (
	"2021_2_LostPointer/pkg/mock"
	"2021_2_LostPointer/pkg/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

func TestUserUseCase_Login(t *testing.T) {
	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepositoryIFace
		redisMock 	  *mock.MockRedisStoreIFace
		fsMock		  *mock.MockFileSystemIFace
		input 	  	  models.Auth
		expected  	  string
		expectedErr   bool
	}{
		{
			name: "Successful login",
			dbMock: &mock.MockUserRepositoryIFace{
				DoesUserExistFunc: func(models.Auth) (uint64, error) {
					return 1, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "some_token", nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: models.Auth{
				Email: "alex1234@gmail.com",
				Password: "alex1234",
			},
			expected: "some_token",
		},
		{
			name: "Invalid credentials",
			dbMock: &mock.MockUserRepositoryIFace{
				DoesUserExistFunc: func(models.Auth) (uint64, error) {
					return 0, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "some_token", nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: models.Auth{
				Email: "alex1234@gmail.com",
				Password: "alex1234",
			},
			expected: "",
			expectedErr: true,
		},
		{
			name: "StoreSession error",
			dbMock: &mock.MockUserRepositoryIFace{
				DoesUserExistFunc: func(models.Auth) (uint64, error) {
					return 1, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "", errors.New("redis_error")
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: models.Auth{
				Email: "alex1234@gmail.com",
				Password: "alex1234",
			},
			expected: "",
			expectedErr: true,
		},
		{
			name: "DoesUserExist error",
			dbMock: &mock.MockUserRepositoryIFace{
				DoesUserExistFunc: func(models.Auth) (uint64, error) {
					return 0, errors.New("sql_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{ },
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
	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepositoryIFace
		redisMock 	  *mock.MockRedisStoreIFace
		fsMock 		  *mock.MockFileSystemIFace
		input 	  	  string
		expected  	  bool
		expectedErr   bool
	}{
		{
			name: "User is authorized",
			dbMock: &mock.MockUserRepositoryIFace{ },
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: "some_cookie",
			expected: true,
		},
		{
			name: "User is not authorized",
			dbMock: &mock.MockUserRepositoryIFace{ },
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 0, errors.New("redis_error")
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: "some_cookie",
			expected: false,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock, testCase.fsMock)

			got, customError := r.IsAuthorized(testCase.input)

			if testCase.expectedErr {
				assert.NotNil(t, customError)
			} else {
				assert.Nil(t, customError)
				assert.Equal(t, testCase.expected, got)
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
	dbMock := &mock.MockUserRepositoryIFace{ }
	redisMock := &mock.MockRedisStoreIFace{
		DeleteSessionFunc: func(string) {},
	}
	fsMock := &mock.MockFileSystemIFace{}

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
		dbMock 	  	  *mock.MockUserRepositoryIFace
		redisMock 	  *mock.MockRedisStoreIFace
		fsMock        *mock.MockFileSystemIFace
		input 	  	  models.User
		expected  	  response
		expectedErr   bool
	}{
		{
			name: "Successfully",
			dbMock: &mock.MockUserRepositoryIFace{
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
			redisMock: &mock.MockRedisStoreIFace{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "some_token", nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
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
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				CreateUserFunc: func(models.User, ...string) (uint64, error) {
					return 0, errors.New("some_error_in_CreateUserFunct")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "some_token", nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "Invalid credentials",
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{ },
			fsMock: &mock.MockFileSystemIFace{},
			input: models.User{
				Email: "alexeikosenkogmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "Email is not unique",
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{ },
			fsMock: &mock.MockFileSystemIFace{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "IsEmailUniqueFunc returns error",
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, errors.New("some_error_in_IsEmailUniqueFunc")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{ },
			fsMock: &mock.MockFileSystemIFace{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "Nickname is not unique",
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{ },
			fsMock: &mock.MockFileSystemIFace{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "IsNicknameUnique returns error",
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, errors.New("some_error_in_IsNicknameUniqueFunc")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{ },
			fsMock: &mock.MockFileSystemIFace{},
			input: models.User{
				Email: "alexeikosenko@gmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expectedErr: true,
		},
		{
			name: "StoreSession returns error",
			dbMock: &mock.MockUserRepositoryIFace{
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
			redisMock: &mock.MockRedisStoreIFace{
				StoreSessionFunc: func(uint64, ...string) (string, error) {
					return "", errors.New("some_error_in_StoreSessionFunc")
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
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

func TestUserUseCase_GetSettings(t *testing.T) {
	type response struct {
		settings *models.SettingsGet
		err 	 *models.CustomError
	}

	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepositoryIFace
		redisMock 	  *mock.MockRedisStoreIFace
		fsMock        *mock.MockFileSystemIFace
		input 	  	  string
		expected  	  response
		expectedErr   bool
	}{
		{
			name: "Successfully got settings",
			dbMock: &mock.MockUserRepositoryIFace{
				GetSettingsFunc: func(int) (*models.SettingsGet, error) {
					return &models.SettingsGet{
						Email: "alex1234@gmail.com",
						Nickname: "alex1234",
						SmallAvatar: "SmallAvatar",
					}, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: "some_cookie_value",
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
			name: "GetSessionUserId returns an error",
			dbMock: &mock.MockUserRepositoryIFace{},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 0, errors.New("some_error_in_redis")
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: "some_cookie_value",
			expected: response{
				settings: nil,
				err: &models.CustomError{
					ErrorType: 401,
				},
			},
			expectedErr: true,
		},
		{
			name: "GetSettings returns an error",
			dbMock: &mock.MockUserRepositoryIFace{
				GetSettingsFunc: func(int) (*models.SettingsGet, error) {
					return nil, errors.New("some_error_in_GetSettings")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: "some_cookie_value",
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
		cookieValue string
		oldSettings *models.SettingsGet
		newSettings *models.SettingsUpload
	}

	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepositoryIFace
		redisMock 	  *mock.MockRedisStoreIFace
		fsMock        *mock.MockFileSystemIFace
		input 	  	  *inputStruct
		expected  	  *models.CustomError
		expectedErr   bool
	}{
		// ----------EMAIL----------
		{
			name: "Successfully update email",
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateEmailFunc: func(int, string) error {
					return nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return false, errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				IsEmailUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateEmailFunc: func(int, string) error {
					return errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateNicknameFunc: func(int, string) error {
					return nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return false, errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				IsNicknameUniqueFunc: func(string) (bool, error) {
					return true, nil
				},
				UpdateNicknameFunc: func(int, string) error {
					return errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return true, nil
				},
				UpdatePasswordFunc: func(int, string, ...string) error {
					return nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
				oldSettings: &models.SettingsGet{},
				newSettings: &models.SettingsUpload{
					OldPassword: "alex1234",
					NewPassword: "Alex1234!",
				},
			},
		},
		{
			name: "Unsuccessfully update password, CheckPasswordByUserID returns false",
			dbMock: &mock.MockUserRepositoryIFace{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return false, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return false, errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return true, nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				CheckPasswordByUserIDFunc: func(int, string) (bool, error) {
					return true, nil
				},
				UpdatePasswordFunc: func(int, string, ...string) error {
					return errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{ },
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{ },
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "some_filename", nil
				},
				UpdateAvatarFunc: func(int, string) error {
					return nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "some_filename_new", nil
				},
				DeleteImageFunc: func(string) error {
					return nil
				},
			},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "", errors.New("some_error")
				},
			},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "", errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "some_filename_new", nil
				},
			},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "some_filename", nil
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "some_filename_new", nil
				},
				DeleteImageFunc: func(string) error {
					return errors.New("some_error")
				},
			},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			dbMock: &mock.MockUserRepositoryIFace{
				GetAvatarFilenameFunc: func(int) (string, error) {
					return "some_filename", nil
				},
				UpdateAvatarFunc: func(int, string) error {
					return errors.New("some_error")
				},
			},
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			fsMock: &mock.MockFileSystemIFace{
				CreateImageFunc: func(*multipart.FileHeader) (string, error) {
					return "some_filename_new", nil
				},
				DeleteImageFunc: func(string) error {
					return nil
				},
			},
			input: &inputStruct{
				cookieValue: "some_cookie",
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
			name: "Unsuccessfully update settings, redis returns error",
			dbMock: &mock.MockUserRepositoryIFace{ },
			redisMock: &mock.MockRedisStoreIFace{
				GetSessionUserIdFunc: func(string) (int, error) {
					return 0, errors.New("some_error")
				},
			},
			fsMock: &mock.MockFileSystemIFace{ },
			input: &inputStruct{
				cookieValue: "some_cookie",
				oldSettings: &models.SettingsGet{
					SmallAvatar: "old_filename",
				},
				newSettings: &models.SettingsUpload{
					AvatarFileName: "new_filename",
				},
			},
			expected: &models.CustomError{
				ErrorType: 401,
			},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock, testCase.fsMock)

			customError := r.UpdateSettings(testCase.input.cookieValue, testCase.input.oldSettings, testCase.input.newSettings)

			if testCase.expectedErr {
				assert.NotNil(t, customError)
				assert.Equal(t, testCase.expected.ErrorType, customError.ErrorType)
			} else {
				assert.Nil(t, customError)
			}
		})
	}
}
