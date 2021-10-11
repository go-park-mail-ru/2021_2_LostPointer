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
		name 	  string
		dbMock 	  *mock.MockUserRepositoryIFace
		redisMock *mock.MockRedisStoreIFace
		input 	  models.Auth
		expected  string
		expectedErr   bool
	}{
		{
			name: "Successful login",
			dbMock: &mock.MockUserRepositoryIFace{
				UserExitsFunc: func(models.Auth) (uint64, error) {
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
				UserExitsFunc: func(models.Auth) (uint64, error) {
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
		},
		{
			name: "StoreSession error",
			dbMock: &mock.MockUserRepositoryIFace{
				UserExitsFunc: func(models.Auth) (uint64, error) {
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
			name: "UserExist error",
			dbMock: &mock.MockUserRepositoryIFace{
				UserExitsFunc: func(models.Auth) (uint64, error) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			r := NewUserUserCase(tt.dbMock, tt.redisMock)

			got, err := r.Login(tt.input)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestUserUseCase_IsAuthorized(t *testing.T) {
	tests := []struct {
		name 	  string
		dbMock 	  *mock.MockUserRepositoryIFace
		redisMock *mock.MockRedisStoreIFace
		input 	  string
		expected  bool
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
					return 0, nil
				},
			},
			input: "some_cookie",
			expected: false,
		},
		{
			name: "Redis error",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewUserUserCase(tt.dbMock, tt.redisMock)

			got, err := r.IsAuthorized(tt.input)

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
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
		name 	 string
		password string
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
				message: "Password must contain at least 8 characters",
				err: nil,
			},
		},
		{
			name: "No lowercase letter",
			password: "AVTFAUDSADIAODIS8!",
			expected: response{
				isValid: false,
				message: "Password must contain at least one lowercase letter",
				err: nil,
			},
		},
		{
			name: "No uppercase letter",
			password: "avtdsaopdsodpasdos8!",
			expected: response{
				isValid: false,
				message: "Password must contain at least one uppercase letter",
				err: nil,
			},
		},
		{
			name: "No digit",
			password: "Avtdskdksdlskdladkl!",
			expected: response{
				isValid: false,
				message: "Password must contain at least one digit",
				err: nil,
			},
		},
		{
			name: "No special symbol",
			password: "Avtdskdksdlskd8",
			expected: response{
				isValid: false,
				message: "Password must contain as least one special character",
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := response{}

			got.isValid, got.message, got.err = ValidatePassword(tt.password)

			assert.Equal(t, tt.expected, got)
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
		name string
		input models.User
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
				message: "The length of nickname must be from 3 to 15 characters",
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
				message: "Invalid email",
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
				message: "Password must contain at least 8 characters",
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := response{}

			got.isValid, got.message, got.err = ValidateRegisterCredentials(tt.input)

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestUserUseCase_Register(t *testing.T) {
	type response struct {
		sessionToken string
		message 	 string
		err 		 error
	}

	tests := []struct {
		name 	  string
		dbMock 	  *mock.MockUserRepositoryIFace
		redisMock *mock.MockRedisStoreIFace
		input 	  models.User
		expected  response
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
				message: "",
				err: nil,
			},
		},
		{
			name: "Invalid credentials",
			dbMock: &mock.MockUserRepositoryIFace{ },
			redisMock: &mock.MockRedisStoreIFace{ },
			input: models.User{
				Email: "alexeikosenkogmail.com",
				Password: "AlexeiKosenkaRulitTankom1234!",
				Nickname: "alexeiKosenka",
			},
			expected: response{
				sessionToken: "",
				message: "Invalid email",
				err: nil,
			},
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
			expected: response{
				sessionToken: "",
				message: "Email is already taken",
				err: nil,
			},
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
			expected: response{
				sessionToken: "",
				message: "",
				err: errors.New("some_error_in_IsEmailUniqueFunc"),
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
			expected: response{
				sessionToken: "",
				message: "Nickname is already taken",
				err: nil,
			},
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
			expected: response{
				sessionToken: "",
				message: "",
				err: errors.New("some_error_in_IsNicknameUniqueFunc"),
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
			expected: response{
				sessionToken: "",
				message: "",
				err: errors.New("some_error_in_StoreSessionFunc"),
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewUserUserCase(tt.dbMock, tt.redisMock)
			got := response{}

			got.sessionToken, got.message, got.err = r.Register(tt.input)

			if tt.expectedErr {
				assert.Error(t, got.err)
			} else {
				assert.NoError(t, got.err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}
