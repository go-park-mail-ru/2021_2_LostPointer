package repository

import (
	"database/sql/driver"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/authorization/proto"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	postgresDB, _, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	redisDB, mock := redismock.NewClientMock()
	repository := NewAuthStorage(postgresDB, redisDB)

	var id int64 = 1

	type inputStruct struct {
		id          int64
		cookieValue string
	}
	tests := []struct {
		name          string
		mock          func()
		input         inputStruct
		expected      error
		expectedError bool
	}{
		{
			name: "Successfully stored session",
			mock: func() {
				mock.ExpectSet("cookie_value", id, constants.CookieLifetime).SetVal("")
			},
			input: inputStruct{
				id:          1,
				cookieValue: "cookie_value",
			},
			expectedError: false,
		},
		{
			name: "Error occurred in redis.Set",
			mock: func() {
				mock.ExpectSet("cookie_value", id, constants.CookieLifetime).SetErr(errors.New("error"))
			},
			input: inputStruct{
				id:          1,
				cookieValue: "cookie_value",
			},
			expectedError: true,
		},
	}
	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.CreateSession(currentTest.input.id, currentTest.input.cookieValue)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionRepository_GetUserIdByCookie(t *testing.T) {
	postgresDB, _, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	redisDB, mock := redismock.NewClientMock()
	repository := NewAuthStorage(postgresDB, redisDB)

	tests := []struct {
		name        string
		mock        func()
		input       string
		expected    int64
		expectedErr bool
	}{
		{
			name: "Successfully returned user id by cookie value",
			mock: func() {
				mock.ExpectGet("cookie_value").SetVal("1")
			},
			input:    "cookie_value",
			expected: 1,
		},
		{
			name: "Error occurred in redis.Get",
			mock: func() {
				mock.ExpectGet("cookie_value").SetErr(errors.New("error"))
			},
			input:       "cookie_value",
			expectedErr: true,
		},
	}
	for _, testCase := range tests {
		currentTest := testCase
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			res, err := repository.GetUserByCookie(currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, currentTest.expected, res)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionRepository_DeleteSession(t *testing.T) {
	postgresDB, _, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	redisDB, mock := redismock.NewClientMock()
	repository := NewAuthStorage(postgresDB, redisDB)

	tests := []struct {
		name        string
		mock        func()
		expected    error
		input       string
		expectedErr bool
	}{
		{
			name: "Successfully deleted session",
			mock: func() {
				mock.ExpectDel("cookie_value").SetVal(1)
			},
			input: "cookie_value",
		},
		{
			name: "Error occurred in redis.Del",
			mock: func() {
				mock.ExpectDel("cookie_value").SetErr(errors.New("error"))
			},
			input:       "cookie_value",
			expectedErr: true,
		},
	}
	for _, testCase := range tests {
		currentTest := testCase
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.DeleteSession(currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthStorage_GetUserByPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	redisDB, _ := redismock.NewClientMock()
	repository := NewAuthStorage(db, redisDB)

	authData := &proto.AuthData{
		Email:    "testEmail",
		Password: "testPassword",
	}

	const (
		ID   = 1
		salt = "testSalt"
	)
	dbPassword, _ := bcrypt.GenerateFromPassword([]byte("testPassword"+salt), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		mock          func()
		expected      int64
		expectedError bool
	}{
		{
			name: "get user success",
			mock: func() {
				row := mock.NewRows([]string{"email", "password", "salt"})
				row.AddRow(ID, dbPassword, salt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).WithArgs(driver.Value(authData.Email)).WillReturnRows(row)
			},
			expected: ID,
		},
		{
			name: "get use fail",
			mock: func() {
				row := mock.NewRows([]string{"email", "password", "salt"})
				row.AddRow(ID, "wornPassword", salt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).WithArgs(driver.Value(authData.Email)).WillReturnRows(row)
			},
			expectedError: true,
		},
		{
			name: "no rows returned; func returns error",
			mock: func() {
				row := mock.NewRows([]string{"email", "password", "salt"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).WithArgs(driver.Value(authData.Email)).WillReturnRows(row)
			},
			expectedError: true,
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"email", "password", "salt"})
				row.AddRow(ID, "wornPassword", salt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).WithArgs(driver.Value(authData.Email)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
		{
			name: "scan returns error",
			mock: func() {
				const newArg = 1
				row := mock.NewRows([]string{"email", "password", "salt", "newArg"})
				row.AddRow(ID, "wornPassword", salt, newArg)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).WithArgs(driver.Value(authData.Email)).WillReturnRows(row)
			},
			expectedError: true,
		},
		{
			name: "row.Err() returns error",
			mock: func() {
				row := mock.NewRows([]string{"email", "password", "salt"})
				row.AddRow(ID, dbPassword, salt).RowError(0, errors.New("error"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).WithArgs(driver.Value(authData.Email)).WillReturnRows(row)
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.GetUserByPassword(authData)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestAuthStorage_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	redisDB, _ := redismock.NewClientMock()
	repository := NewAuthStorage(db, redisDB)

	const ID = 1

	tests := []struct {
		name          string
		data          *proto.RegisterData
		mock          func()
		expected      int64
		expectedError bool
	}{
		{
			name: "create user success",
			data: &proto.RegisterData{
				Email:    "testEmail",
				Password: "testPassword",
				Nickname: "testNickname",
			},
			mock: func() {
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users(email, password, nickname, salt, avatar) VALUES($1, $2, $3, $4, $5) RETURNING id`)).WillReturnRows(row)
			},
			expected: ID,
		},
		{
			name: "create user fail",
			data: &proto.RegisterData{
				Email:    "testEmail",
				Password: "testPassword",
				Nickname: "testNickname",
			},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users(email, password, nickname, salt, avatar) VALUES($1, $2, $3, $4, $5) RETURNING id`)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.CreateUser(currentTest.data)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestAuthStorage_IsEmailUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	redisDB, _ := redismock.NewClientMock()
	repository := NewAuthStorage(db, redisDB)
	const (
		email         = "testEmail"
		emailArgument = "testemail"
	)

	tests := []struct {
		name          string
		mock          func()
		expected      bool
		expectedError bool
	}{
		{
			name: "email is unique",
			mock: func() {
				row := mock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).WithArgs(driver.Value(emailArgument)).WillReturnRows(row)
			},
			expected: true,
		},
		{
			name: "email is not unique",
			mock: func() {
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).WithArgs(driver.Value(emailArgument)).WillReturnRows(row)
			},
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"id"})
				row.AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).WithArgs(driver.Value(emailArgument)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.IsEmailUnique(email)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestAuthStorage_IsNicknameUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	redisDB, _ := redismock.NewClientMock()
	repository := NewAuthStorage(db, redisDB)

	const (
		nickname         = "testEmail"
		nicknameArgument = "testemail"
	)

	tests := []struct {
		name          string
		mock          func()
		expected      bool
		expectedError bool
	}{
		{
			name: "nickname is unique",
			mock: func() {
				row := mock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).WithArgs(driver.Value(nicknameArgument)).WillReturnRows(row)
			},
			expected: true,
		},
		{
			name: "nickname is not unique",
			mock: func() {
				row := mock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).WithArgs(driver.Value(nicknameArgument)).WillReturnRows(row)
			},
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"id"})
				row.AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).WithArgs(driver.Value(nicknameArgument)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.IsNicknameUnique(nickname)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestAuthStorage_GetAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	redisDB, _ := redismock.NewClientMock()
	repository := NewAuthStorage(db, redisDB)

	const (
		avatar = "testAvatar"
		userID = 1
	)

	tests := []struct {
		name          string
		mock          func()
		expected      string
		expectedError bool
	}{
		{
			name: "get avatar success",
			mock: func() {
				row := mock.NewRows([]string{"avatar"})
				row.AddRow(avatar)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expected: avatar,
		},
		{
			name: "get avatar fail",
			mock: func() {
				row := mock.NewRows([]string{"avatar"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expected: "",
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"avatar"})
				row.AddRow(avatar)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
		{
			name: "scan returns error",
			mock: func() {
				const newArg = 1
				row := mock.NewRows([]string{"avatar", "newArg"})
				row.AddRow(avatar, newArg)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expectedError: true,
		},
		{
			name: "row.Err() returns error",
			mock: func() {
				row := mock.NewRows([]string{"avatar"})
				row.AddRow(avatar).RowError(0, errors.New("error"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.GetAvatar(userID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}
