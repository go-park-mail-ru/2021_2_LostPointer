package usecase

import (
	"2021_2_LostPointer/pkg/mock"
	"2021_2_LostPointer/pkg/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserUseCase_Login(t *testing.T) {
	tests := []struct {
		name 	  	  string
		dbMock 	  	  *mock.MockUserRepositoryIFace
		redisMock 	  *mock.MockRedisStoreIFace
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
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock)

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
			input: "some_cookie",
			expected: false,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock)

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

	r := NewUserUserCase(dbMock, redisMock)
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
			r := NewUserUserCase(testCase.dbMock, testCase.redisMock)
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
