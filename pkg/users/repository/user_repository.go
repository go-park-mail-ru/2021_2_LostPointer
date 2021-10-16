package repository

import (
	"2021_2_LostPointer/pkg/models"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/chai2010/webp"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sunshineplan/imgconv"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

const SaltLength = 8
const SessionTokenLength = 40
const AvatarWidthBig = 500
const AvatarWidthLittle = 150
const AvatarDefaultPath = "placeholder"

var ctx = context.Background()

type UserRepository struct {
	userDB 	*sql.DB
}

type RedisStore struct {
	redisConnection *redis.Client
}

type FileSystem struct {}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{userDB: db}
}

func NewRedisStore(redisConnection *redis.Client) RedisStore {
	return RedisStore{
		redisConnection: redisConnection,
	}
}

func NewFileSystem() FileSystem {
	return FileSystem{}
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

func GetHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (Data UserRepository) CreateUser(userData models.User, customSalt ...string) (uint64, error) {
	var id uint64
	var salt string

	if len(customSalt) != 0 {
		salt = customSalt[0]
	} else {
		salt = GetRandomString(SaltLength)
	}
	err := Data.userDB.QueryRow(
		`INSERT INTO users(email, password, nickname, salt) VALUES($1, $2, $3, $4) RETURNING id`,
		strings.ToLower(userData.Email), GetHash(userData.Password + salt), userData.Nickname, salt,
		).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (Data UserRepository) DoesUserExist(authData models.Auth) (uint64, error) {
	var id uint64
	var password, salt string

	rows, err := Data.userDB.Query(`SELECT id, password, salt FROM users WHERE email=$1`, authData.Email)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	// Пользователя с таким email нет в базе
	if !rows.Next() {
		return 0, nil
	}
	if err := rows.Scan(&id, &password, &salt); err != nil {
		return 0, err
	}
	// Не совпадает пароль
	if GetHash(authData.Password + salt) != password {
		return 0, nil
	}

	return id, nil
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

func (Data UserRepository) GetSettings(userID int) (*models.SettingsGet, error) {
	var settings models.SettingsGet

	rows, err := Data.userDB.Query(`SELECT email, avatar, nickname FROM users WHERE id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, err
	}

	var avatarNULL sql.NullString

	if err = rows.Scan(&settings.Email, &avatarNULL, &settings.Nickname); err != nil {
		return nil, err
	}
	if !avatarNULL.Valid {
		settings.Avatar = os.Getenv("AVATAR_STORAGE") + AvatarDefaultPath + ".webp"
	} else {
		settings.Avatar = os.Getenv("AVATAR_STORAGE") + avatarNULL.String + "_500px.webp"
	}

	return &settings, nil
}

func (Data UserRepository) CheckPasswordByUserID(userID int, oldPassword string) (bool, error) {
	var password, salt string

	rows, err := Data.userDB.Query(`SELECT password, salt FROM users WHERE id=$1`, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() {
		return false, nil
	}

	if err := rows.Scan(&password, &salt); err != nil {
		return false, err
	}
	// Не совпадает пароль
	if GetHash(oldPassword + salt) != password {
		return false, nil
	}

	return true, nil
}

func (Data UserRepository) UpdateEmail(userID int, email string) error {
	err := Data.userDB.QueryRow(`UPDATE users SET email=$1 WHERE id=$2`, strings.ToLower(email), userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) UpdateNickname(userID int, nickname string) error {
	err := Data.userDB.QueryRow(`UPDATE users SET nickname=$1 WHERE id=$2`, nickname, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) UpdatePassword(userID int, password string, customSalt ...string) error {
	var salt string

	if len(customSalt) != 0 {
		salt = customSalt[0]
	} else {
		salt = GetRandomString(SaltLength)
	}

	err := Data.userDB.QueryRow(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`, GetHash(password + salt), salt, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) UpdateAvatar(userID int, fileName string) error {
	log.Println(fileName)
	err := Data.userDB.QueryRow(`UPDATE users SET avatar=$1 WHERE id=$2`, fileName, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) GetAvatarFilename(userID int) (string, error) {
	var filename string

	rows, err := Data.userDB.Query(`SELECT avatar FROM users WHERE id=$1`, userID)
	if err != nil {
		return "", err
	}

	if !rows.Next() {
		return "", nil
	}

	var avatarNULL sql.NullString
	if err := rows.Scan(&avatarNULL); err != nil {
		return "", err
	}
	if !avatarNULL.Valid {
		filename = "placeholder"
	} else {
		filename = avatarNULL.String
	}

	return filename, nil
}

func (File FileSystem) CreateImage(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	reader := io.Reader(f)
	src, err := imgconv.Decode(reader)
	if err != nil {
		return "", err
	}

	fileName := uuid.NewString()

	avatarLarge := imgconv.Resize(src, imgconv.ResizeOption{Width: AvatarWidthBig, Height: AvatarWidthBig})
	out, err := os.Create(os.Getenv("AVATAR_STORAGE") + fileName + "_500px.webp")
	if err != nil {
		return "", err
	}
	writer := io.Writer(out)
	err = webp.Encode(writer, avatarLarge, &webp.Options{Quality: 85})
	if err != nil {
		return "", err
	}

	avatarSmall := imgconv.Resize(src, imgconv.ResizeOption{Width: AvatarWidthLittle, Height: AvatarWidthLittle})
	out, err = os.Create(os.Getenv("AVATAR_STORAGE") + fileName + "_150px.webp")
	if err != nil {
		return "", err
	}
	writer = io.Writer(out)
	err = webp.Encode(writer, avatarSmall, &webp.Options{Quality: 85})
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (File FileSystem) DeleteImage(filename string) error {
	// 1) Проверяем, что файл существует
	doesFileExist := true
	if _, err := os.Stat(filename + "_150px.webp"); os.IsNotExist(err){
		doesFileExist = false
	}

	// 2) Удаляем файл со старой аватаркой
	if filename != "placeholder" && doesFileExist {
		err := os.Remove(filename + "_150px.webp")
		if err != nil {
			return err
		}
		err = os.Remove(filename + "_500px.webp")
		if err != nil {
			return err
		}
	}

	return nil
}

func (r RedisStore) StoreSession(userID uint64, customSessionToken ...string) (string, error) {
	var sessionToken string
	if len(customSessionToken) != 0 {
		sessionToken = customSessionToken[0]
	} else {
		sessionToken = GetRandomString(SessionTokenLength)
	}
	err := r.redisConnection.Set(ctx, sessionToken, userID, time.Hour).Err()
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}

func (r RedisStore) GetSessionUserId(session string) (int, error) {
	res, err := r.redisConnection.Get(ctx, session).Result()
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(res)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (r RedisStore) DeleteSession(cookieValue string) {
	r.redisConnection.Del(ctx, cookieValue)
}
