package repository

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/utils/constants"
	"2021_2_LostPointer/internal/utils/hash"
	sanitizer "2021_2_LostPointer/internal/utils/sanitize"
	"database/sql"
	"os"
	"strings"
)

type UserRepository struct {
	userDB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{userDB: db}
}

func (Data UserRepository) CreateUser(userData *models.User) (int, error) {
	var id int

	salt := hash.GetRandomString(constants.SaltLength)
	sanitizedData := sanitizer.SanitizeUserData(*userData)
	err := Data.userDB.QueryRow(
		`INSERT INTO users(email, password, nickname, salt, avatar) VALUES($1, $2, $3, $4, $5) RETURNING id`,
		strings.ToLower(sanitizedData.Email), hash.GetHash(sanitizedData.Password+salt), sanitizedData.Nickname, salt, constants.AvatarDefaultFileName,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (Data UserRepository) DoesUserExist(authData *models.Auth) (int, error) {
	var (
		id int
		password, salt string
	)

	rows, err := Data.userDB.Query(`SELECT id, password, salt FROM users WHERE email=$1`, authData.Email)
	if err != nil {
		return 0, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if !rows.Next() {
		return 0, nil
	}
	if err := rows.Scan(&id, &password, &salt); err != nil {
		return 0, err
	}
	if hash.GetHash(authData.Password+salt) != password {
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
	if hash.GetHash(oldPassword+salt) != password {
		return false, nil
	}

	return true, nil
}

func (Data UserRepository) UpdateEmail(userID int, email string) error {
	sanitizedEmail := sanitizer.SanitizeEmail(email)
	err := Data.userDB.QueryRow(`UPDATE users SET email=$1 WHERE id=$2`, strings.ToLower(sanitizedEmail), userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) UpdateNickname(userID int, nickname string) error {
	sanitizedNickname := sanitizer.SanitizeNickname(nickname)
	err := Data.userDB.QueryRow(`UPDATE users SET nickname=$1 WHERE id=$2`, sanitizedNickname, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (Data UserRepository) UpdatePassword(userID int, password string) error {
	salt := hash.GetRandomString(constants.SaltLength)

	err := Data.userDB.QueryRow(`UPDATE users SET password=$1, salt=$2 WHERE id=$3`, hash.GetHash(password+salt), salt, userID).Err()
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

//func sanitizeUserData(userData models.User) models.User {
//	var sanitizedData models.User
//
//	sanitizedData.Nickname = sanitize.HTML(userData.Nickname)
//	sanitizedData.Email = sanitize.HTML(userData.Email)
//	sanitizedData.Password = userData.Password
//
//	return sanitizedData
//}
//
//func sanitizeEmail(email string) string {
//	return sanitize.HTML(email)
//}
//
//func sanitizeNickname(nickname string) string {
//	return sanitize.HTML(nickname)
//}
