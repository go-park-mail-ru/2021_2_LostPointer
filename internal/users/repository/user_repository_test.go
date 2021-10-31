package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestSanitizeUserData(t *testing.T) {
	dangData := models.User{
		Nickname: "<script>nickname</script>",
		Email:    "<script>email</script>",
		Password: "password",
	}
	result := sanitizeUserData(dangData)
	expected := models.User{
		Nickname: "nickname",
		Email:    "email",
		Password: "password",
	}

	assert.Equal(t, expected, result)
}

func TestSanitizeEmail(t *testing.T) {
	dangEmail := "<script>email</script>"
	result := sanitizeEmail(dangEmail)
	expected := "email"

	assert.Equal(t, expected, result)
}

func TestSanitizeNickname(t *testing.T) {
	dangEmail := "<script>nickname</script>"
	result := sanitizeNickname(dangEmail)
	expected := "nickname"

	assert.Equal(t, expected, result)
}


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
		input 		models.Auth
		expected 	uint64
		expectedErr bool
	}{
		{
			name: "User was found in db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id", "password", "salt"}).
					AddRow("1", GetHash("alex1234" + "1234"), "1234")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: models.Auth{Email: "alex1234@gmail.com", Password: "alex1234"},
			expected: 1,
		},
		{
			name: "User was not found in db, wrong email",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id", "password", "salt"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: models.Auth{Email: "alex1234@gmail.com", Password: "alex1234"},
			expected: 0,
		},
		{
			name: "The password in the database did not match the received password",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id", "password", "salt"}).
					AddRow("1", GetHash("alex123" + "1234"), "1234")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: models.Auth{Email: "alex1234@gmail.com", Password: "alex1234"},
			expected: 0,
		},
		{
			name: "Error occurred during SELECT request",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnError(errors.New("Error occurred during request "))
			},
			input: models.Auth{Email: "alex1234@gmail.com", Password: "alex1234"},
			expected: 0,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.DoesUserExist(&testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

//func TestUserRepository_IsEmailUnique(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input   	string
//		expected    bool
//		expectedErr bool
//	}{
//		{
//			name: "Received email is unique",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"id"})
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
//					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
//			},
//			input: "alex1234@gmail.com",
//			expected: true,
//		},
//		{
//			name: "Received email is not unique",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
//					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
//			},
//			input: "alex1234@gmail.com",
//			expected: false,
//		},
//		{
//			name: "Error occurred during SELECT request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
//					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnError(errors.New("Error occurred during request "))
//			},
//			input: "alex1234@gmail.com",
//			expected: false,
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			got, err := r.IsEmailUnique(testCase.input)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, testCase.expected, got)
//			}
//		})
//	}
//}
//
//func TestUserRepository_IsNicknameUnique(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input   	string
//		expected    bool
//		expectedErr bool
//	}{
//		{
//			name: "Received nickname is unique",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"id"})
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
//					WithArgs(driver.Value("alex1234")).WillReturnRows(rows)
//			},
//			input: "alex1234",
//			expected: true,
//		},
//		{
//			name: "Received nickname is not unique",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
//					WithArgs(driver.Value("alex1234")).WillReturnRows(rows)
//			},
//			input: "alex1234",
//			expected: false,
//		},
//		{
//			name: "Error occurred during SELECT request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
//					WithArgs(driver.Value("alex1234")).WillReturnError(errors.New("Error occurred during request "))
//			},
//			input: "alex1234",
//			expected: false,
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			got, err := r.IsNicknameUnique(testCase.input)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, testCase.expected, got)
//			}
//		})
//	}
//}
//
//func TestUserRepository_CreateUser(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input   	models.User
//		expected    uint64
//		expectedErr bool
//	}{
//		{
//			name: "User was successfully created in db",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
//				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users(email, password, nickname, salt, avatar) VALUES($1, $2, $3, $4, $5) RETURNING id`)).
//					WithArgs(
//						driver.Value(strings.ToLower("alex1234@gmail.com")),
//						driver.Value(GetHash("alex1234" + "1234")),
//						driver.Value("alex1234"),
//						driver.Value("1234"),
//						driver.Value(AvatarDefaultFileName),
//					).WillReturnRows(rows)
//			},
//			input: models.User{Email: "alex1234@gmail.com", Password: "alex1234", Nickname: "alex1234"},
//			expected: 1,
//		},
//		{
//			name: "Error occurred during INSERT request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users(email, password, nickname, salt) VALUES($1, $2, $3, $4) RETURNING id`)).
//					WithArgs(
//						driver.Value(strings.ToLower("alex1234@gmail.com")),
//						driver.Value(GetHash("alex1234" + "1234")),
//						driver.Value("alex1234"),
//						driver.Value("1234"),
//					).WillReturnError(errors.New("Error occurred during transaction "))
//			},
//			input: models.User{Email: "alex1234@gmail.com", Password: "alex1234", Nickname: "alex1234"},
//			expected: 0,
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			got, err := r.CreateUser(testCase.input, "1234")
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, testCase.expected, got)
//			}
//		})
//	}
//}
//
//func TestUserRepository_CheckPasswordByUserID(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	type inputStruct struct {
//		ID 		 int
//		Password string
//	}
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input   	inputStruct
//		expected    bool
//		expectedErr bool
//	}{
//		{
//			name: "Received password was found in db",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"password", "salt"}).AddRow(GetHash("alex1234" + "1234"), "1234")
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).
//					WithArgs(driver.Value(1)).WillReturnRows(rows)
//			},
//			input: inputStruct{ID: 1, Password: "alex1234"},
//			expected: true,
//		},
//		{
//			name: "Received password was not found in db",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"password", "salt"})
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).
//					WithArgs(driver.Value(1)).WillReturnRows(rows)
//			},
//			input: inputStruct{ID: 1, Password: "alex1234"},
//			expected: false,
//		},
//		{
//			name: "Error occurred during SELECT request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT password, salt FROM users WHERE id=$1`)).
//					WithArgs(driver.Value(1)).WillReturnError(errors.New("some_error_during_request"))
//			},
//			input: inputStruct{ID: 1, Password: "alex1234"},
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			got, err := r.CheckPasswordByUserID(testCase.input.ID, testCase.input.Password)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, testCase.expected, got)
//			}
//		})
//	}
//}
//
//func TestUserRepository_GetSettings(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input   	int
//		expected    *models.SettingsGet
//		expectedErr bool
//	}{
//		{
//			name: "Settings were successfully returned from db",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"email", "avatar", "nickname"}).
//					AddRow("alex1234@gmail.com", "default", "alex1234")
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).
//					WithArgs(driver.Value(1)).WillReturnRows(rows)
//			},
//			input: 1,
//			expected: &models.SettingsGet{
//				Email: "alex1234@gmail.com",
//				SmallAvatar: "default_150px.webp",
//				BigAvatar: "default_500px.webp",
//				Nickname: "alex1234",
//			},
//		},
//		{
//			name: "No settings were found in db",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"email", "avatar", "nickname"})
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT email, avatar, nickname FROM users WHERE id=$1`)).
//					WithArgs(driver.Value(1)).WillReturnRows(rows)
//			},
//			input: 1,
//			expected: nil,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			got, err := r.GetSettings(testCase.input)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, testCase.expected, got)
//			}
//		})
//	}
//}
//
//func TestUserRepository_UpdateEmail(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	type inputStruct struct {
//		userID int
//		email  string
//	}
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input		*inputStruct
//		expectedErr bool
//	}{
//		{
//			name: "Email was updated successfully",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{})
//				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET email=$1 WHERE id=$2`)).
//					WithArgs(
//						driver.Value(strings.ToLower("alex1234@gmail.com")),
//						driver.Value(1),
//					).WillReturnRows(rows)
//			},
//			input: &inputStruct{userID: 1, email: "alex1234@gmail.com"},
//		},
//		{
//			name: "Error occurred during UPDATE request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET email=$1 WHERE id=$2`)).
//					WithArgs(
//						driver.Value(strings.ToLower("alex1234@gmail.com")),
//						driver.Value(1),
//					).WillReturnError(errors.New("some_error_during_request"))
//			},
//			input: &inputStruct{userID: 1, email: "alex1234@gmail.com"},
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			err := r.UpdateEmail(testCase.input.userID, testCase.input.email)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUserRepository_UpdateNickname(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	type inputStruct struct {
//		userID 	  int
//		nickname  string
//	}
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input		*inputStruct
//		expectedErr bool
//	}{
//		{
//			name: "Nickname was updated successfully",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{})
//				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).
//					WithArgs(
//						driver.Value("alex1234"),
//						driver.Value(1),
//					).WillReturnRows(rows)
//			},
//			input: &inputStruct{userID: 1, nickname: "alex1234"},
//		},
//		{
//			name: "Error occurred during UPDATE request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).
//					WithArgs(
//						driver.Value("alex1234"),
//						driver.Value(1),
//					).WillReturnError(errors.New("some_error_during_request"))
//			},
//			input: &inputStruct{userID: 1, nickname: "alex1234"},
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			err := r.UpdateNickname(testCase.input.userID, testCase.input.nickname)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUserRepository_UpdatePassword(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	type inputStruct struct {
//		userID 	  int
//		password  string
//		salt	  string
//	}
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input		*inputStruct
//		expectedErr bool
//	}{
//		{
//			name: "Email was updated successfully",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{})
//				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`)).
//					WithArgs(
//						driver.Value(GetHash("alex1234" + "1234")),
//						driver.Value("1234"),
//						driver.Value(1),
//					).WillReturnRows(rows)
//			},
//			input: &inputStruct{userID: 1, password: "alex1234", salt: "1234"},
//		},
//		{
//			name: "Error occurred during UPDATE request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`)).
//					WithArgs(
//					driver.Value(GetHash("alex1234" + "1234")),
//					driver.Value("1234"),
//					driver.Value(1),
//				).WillReturnError(errors.New("some_error_during_request"))
//			},
//			input: &inputStruct{userID: 1, password: "alex1234", salt: "1234"},
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			err := r.UpdatePassword(testCase.input.userID, testCase.input.password, testCase.input.salt)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUserRepository_UpdateAvatar(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	type inputStruct struct {
//		userID 	  int
//		filename  string
//	}
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input		*inputStruct
//		expectedErr bool
//	}{
//		{
//			name: "Avatar was successfully updated",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{})
//				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET avatar=$1 WHERE id=$2`)).
//					WithArgs(
//						driver.Value("avatar"),
//						driver.Value(1),
//					).WillReturnRows(rows)
//			},
//			input: &inputStruct{userID: 1, filename: "avatar"},
//		},
//		{
//			name: "Error occurred during UPDATE request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`UPDATE users SET nickname=$1 WHERE id=$2`)).
//					WithArgs(
//						driver.Value("avatar"),
//						driver.Value(1),
//					).WillReturnError(errors.New("some_error_during_request"))
//			},
//			input: &inputStruct{userID: 1, filename: "avatar"},
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			err := r.UpdateAvatar(testCase.input.userID, testCase.input.filename)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUserRepository_GetAvatarFilename(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer func(db *sql.DB) {
//		_ = db.Close()
//	}(db)
//
//	r := NewUserRepository(db)
//
//	tests := []struct {
//		name    	string
//		mock    	func()
//		input		int
//		expected    string
//		expectedErr bool
//	}{
//		{
//			name: "Avatar filename was successfully returned from db",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"avatar"}).
//					AddRow("avatar")
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).
//					WithArgs(driver.Value(1)).WillReturnRows(rows)
//			},
//			input: 1,
//			expected: "avatar",
//		},
//		{
//			name: "User with received id was not found in db",
//			mock: func(){
//				rows := sqlmock.NewRows([]string{"avatar"})
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).
//					WithArgs(driver.Value(1)).WillReturnRows(rows)
//			},
//			input: 1,
//			expected: "",
//		},
//		{
//			name: "Error occurred during SELECT request",
//			mock: func(){
//				mock.ExpectQuery(regexp.QuoteMeta(`SELECT avatar FROM users WHERE id=$1`)).
//					WithArgs(driver.Value(1)).WillReturnError(errors.New("some_error_during_request"))
//			},
//			input: 1,
//			expectedErr: true,
//		},
//	}
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T) {
//			testCase.mock()
//
//			got, err := r.GetAvatarFilename(testCase.input)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, testCase.expected, got)
//			}
//		})
//	}
//}
//
//func TestRedisStore_DeleteSession(t *testing.T) {
//	var err error
//
//	db, mock := redismock.NewClientMock()
//	cookieValue := "feeuhfuy3748478djakdj"
//	mock.Regexp().ExpectDel(cookieValue).RedisNil()
//	r := NewRedisStore(db)
//	r.DeleteSession(cookieValue)
//	if err = mock.ExpectationsWereMet(); err != nil {
//		t.Error("Error occurred during test case", err)
//	}
//
//	assert.NoError(t, err)
//}
//
//func TestRedisStore_StoreSession(t *testing.T) {
//	db, mock := redismock.NewClientMock()
//	defer func(db *redis.Client) {
//		_ = db.Close()
//	}(db)
//
//	sessionToken := GetRandomString(SessionTokenLength)
//	var userID uint64 = 1
//
//	tests := []struct {
//		name 		string
//		mock 		func()
//		input 		uint64
//		expected 	string
//		expectedErr bool
//	}{
//		{
//			name: "Successfully stored sessions in redis",
//			mock: func() {
//				mock.ExpectSet(sessionToken, userID, time.Hour).SetVal("")
//			},
//			input: 1,
//			expected: sessionToken,
//		},
//		{
//			name: "Error occurred in Set method",
//			mock: func() {
//				mock.ExpectSet(sessionToken, userID, time.Hour).SetErr(errors.New("error"))
//			},
//			input: 1,
//			expectedErr: true,
//		},
//	}
//
//	r := NewRedisStore(db)
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T){
//			testCase.mock()
//
//			got, err := r.StoreSession(testCase.input, sessionToken)
//			if testCase.expectedErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, testCase.expected, got)
//			}
//		})
//	}
//}
//
//func TestRedisStore_GetSessionUserId(t *testing.T) {
//	db, mock := redismock.NewClientMock()
//	defer func(db *redis.Client) {
//		_ = db.Close()
//	}(db)
//
//	tests := []struct {
//		name 		string
//		mock 		func()
//		input 		string
//		expected 	int
//		expectedErr bool
//	}{
//		{
//			name: "Successfully returned sessions",
//			mock: func() {
//				mock.ExpectGet("cookie").SetVal("1")
//			},
//			input: "cookie",
//			expected: 1,
//		},
//		{
//			name: "Error occurred in redis Get method",
//			mock: func() {
//				mock.ExpectGet("some_cookie_value").SetErr(errors.New("error"))
//			},
//			input: "some_cookie_value",
//			expected: -1,
//			expectedErr: true,
//		},
//		{
//			name: "Error occurred in redis Get method",
//			mock: func() {
//				mock.ExpectGet("some_cookie_value").RedisNil()
//			},
//			input: "some_cookie_value",
//			expected: 0,
//			expectedErr: true,
//		},
//	}
//
//	r := NewRedisStore(db)
//
//	for _, testCase := range tests {
//		t.Run(testCase.name, func(t *testing.T){
//			testCase.mock()
//
//			got, err := r.GetSessionUserId(testCase.input)
//			if testCase.expectedErr {
//				assert.NotNil(t, err)
//				assert.Equal(t, testCase.expected, got)
//			} else {
//				assert.Nil(t, err)
//				assert.Equal(t, testCase.expected, got)
//			}
//		})
//	}
//}
//
//func TestGetHashReturnsStringWithCorrectLength(t *testing.T) {
//	expected := "b603de426ed0b347d8ca096fb13ba40057d1cb21c9767f231cbb490f09fee088"
//
//	inputStr := "alex1234"
//	got := GetHash(inputStr)
//
//	assert.Equal(t, got, expected)
//}
//
//func TestGetRandomString(t *testing.T) {
//	uniqueStrings := make(map[string]bool)
//
//	const testCasesAmount = 10000
//	const length = 10
//	for i := 0; i < testCasesAmount; i++ {
//		got := GetRandomString(length)
//		uniqueStrings[got] = true
//	}
//
//	assert.Equal(t, len(uniqueStrings), testCasesAmount)
//}
//
//func TestRandInt(t *testing.T) {
//	inputMin, inputMax := -10, 10
//
//	got := RandInt(inputMin, inputMax)
//
//	assert.True(t, got >= inputMin && got <= inputMax)
//}
//
//
