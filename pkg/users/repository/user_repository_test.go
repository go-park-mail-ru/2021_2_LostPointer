package repository

import (
	"2021_2_LostPointer/pkg/models"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestUserRepository_DoesUserExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	tests := []struct {
		name 		string
		mock 		func()
		input 		models.Auth
		expected 	uint64
		expectedErr bool
	}{
		{
			name: "User exists in database",
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
			name: "Wrong email",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id", "password", "salt"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: models.Auth{Email: "alex1234@gmail.com", Password: "alex1234"},
			expected: 0,
		},
		{
			name: "Wrong password",
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
			name: "Func returns error",
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
	defer db.Close()

	r := NewUserRepository(db)

	tests := []struct {
		name    	string
		mock    	func()
		input   	string
		expected    bool
		expectedErr bool
	}{
		{
			name: "Email is unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: "alex1234@gmail.com",
			expected: true,
		},
		{
			name: "Email is not unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: "alex1234@gmail.com",
			expected: false,
		},
		{
			name: "Func returns error",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnError(errors.New("Error occurred during request "))
			},
			input: "alex1234@gmail.com",
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
	defer db.Close()

	r := NewUserRepository(db)

	tests := []struct {
		name    	string
		mock    	func()
		input   	string
		expected    bool
		expectedErr bool
	}{
		{
			name: "Nickname is unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
					WithArgs(driver.Value("alex1234")).WillReturnRows(rows)
			},
			input: "alex1234",
			expected: true,
		},
		{
			name: "Nickname is not unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
					WithArgs(driver.Value("alex1234")).WillReturnRows(rows)
			},
			input: "alex1234",
			expected: false,
		},
		{
			name: "Func returns error",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
					WithArgs(driver.Value("alex1234")).WillReturnError(errors.New("Error occurred during request "))
			},
			input: "alex1234",
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
	defer db.Close()

	r := NewUserRepository(db)

	tests := []struct {
		name    	string
		mock    	func()
		input   	models.User
		expected    uint64
		expectedErr bool
	}{
		{
			name: "Add user to db",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"}).AddRow("1")
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users(email, password, nickname, salt) VALUES($1, $2, $3, $4) RETURNING id`)).
					WithArgs(
						driver.Value(strings.ToLower("alex1234@gmail.com")),
						driver.Value(GetHash("alex1234" + "1234")),
						driver.Value("alex1234"),
						driver.Value("1234")).WillReturnRows(rows)
			},
			input: models.User{Email: "alex1234@gmail.com", Password: "alex1234", Nickname: "alex1234"},
			expected: 1,
		},
		{
			name: "Func returns error",
			mock: func(){
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users(email, password, nickname, salt) VALUES($1, $2, $3, $4) RETURNING id`)).
					WithArgs(
						driver.Value(strings.ToLower("alex1234@gmail.com")),
						driver.Value(GetHash("alex1234" + "1234")),
						driver.Value("alex1234"),
						driver.Value("1234")).WillReturnError(errors.New("Error occurred during transaction "))
			},
			input: models.User{Email: "alex1234@gmail.com", Password: "alex1234", Nickname: "alex1234"},
			expected: 0,
			expectedErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()

			got, err := r.CreateUser(testCase.input, "1234")
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestRedisStore_DeleteSession(t *testing.T) {
	var err error

	db, mock := redismock.NewClientMock()

	cookieValue := "feeuhfuy3748478djakdj"

	mock.Regexp().ExpectDel(cookieValue).RedisNil()

	r := NewRedisStore(db)

	r.DeleteSession(cookieValue)

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Error("Error occurred during test case", err)
	}

	assert.NoError(t, err)
}

func TestRedisStore_StoreSession(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	sessionToken := GetRandomString(40)
	var userID uint64 = 1
	mock.ExpectSet(sessionToken, userID, time.Hour).SetVal("")

	r := NewRedisStore(db)

	_, err := r.StoreSession(userID, sessionToken)
	assert.NoError(t, err)
}

func TestRedisStore_GetSessionUserId(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	tests := []struct {
		name 		string
		mock 		func()
		input 		string
		expected 	int
		expectedErr bool
	}{
		{
			name: "Session exists",
			mock: func() {
				mock.ExpectGet("some_cookie_value").SetVal("1")
			},
			input: "some_cookie_value",
			expected: 1,
		},
		{
			name: "Func returns error",
			mock: func() {
				mock.ExpectGet("some_cookie_value").RedisNil()
			},
			input: "some_cookie_value",
			expected: 0,
			expectedErr: true,
		},
	}

	r := NewRedisStore(db)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T){
			testCase.mock()

			got, err := r.GetSessionUserId(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}

func TestGetHashReturnsStringWithCorrectLength(t *testing.T) {
	expected := "b603de426ed0b347d8ca096fb13ba40057d1cb21c9767f231cbb490f09fee088"

	inputStr := "alex1234"
	got := GetHash(inputStr)

	assert.Equal(t, got, expected)
}

func TestGetRandomString(t *testing.T) {
	uniqueStrings := make(map[string]bool)

	const testCasesAmount = 10000
	const length = 10
	for i := 0; i < testCasesAmount; i++ {
		got := GetRandomString(length)
		uniqueStrings[got] = true
	}

	assert.Equal(t, len(uniqueStrings), testCasesAmount)
}

func TestRandInt(t *testing.T) {
	inputMin, inputMax := -10, 10

	got := RandInt(inputMin, inputMax)

	assert.True(t, got >= inputMin && got <= inputMax)
}
