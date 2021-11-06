package repository

import (
	"2021_2_LostPointer/internal/constants"
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type AuthStorage struct {
	db *sql.DB
	redis *redis.Client
}

func NewAuthStorage(db *sql.DB, redis *redis.Client) AuthStorage {
	return AuthStorage{db: db, redis: redis}
}

func (storage *AuthStorage) CreateSession(ID int64, cookieValue string) error {
	err := storage.redis.Set(context.Background(), cookieValue, ID, constants.CookieLifetime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (storage *AuthStorage) GetUserByPassword(login string, password string) (int64, error) {
	query := `SELECT id, password, salt FROM users WHERE email=$1`
	rows, err := storage.db.Query(query, login)
	if err != nil {
		return 0, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if !rows.Next() {
		return 0, nil
	}
	var (
		dbUserID int
		dbPassword, dbSalt string
	)
	if err := rows.Scan(&dbUserID, &dbPassword, &dbSalt); err != nil {
		return 0, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password + dbSalt)); err != nil {
		return 0, nil
	}
	return int64(dbUserID), nil
}
