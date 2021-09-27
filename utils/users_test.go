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
		Username: "alex",
		Password: "1234",
		Salt: GetRandomString(SaltLength),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, username, password, salt
		FROM users
		WHERE username=$1
	`)).WithArgs(
		driver.Value(user.Username),
	).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"id", "username", "password", "salt"})
		rr.AddRow(user.ID, user.Username, GetHash(user.Password + user.Salt), user.Salt)
		return rr
	}())

	resultId, err := UserExistsLogin(db, user)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, user)
	}
	assert.Equal(t, user.ID, resultId)
}

func TestIsUserUnique(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Data for testing
	user := models.User{
		ID: 1,
		Username: "alex",
		Password: "1234",
		Salt: GetRandomString(SaltLength),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id FROM users WHERE username=$1
	`)).WithArgs(
		driver.Value(user.Username),
	).WillReturnRows(func() *sqlmock.Rows {
		rr := sqlmock.NewRows([]string{"id"})
		rr.AddRow(user.ID)
		return rr
	}())

	exists, err := IsUserUnique(db, user)
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
		Username: "alex",
		Password: "1234",
		Salt: GetRandomString(SaltLength),
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users(username, password, salt)
		VALUES($1, $2, $3)
		RETURNING id
	`)).WithArgs(
		driver.Value(user.Username),
		driver.Value(GetHash(user.Password + user.Salt)),
		driver.Value(user.Salt),
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
