package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/profile/proto"
	"database/sql/driver"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"regexp"
	"testing"
)

func TestUserSettingsStorage_GetSettings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewUserSettingsStorage(db)

	const (
		avatar = "testAvatar"
		userID = 1
	)
	settings := &proto.UserSettings{
		Email:    "testEmail",
		Nickname: "testNickname",
	}
	expectedSettings := &proto.UserSettings{
		Email:       "testEmail",
		Nickname:    "testNickname",
		SmallAvatar: os.Getenv("USERS_ROOT_PREFIX") + avatar + constants.UserAvatarExtension150px,
		BigAvatar:   os.Getenv("USERS_ROOT_PREFIX") + avatar + constants.UserAvatarExtension500px,
	}

	tests := []struct {
		name          string
		mock          func()
		expected      *proto.UserSettings
		expectedError bool
	}{
		{
			name: "get settings success",
			mock: func() {
				row := mock.NewRows([]string{"email", "avatar", "name"})
				row.AddRow(settings.Email, avatar, settings.Nickname)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expected: expectedSettings,
		},
		{
			name: "can not get settings",
			mock: func() {
				row := mock.NewRows([]string{"email", "avatar", "name"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expectedError: true,
		},
		{
			name: "query returns error",
			mock: func() {
				row := mock.NewRows([]string{"email", "avatar", "name"})
				row.AddRow(settings.Email, avatar, settings.Nickname)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
		{
			name: "scan returns error",
			mock: func() {
				const newArg = 1
				row := mock.NewRows([]string{"email", "avatar", "name", "newArg"})
				row.AddRow(settings.Email, avatar, settings.Nickname, newArg)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expectedError: true,
		},
		{
			name: "row.Err() returns error",
			mock: func() {
				row := mock.NewRows([]string{"email", "avatar", "name"}).RowError(0, errors.New("error"))
				row.AddRow(settings.Email, avatar, settings.Nickname)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.GetSettings(userID)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}

func TestUserSettingsStorage_UpdateEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewUserSettingsStorage(db)

	const userID = 1

	tests := []struct {
		name          string
		email         string
		mock          func()
		expectedError bool
	}{
		{
			name:  "update email success",
			email: "test@test.com",
			mock: func() {
				row := mock.NewRows([]string{"success"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET email=$1 WHERE id=$2`)).WithArgs(driver.Value("test@test.com"), driver.Value(userID)).WillReturnRows(row)
			},
		},
		{
			name:  "update email with xss success",
			email: "<script>alert()</script>",
			mock: func() {
				row := mock.NewRows([]string{"success"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET email=$1 WHERE id=$2`)).WithArgs(driver.Value("alert()"), driver.Value(userID)).WillReturnRows(row)
			},
		},
		{
			name:  "query returns error",
			email: "test@test.com",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET email=$1 WHERE id=$2`)).WithArgs(driver.Value("test@test.com"), driver.Value(userID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.UpdateEmail(userID, currentTest.email)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserSettingsStorage_UpdateNickname(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewUserSettingsStorage(db)

	const userID = 1

	tests := []struct {
		name          string
		nickname      string
		mock          func()
		expectedError bool
	}{
		{
			name:     "update nickname success",
			nickname: "test",
			mock: func() {
				row := mock.NewRows([]string{"success"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).WithArgs(driver.Value("test"), driver.Value(userID)).WillReturnRows(row)
			},
		},
		{
			name:     "update nickname with xss success",
			nickname: "<script>alert()</script>",
			mock: func() {
				row := mock.NewRows([]string{"success"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).WithArgs(driver.Value("alert()"), driver.Value(userID)).WillReturnRows(row)
			},
		},
		{
			name:     "query returns error",
			nickname: "test",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).WithArgs(driver.Value("test"), driver.Value(userID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.UpdateNickname(userID, currentTest.nickname)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserSettingsStorage_UpdateAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewUserSettingsStorage(db)

	const userID = 1

	tests := []struct {
		name          string
		filename      string
		mock          func()
		expectedError bool
	}{
		{
			name:     "update avatar success",
			filename: "test",
			mock: func() {
				row := mock.NewRows([]string{"success"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET avatar=$1 WHERE id=$2`)).WithArgs(driver.Value("test"), driver.Value(userID)).WillReturnRows(row)
			},
		},
		{
			name:     "query returns error",
			filename: "test",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET avatar=$1 WHERE id=$2`)).WithArgs(driver.Value("test"), driver.Value(userID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.UpdateAvatar(userID, currentTest.filename)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserSettingsStorage_UpdatePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewUserSettingsStorage(db)

	const (
		userID = 1
	)

	tests := []struct {
		name          string
		password      string
		mock          func()
		expectedError bool
	}{
		{
			name:     "update password success",
			password: "test",
			mock: func() {
				row := mock.NewRows([]string{"success"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`)).WillReturnRows(row)
			},
		},
		{
			name:     "query returns error",
			password: "test",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			err := repository.UpdatePassword(userID, currentTest.password)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserSettingsStorage_IsEmailUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewUserSettingsStorage(db)

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

func TestUserSettingsStorage_IsNicknameUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewUserSettingsStorage(db)

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

func TestUserSettingsStorage_CheckPasswordByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	repository := NewUserSettingsStorage(db)

	const (
		userID      = 1
		oldPassword = "testPassword"
		salt        = "testSalt"
	)
	password, _ := bcrypt.GenerateFromPassword([]byte(oldPassword+salt), bcrypt.DefaultCost)

	tests := []struct {
		name          string
		mock          func()
		expected      bool
		expectedError bool
	}{
		{
			name: "password exists",
			mock: func() {
				row := mock.NewRows([]string{"password", "salt"}).AddRow(password, salt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expected: true,
		},
		{
			name: "can not find password in database",
			mock: func() {
				row := mock.NewRows([]string{"password", "salt"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
		},
		{
			name: "query returns error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnError(errors.New("error"))
			},
			expectedError: true,
		},
		{
			name: "wrong credentials",
			mock: func() {
				row := mock.NewRows([]string{"password", "salt"}).AddRow("wrongCredentials", salt)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expectedError: true,
		},
		{
			name: "scan returns error",
			mock: func() {
				const newArg = 1
				row := mock.NewRows([]string{"password", "salt", "newArg"}).AddRow(password, salt, newArg)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expectedError: true,
		},
		{
			name: "row.Err() returns error",
			mock: func() {
				const newArg = 1
				row := mock.NewRows([]string{"password", "salt", "newArg"}).AddRow(password, salt, newArg).RowError(0, errors.New("error"))
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).WithArgs(driver.Value(userID)).WillReturnRows(row)
			},
			expectedError: true,
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			currentTest.mock()
			result, err := repository.CheckPasswordByUserID(userID, oldPassword)
			if currentTest.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, result)
			}
		})
	}
}
