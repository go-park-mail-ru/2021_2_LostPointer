package repository

import (
	"2021_2_LostPointer/internal/microservices/authorization/proto"
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/kennygrant/sanitize"
	"golang.org/x/crypto/bcrypt"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/errors"
	"2021_2_LostPointer/pkg/utils"
)

type AuthStorage struct {
	db    *sql.DB
	redis *redis.Client
}

func NewAuthStorage(db *sql.DB, redis *redis.Client) *AuthStorage {
	return &AuthStorage{db: db, redis: redis}
}

func (storage *AuthStorage) CreateSession(id int64, cookieValue string) error {
	err := storage.redis.Set(context.Background(), cookieValue, id, constants.CookieLifetime).Err()
	if err != nil {
		return err
	}
	return nil
}

func (storage *AuthStorage) GetUserByCookie(cookieValue string) (int64, error) {
	idStr, err := storage.redis.Get(context.Background(), cookieValue).Result()
	id, _ := strconv.Atoi(idStr)
	return int64(id), err
}

func (storage *AuthStorage) DeleteSession(cookieValue string) error {
	err := storage.redis.Del(context.Background(), cookieValue).Err()
	return err
}

func (storage *AuthStorage) GetUserByPassword(authData *proto.AuthData) (int64, error) {
	query := `SELECT id, password, salt FROM users WHERE email=$1`
	rows, err := storage.db.Query(query, authData.Email)
	if err != nil {
		return 0, err
	}
	err = rows.Err()
	if err != nil {
		return 0, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	if !rows.Next() {
		return 0, errors.ErrWrongCredentials
	}
	var (
		dbUserID           int
		dbPassword, dbSalt string
	)
	if err = rows.Scan(&dbUserID, &dbPassword, &dbSalt); err != nil {
		return 0, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(authData.Password+dbSalt)); err != nil {
		return 0, errors.ErrWrongCredentials
	}
	return int64(dbUserID), nil
}

func (storage *AuthStorage) CreateUser(registerData *proto.RegisterData) (int64, error) {
	query := `INSERT INTO users(email, password, nickname, salt, avatar) VALUES($1, $2, $3, $4, $5) RETURNING id`

	salt, err := utils.GetRandomString(constants.SaltLength)
	if err != nil {
		return 0, err
	}
	sanitizedEmail := SanitizeEmail(registerData.Email)
	sanitizedNickname := SanitizeNickname(registerData.Nickname)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password+salt), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int64
	err = storage.db.QueryRow(query, strings.ToLower(sanitizedEmail), hashedPassword, sanitizedNickname,
		salt, constants.AvatarDefaultFileName).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (storage *AuthStorage) IsEmailUnique(email string) (bool, error) {
	query := `SELECT id FROM users WHERE lower(email)=$1`

	rows, err := storage.db.Query(query, strings.ToLower(email))
	if err != nil {
		return false, err
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

func (storage *AuthStorage) IsNicknameUnique(nickname string) (bool, error) {
	query := `SELECT id FROM users WHERE lower(nickname)=$1`

	rows, err := storage.db.Query(query, strings.ToLower(nickname))
	if err != nil {
		return false, err
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()
	if rows.Next() {
		return false, nil
	}

	return true, nil
}

func (storage *AuthStorage) GetAvatar(userID int64) (string, error) {
	query := `SELECT avatar FROM users WHERE id=$1`

	rows, err := storage.db.Query(query, userID)
	if err != nil {
		return "", err
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	if !rows.Next() {
		return "", nil
	}

	var filename string
	if err := rows.Scan(&filename); err != nil {
		return "", err
	}

	return filename, nil
}

func SanitizeEmail(email string) string {
	return sanitize.HTML(email)
}

func SanitizeNickname(nickname string) string {
	return sanitize.HTML(nickname)
}
