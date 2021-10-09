package utils

import (
	"2021_2_LostPointer/models"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestUserExistsLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Data for testing
	user := models.User{
		ID: 1,
		Email: "alex",
		Password: "1234",
		Salt: GetRandomString(SaltLength),
		Nickname: "stas_gena_turbo",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, email, password, salt
		FROM users
		WHERE email=$1
	`)).WithArgs(
		driver.Value(user.Email),
	).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"id", "username", "password", "salt"})
		rr.AddRow(user.ID, user.Email, GetHash(user.Password + user.Salt), user.Salt)
		return rr
	}())

	resultId, err := UserExistsLogin(db, user)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.Equal(t, user.ID, resultId)
}

func TestIsUserEmailUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Data for testing
	user := models.User{
		ID: 1,
		Email: "alex",
		Password: "1234",
		Salt: GetRandomString(SaltLength),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id FROM users WHERE email=$1
	`)).WithArgs(
		driver.Value(user.Email),
	).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"id"})
		rr.AddRow(user.ID)
		return rr
	}())

	exists, err := IsUserEmailUnique(db, user.Email)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.Equal(t, false, exists)
}

func TestIsUserNicknameUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Data for testing
	user := models.User{
		ID: 1,
		Email: "alex",
		Password: "1234",
		Salt: GetRandomString(SaltLength),
		Nickname: "stas_gena_turbo",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id FROM users WHERE nickname=$1
	`)).WithArgs(
		driver.Value(user.Nickname),
	).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"id"})
		rr.AddRow(user.ID)
		return rr
	}())

	exists, err := IsUserNicknameUnique(db, user.Nickname)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.Equal(t, false, exists)
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Data for testing
	user := models.User{
		ID: 1,
		Email: "alex",
		Password: "1234",
		Salt: GetRandomString(SaltLength),
		Nickname: "stas_gena_turbo",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users(email, password, salt, nickname)
		VALUES($1, $2, $3, $4)
		RETURNING id
	`)).WithArgs(
		driver.Value(user.Email),
		driver.Value(GetHash(user.Password + user.Salt)),
		driver.Value(user.Salt),
		driver.Value(user.Nickname),
	).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"id"})
		rr.AddRow(user.ID)
		return rr
	}())

	resultID, err := CreateUser(db, user, user.Salt)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.Equal(t, user.ID, resultID)
}
