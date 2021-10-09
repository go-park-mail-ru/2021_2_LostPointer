package repository

import (
	"2021_2_LostPointer/internal/models"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const SessionTokenLength = 40
const SaltLength = 8

type UserRepositoryRealisation struct {
	userDB 			*sql.DB
	redisConnection *redis.Client
}

func NewUserRepositoryRealization(db *sql.DB, redisConnection *redis.Client) UserRepositoryRealisation {
	return UserRepositoryRealisation{userDB: db, redisConnection: redisConnection}
}

func GetRandomString(l int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(97, 122))
	}
	return string(bytes)
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func StoreSession(redisConnection *redis.Client, session *models.Session) error {
	err := redisConnection.Set(session.Session, session.UserId, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetSessionUser(redisConnection *redis.Client, session string) (int, error) {
	res, err := redisConnection.Get(session).Result()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	id, err := strconv.Atoi(res)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return id, err
}

func GetHash(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (Data UserRepositoryRealisation) CreateUser(userData models.User) (string, error) {
	var id uint64 = 0

	salt := GetRandomString(SaltLength)
	err := Data.userDB.QueryRow(
		`INSERT INTO users(email, password, name, salt) VALUES($1, $2, $3, $4) RETURNING id`,
		userData.Email, GetHash(userData.Password + salt), userData.NickName , salt,
		).Scan(&id)
	if err != nil {
		return "", err
	}
	sessionToken := GetRandomString(SessionTokenLength)
	err = StoreSession(Data.redisConnection, &models.Session{UserId: id, Session: sessionToken})
	if err != nil {
		return "", err
	}

	return sessionToken, err
}

func (Data UserRepositoryRealisation) IsEmailUnique(email string) (bool, error) {
	rows, err := Data.userDB.Query(`SELECT id FROM users WHERE email=$1`, email)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

func (Data UserRepositoryRealisation) IsNicknameUnique(nickname string) (bool, error) {
	rows, err := Data.userDB.Query(`SELECT id FROM users WHERE name=$1`, nickname)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

