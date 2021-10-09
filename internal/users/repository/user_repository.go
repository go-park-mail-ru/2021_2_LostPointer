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
	redisConnection *redis.Client
}

func NewUserRepository(db *sql.DB, redisConnection *redis.Client) UserRepository {
	return UserRepository{userDB: db, redisConnection: redisConnection}
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

func (Data UserRepository) CreateUser(userData models.User) (string, error) {
	var id uint64

	salt := GetRandomString(SaltLength)
	err := Data.userDB.QueryRow(
		`INSERT INTO users(email, password, name, salt) VALUES($1, $2, $3, $4) RETURNING id`,
		strings.ToLower(userData.Email), GetHash(userData.Password + salt), userData.Nickname, salt,
		).Scan(&id)
	if err != nil {
		return "", err
	}
	sessionToken := GetRandomString(SessionTokenLength)
	err = StoreSession(Data.redisConnection, &models.Session{UserID: id, Session: sessionToken})
	if err != nil {
		return "", err
	}

	return sessionToken, err
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
	rows, err := Data.userDB.Query(`SELECT id FROM users WHERE lower(name)=$1`, strings.ToLower(nickname))

	if err != nil {
		return false, err
	}
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

