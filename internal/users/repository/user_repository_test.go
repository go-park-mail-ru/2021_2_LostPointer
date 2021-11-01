package repository

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/utils/constants"
	"2021_2_LostPointer/internal/utils/hash"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
)

func TestUserRepository_DoesUserExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	tests := []struct {
		name 		string
		mock 		func()
		input 		*models.Auth
		expected 	int
		expectedErr bool
	}{
		{
			name: "User was found in db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id", "password", "salt"}).
					AddRow("1", hash.GetHash("JesusLoveMe" + "1337"), "1337")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("LaHaine@gmail.com")).WillReturnRows(rows)
			},
			input: &models.Auth{Email: "LaHaine@gmail.com", Password: "JesusLoveMe"},
			expected: 1,
		},
		{
			name: "User was not found in db, wrong email",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id", "password", "salt"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("LaHaine@gmail.com")).WillReturnRows(rows)
			},
			input: &models.Auth{Email: "LaHaine@gmail.com", Password: "JesusLoveMe"},
			expected: 0,
		},
		{
			name: "The password in the database did not match the received password",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id", "password", "salt"}).
					AddRow("1", hash.GetHash("JesusLoveMe" + "1337"), "1337")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("LaHaine@gmail.com")).WillReturnRows(rows)
			},
			input: &models.Auth{Email: "LaHaine@gmail.com", Password: "JesusLoveMe1488"},
			expected: 0,
		},
		{
			name: "Error occurred during SELECT request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("LaHaine@gmail.com")).WillReturnError(errors.New("Error occurred during request"))
			},
			input: &models.Auth{Email: "LaHaine@gmail.com", Password: "JesusLoveMe"},
			expected: 0,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.DoesUserExist(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserRepository_IsEmailUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	tests := []struct {
		name    	string
		mock    	func()
		input   	string
		expected    bool
		expectedErr bool
	}{
		{
			name: "Received email is unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
					WithArgs(driver.Value(strings.ToLower("LaHaine@gmail.com"))).WillReturnRows(rows)
			},
			input: "LaHaine@gmail.com",
			expected: true,
		},
		{
			name: "Received email is not unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
					WithArgs(driver.Value(strings.ToLower("LaHaine@gmail.com"))).WillReturnRows(rows)
			},
			input: "LaHaine@gmail.com",
			expected: false,
		},
		{
			name: "Error occurred during SELECT request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
					WithArgs(driver.Value(strings.ToLower("LaHaine@gmail.com"))).WillReturnError(errors.New("error"))
			},
			input: "LaHaine@gmail.com",
			expected: false,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.IsEmailUnique(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserRepository_IsNicknameUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	tests := []struct {
		name    	string
		mock    	func()
		input   	string
		expected    bool
		expectedErr bool
	}{
		{
			name: "Received nickname is unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
					WithArgs(driver.Value(strings.ToLower("LaHaine"))).WillReturnRows(rows)
			},
			input: "LaHaine",
			expected: true,
		},
		{
			name: "Received nickname is not unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
					WithArgs(driver.Value(strings.ToLower("LaHaine"))).WillReturnRows(rows)
			},
			input: "LaHaine",
			expected: false,
		},
		{
			name: "Error occurred during SELECT request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
					WithArgs(driver.Value(strings.ToLower("LaHaine"))).WillReturnError(errors.New("error"))
			},
			input: "LaHaine",
			expected: false,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.IsNicknameUnique(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	tests := []struct {
		name    	string
		mock    	func()
		input   	*models.User
		expected    int
		expectedErr bool
	}{
		{
			name: "Error occurred during INSERT request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users(email, password, nickname, salt, avatar) VALUES($1, $2, $3, $4) RETURNING id`)).
					WithArgs(
						driver.Value(strings.ToLower("LaHaine@gmail.com")),
						driver.Value(hash.GetHash("JesusLovesMe" + "1337")),
						driver.Value("LaHaine"),
						driver.Value("1337"),
						driver.Value(constants.AvatarDefaultFileName),
					).WillReturnError(errors.New("Error occurred during request"))
			},
			input: &models.User{Email: "LaHaine@gmail.com", Password: "JesusLovesMe", Nickname: "LaHaine"},
			expected: 0,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.CreateUser(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserRepository_CheckPasswordByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	type inputStruct struct {
		ID 		 int
		Password string
	}

	tests := []struct {
		name    	string
		mock    	func()
		input   	inputStruct
		expected    bool
		expectedErr bool
	}{
		{
			name: "Received password was found in db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"password", "salt"}).AddRow(hash.GetHash("JesusLovesMe" + "1337"), "1337")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			input: inputStruct{ID: 1, Password: "JesusLovesMe"},
			expected: true,
		},
		{
			name: "Received password was not found in db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"password", "salt"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			input: inputStruct{ID: 1, Password: "JesusLovesMe"},
			expected: false,
		},
		{
			name: "Error occurred during SELECT request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnError(errors.New("error"))
			},
			input: inputStruct{ID: 1, Password: "JesusLovesMe"},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.CheckPasswordByUserID(testCase.input.ID, testCase.input.Password)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserRepository_GetSettings(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	tests := []struct {
		name    	string
		mock    	func()
		input   	int
		expected    *models.SettingsGet
		expectedErr bool
	}{
		{
			name: "Settings were successfully returned from db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"email", "avatar", "nickname"}).
					AddRow("lahaine@gmail.com", "default", "LaHaine")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			input: 1,
			expected: &models.SettingsGet{
				Email: "lahaine@gmail.com",
				SmallAvatar: "default_150px.webp",
				BigAvatar: "default_500px.webp",
				Nickname: "LaHaine",
			},
		},
		{
			name: "No settings were found in db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"email", "avatar", "nickname"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			input: 1,
			expected: nil,
		},
		{
			name: "Error occurred during SELECT request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnError(errors.New("error"))
			},
			input: 1,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.GetSettings(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestUserRepository_UpdateEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	type inputStruct struct {
		userID int
		email  string
	}

	tests := []struct {
		name    	string
		mock    	func()
		input		*inputStruct
		expectedErr bool
	}{
		{
			name: "Email was updated successfully",
			mock: func(){
				rows := sqlmock.NewRows([]string{})
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET email=$1 WHERE id=$2`)).
					WithArgs(
						driver.Value(strings.ToLower("LaHaine@gmail.com")),
						driver.Value(1),
					).WillReturnRows(rows)
			},
			input: &inputStruct{userID: 1, email: "LaHaine@gmail.com"},
		},
		{
			name: "Error occurred during UPDATE request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET email=$1 WHERE id=$2`)).
					WithArgs(
						driver.Value(strings.ToLower("LaHaine@gmail.com")),
						driver.Value(1),
					).WillReturnError(errors.New("error"))
			},
			input: &inputStruct{userID: 1, email: "alex1234@gmail.com"},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			err := r.UpdateEmail(testCase.input.userID, testCase.input.email)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_UpdateNickname(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	type inputStruct struct {
		userID 	  int
		nickname  string
	}

	tests := []struct {
		name    	string
		mock    	func()
		input		*inputStruct
		expectedErr bool
	}{
		{
			name: "Nickname was updated successfully",
			mock: func(){
				rows := sqlmock.NewRows([]string{})
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).
					WithArgs(
						driver.Value("LaHaine"),
						driver.Value(1),
					).WillReturnRows(rows)
			},
			input: &inputStruct{userID: 1, nickname: "LaHaine"},
		},
		{
			name: "Error occurred during UPDATE request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).
					WithArgs(
						driver.Value("LaHaine"),
						driver.Value(1),
					).WillReturnError(errors.New("error"))
			},
			input: &inputStruct{userID: 1, nickname: "LaHaine"},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			err := r.UpdateNickname(testCase.input.userID, testCase.input.nickname)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_UpdatePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	type inputStruct struct {
		userID 	  int
		password  string
	}

	tests := []struct {
		name    	string
		mock    	func()
		input		*inputStruct
		expectedErr bool
	}{
		//{
		//	name: "Password was updated successfully",
		//	mock: func(){
		//		rows := sqlmock.NewRows([]string{})
		//		mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`)).
		//			WithArgs(
		//				driver.Value(hash.GetHash("JesusLovesMe" + "1337")),
		//				driver.Value("1337"),
		//				driver.Value(1),
		//			).WillReturnRows(rows)
		//	},
		//	input: &inputStruct{userID: 1, password: "JesusLovesMe"},
		//},
		{
			name: "Error occurred during UPDATE request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`)).
					WithArgs(
					driver.Value(hash.GetHash("JesusLovesMe" + "1337")),
					driver.Value("1337"),
					driver.Value(1),
				).WillReturnError(errors.New("some_error_during_request"))
			},
			input: &inputStruct{userID: 1, password: "JesusLovesMe"},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			err := r.UpdatePassword(testCase.input.userID, testCase.input.password)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_UpdateAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	type inputStruct struct {
		userID 	  int
		filename  string
	}

	tests := []struct {
		name    	string
		mock    	func()
		input		*inputStruct
		expectedErr bool
	}{
		{
			name: "Avatar was successfully updated",
			mock: func(){
				rows := sqlmock.NewRows([]string{})
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET avatar=$1 WHERE id=$2`)).
					WithArgs(
						driver.Value("avatar"),
						driver.Value(1),
					).WillReturnRows(rows)
			},
			input: &inputStruct{userID: 1, filename: "avatar"},
		},
		{
			name: "Error occurred during UPDATE request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).
					WithArgs(
						driver.Value("avatar"),
						driver.Value(1),
					).WillReturnError(errors.New("some_error_during_request"))
			},
			input: &inputStruct{userID: 1, filename: "avatar"},
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			err := r.UpdateAvatar(testCase.input.userID, testCase.input.filename)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepository_GetAvatarFilename(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	r := NewUserRepository(db)

	tests := []struct {
		name    	string
		mock    	func()
		input		int
		expected    string
		expectedErr bool
	}{
		{
			name: "Avatar filename was successfully returned from db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"avatar"}).
					AddRow("avatar")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			input: 1,
			expected: "avatar",
		},
		{
			name: "User with received id was not found in db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"avatar"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			input: 1,
		},
		{
			name: "Error occurred in rows.Scan",
			mock: func(){
				rows := sqlmock.NewRows([]string{"avatar", "excessArg"}).
					AddRow("avatar", "someArg")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			input: 1,
			expected: "avatar",
			expectedErr: true,
		},
		{
			name: "Error occurred during SELECT request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).
					WithArgs(driver.Value(1)).WillReturnError(errors.New("some_error_during_request"))
			},
			input: 1,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.GetAvatarFilename(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}