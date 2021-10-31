package repository

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/utils/constants"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/kennygrant/sanitize"
	"github.com/sunshineplan/imgconv"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type UserRepository struct {
	userDB *sql.DB
}

type FileSystem struct{}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{userDB: db}
}

func NewFileSystem() FileSystem {
	return FileSystem{}
}

func GetRandomString(l int) string {
	validCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = validCharacters[RandInt(0, len(validCharacters)-1)]
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

func (Data UserRepository) CreateUser(userData *models.User) (int, error) {
	var id int

	salt := GetRandomString(constants.SaltLength)
	sanitizedData := sanitizeUserData(*userData)
	err := Data.userDB.QueryRow(
		`INSERT INTO users(email, password, nickname, salt, avatar) VALUES($1, $2, $3, $4, $5) RETURNING id`,
		strings.ToLower(sanitizedData.Email), GetHash(sanitizedData.Password+salt), sanitizedData.Nickname, salt, constants.AvatarDefaultFileName,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (Data UserRepository) DoesUserExist(authData *models.Auth) (int, error) {
	var id int
	var password, salt string

	rows, err := Data.userDB.Query(`SELECT id, password, salt FROM users WHERE email=$1`, authData.Email)
	if err != nil {
		return 0, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	// Пользователя с таким email нет в базе
	if !rows.Next() {
		return 0, nil
	}
	if err := rows.Scan(&id, &password, &salt); err != nil {
		return 0, err
	}
	// Не совпадает пароль
	if GetHash(authData.Password+salt) != password {
		return 0, nil
	}

	return id, nil
}

func (Data UserRepository) IsEmailUnique(email string) (bool, error) {
	rows, err := Data.userDB.Query(`SELECT id FROM users WHERE lower(email)=$1`, strings.ToLower(email))
	if err != nil {
		return false, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
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
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
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
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if !rows.Next() {
		return nil, err
	}

	var avatarFilename string
	if err = rows.Scan(&settings.Email, &avatarFilename, &settings.Nickname); err != nil {
		return nil, err
	}
	settings.BigAvatar = os.Getenv("ROOT_PATH_PREFIX") + avatarFilename + constants.BigAvatarPostfix
	settings.SmallAvatar = os.Getenv("ROOT_PATH_PREFIX") + avatarFilename + constants.LittleAvatarPostfix

	return &settings, nil
}

func (Data UserRepository) CheckPasswordByUserID(userID int, oldPassword string) (bool, error) {
	var password, salt string

	rows, err := Data.userDB.Query(`SELECT password, salt FROM users WHERE id=$1`, userID)
	if err != nil {
		return false, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if !rows.Next() {
		return false, nil
	}

	if err := rows.Scan(&password, &salt); err != nil {
		return false, err
	}
	// Не совпадает пароль
	if GetHash(oldPassword+salt) != password {
		return false, nil
	}

	return true, nil
}

func (Data UserRepository) UpdateEmail(userID int, email string) error {
	sanitizedEmail := sanitizeEmail(email)
	err := Data.userDB.QueryRow(`UPDATE users SET email=$1 WHERE id=$2`, strings.ToLower(sanitizedEmail), userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) UpdateNickname(userID int, nickname string) error {
	sanitizedNickname := sanitizeNickname(nickname)
	err := Data.userDB.QueryRow(`UPDATE users SET nickname=$1 WHERE id=$2`, sanitizedNickname, userID).Err()
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
		salt = GetRandomString(constants.SaltLength)
	}

	err := Data.userDB.QueryRow(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`, GetHash(password+salt), salt, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) UpdateAvatar(userID int, fileName string) error {
	err := Data.userDB.QueryRow(`UPDATE users SET avatar=$1 WHERE id=$2`, fileName, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) GetAvatarFilename(userID int) (string, error) {
	rows, err := Data.userDB.Query(`SELECT avatar FROM users WHERE id=$1`, userID)
	if err != nil {
		return "", err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if !rows.Next() {
		return "", nil
	}

	var filename string
	if err := rows.Scan(&filename); err != nil {
		return "", err
	}

	return filename, nil
}

func (File FileSystem) CreateImage(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	reader := io.Reader(f)
	src, err := imgconv.Decode(reader)
	if err != nil {
		return "", err
	}

	fileName := uuid.NewString()

	avatarLarge := imgconv.Resize(src, imgconv.ResizeOption{Height: constants.BigAvatarHeight})
	out, err := os.Create(os.Getenv("FULL_PATH_PREFIX") + fileName + constants.BigAvatarPostfix)
	if err != nil {
		return "", err
	}
	writer := io.Writer(out)
	err = webp.Encode(writer, avatarLarge, &webp.Options{Quality: 85})
	if err != nil {
		return "", err
	}

	avatarSmall := imgconv.Resize(src, imgconv.ResizeOption{Height: constants.LittleAvatarHeight})
	out, err = os.Create(os.Getenv("FULL_PATH_PREFIX") + fileName + constants.LittleAvatarPostfix)
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
	if _, err := os.Stat(filename + constants.LittleAvatarPostfix); os.IsNotExist(err) {
		doesFileExist = false
	}

	// 2) Удаляем файл со старой аватаркой
	if filename != constants.AvatarDefaultFileName && doesFileExist {
		err := os.Remove(filename + constants.LittleAvatarPostfix)
		if err != nil {
			return err
		}
		err = os.Remove(filename + constants.BigAvatarPostfix)
		if err != nil {
			return err
		}
	}

	return nil
}

func sanitizeUserData(userData models.User) models.User {
	var sanitizedData models.User

	sanitizedData.Nickname = sanitize.HTML(userData.Nickname)
	sanitizedData.Email = sanitize.HTML(userData.Email)
	sanitizedData.Password = userData.Password

	return sanitizedData
}

func sanitizeEmail(email string) string {
	return sanitize.HTML(email)
}

func sanitizeNickname(nickname string) string {
	return sanitize.HTML(nickname)
}
