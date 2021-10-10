package repository

import (
	"2021_2_LostPointer/internal/models"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/go-redis/redis"
	"math/rand"
	"strings"
	"time"
)

const SessionTokenLength = 40
const SaltLength = 8

type UserRepository struct {
	userDB 			*sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{userDB: db}
}

func GetRandomString(l int) string {
	validCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = validCharacters[RandInt(0, len(validCharacters) - 1)]
	}
	return string(bytes)
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func StoreSession(redisConnection *redis.Client, session *models.Session) error {
	err := redisConnection.Set(session.Session, session.UserID, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (Data UserRepository) CreateUser(userData models.User) (uint64, error) {
	var id uint64

	salt := GetRandomString(SaltLength)
	err := Data.userDB.QueryRow(
		`INSERT INTO users(email, password, nickname, salt) VALUES($1, $2, $3, $4) RETURNING id`,
		strings.ToLower(userData.Email), GetHash(userData.Password + salt), userData.Nickname, salt,
		).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (Data UserRepository) UserExits(authData models.Auth) (uint64, error) {
	var id uint64
	var password, salt string

	rows, err := Data.userDB.Query(`SELECT id, password, salt FROM users WHERE email=$1`, authData.Email)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	// Пользователь с таким email нет в базе
	if !rows.Next() {
		return 0, nil
	}
	if err := rows.Scan(&id, &password, &salt); err != nil {
		return 0, nil
	}
	// Не совпадает пароль
	if GetHash(authData.Password + salt) != password {
		return 0, nil
	}

	return id, err
}

func (Data UserRepository) IsEmailUnique(email string) (bool, error) {
	rows, err := Data.userDB.Query(`SELECT id FROM users WHERE lower(email)=$1`, strings.ToLower(email))
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

func (Data UserRepository) IsNicknameUnique(nickname string) (bool, error) {
	rows, err := Data.userDB.Query(`SELECT id FROM users WHERE lower(nickname)=$1`, strings.ToLower(nickname))

	if err != nil {
		return false, err
	}
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

