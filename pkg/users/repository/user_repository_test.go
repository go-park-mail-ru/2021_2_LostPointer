package repository

import (
	"2021_2_LostPointer/pkg/models"
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestUserRepository_UserExits(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewUserRepository(db)

	tests := []struct {
		name 	string
		mock 	func()
		input 	models.Auth
		want 	uint64
		wantErr bool
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
			want: 1,
		},
		{
			name: "Wrong email",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id", "password", "salt"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, password, salt FROM users WHERE email=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: models.Auth{Email: "alex1234@gmail.com", Password: "alex1234"},
			want: 0,
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
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.UserExits(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
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
		name    string
		mock    func()
		input   string
		want    bool
		wantErr bool
	}{
		{
			name: "Email is unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: "alex1234@gmail.com",
			want: true,
		},
		{
			name: "Email is not unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(email)=$1`)).
					WithArgs(driver.Value("alex1234@gmail.com")).WillReturnRows(rows)
			},
			input: "alex1234@gmail.com",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.IsEmailUnique(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
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
		name    string
		mock    func()
		input   string
		want    bool
		wantErr bool
	}{
		{
			name: "Nickname is unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
					WithArgs(driver.Value("alex1234")).WillReturnRows(rows)
			},
			input: "alex1234",
			want: true,
		},
		{
			name: "Nickname is not unique",
			mock: func(){
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id FROM users WHERE lower(nickname)=$1`)).
					WithArgs(driver.Value("alex1234")).WillReturnRows(rows)
			},
			input: "alex1234",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.IsNicknameUnique(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
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
		name    string
		mock    func()
		input   models.User
		want    uint64
		wantErr bool
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
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CreateUser(tt.input, "1234")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
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
		name string
		mock func()
		input string
		want int
		wantErr bool
	}{
		{
			name: "Session exists",
			mock: func() {
				mock.ExpectGet("some_cookie_value").SetVal("1")
			},
			input: "some_cookie_value",
			want: 1,
		},
		{
			name: "Session does not exist",
			mock: func() {
				mock.ExpectGet("some_cookie_value").RedisNil()
			},
			input: "some_cookie_value",
			want: 0,
			wantErr: true,
		},
	}

	r := NewRedisStore(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			tt.mock()

			got, err := r.GetSessionUserId(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetHash(t *testing.T) {
	expected := "b603de426ed0b347d8ca096fb13ba40057d1cb21c9767f231cbb490f09fee088"

	inputStr := "alex1234"
	got := GetHash(inputStr)

	assert.Equal(t, got, expected)
}

func TestGetRandomString(t *testing.T) {
	expectedLength := 10

	inputLength := 10
	got := GetRandomString(inputLength)

	assert.Equal(t, len(got), expectedLength)
}

func TestRandInt(t *testing.T) {
	inputMin, inputMax := -10, 10

	got := RandInt(inputMin, inputMax)

	assert.True(t, got >= inputMin && got <= inputMax)
}

func TestNewUserRepository(t *testing.T) {
	db := new(sql.DB)
	r := NewUserRepository(db)
	assert.Equal(t, r.userDB, db)
}
