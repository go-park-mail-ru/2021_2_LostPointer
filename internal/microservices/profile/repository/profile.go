package repository

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/kennygrant/sanitize"
	"golang.org/x/crypto/bcrypt"

	"2021_2_LostPointer/internal/constants"
	customErrors "2021_2_LostPointer/internal/errors"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/pkg/utils"
)

type UserSettingsStorage struct {
	db *sql.DB
}

func NewUserSettingsStorage(db *sql.DB) UserSettingsStorage {
	return UserSettingsStorage{db: db}
}

func (storage *UserSettingsStorage) GetSettings(userID int64) (*models.UserSettings, error) {
	query := `SELECT email, avatar, nickname FROM users WHERE id=$1`

	rows, err := storage.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	if !rows.Next() {
		return nil, customErrors.ErrUserNotFound
	}

	var (
		avatar   string
		settings models.UserSettings
	)
	if err = rows.Scan(&settings.Email, &avatar, &settings.Nickname); err != nil {
		return nil, err
	}
	settings.BigAvatar = os.Getenv("USERS_ROOT_PREFIX") + avatar + constants.BigAvatarPostfix
	settings.SmallAvatar = os.Getenv("USERS_ROOT_PREFIX") + avatar + constants.LittleAvatarPostfix

	return &settings, nil
}

func (storage *UserSettingsStorage) UpdateEmail(userID int64, email string) error {
	query := `UPDATE users SET email=$1 WHERE id=$2`

	sanitizedEmail := SanitizeEmail(email)
	err := storage.db.QueryRow(query, strings.ToLower(sanitizedEmail), userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (storage *UserSettingsStorage) UpdateNickname(userID int64, nickname string) error {
	query := `UPDATE users SET nickname=$1 WHERE id=$2`

	sanitizedNickname := SanitizeNickname(nickname)
	err := storage.db.QueryRow(query, sanitizedNickname, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (storage *UserSettingsStorage) UpdateAvatar(userID int64, fileName string) error {
	query := `UPDATE users SET avatar=$1 WHERE id=$2`

	err := storage.db.QueryRow(query, fileName, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (storage *UserSettingsStorage) UpdatePassword(userID int64, password string) error {
	query := `UPDATE users SET password=$1, salt=$2 WHERE id=$3`

	salt, err := utils.GetRandomString(constants.SaltLength)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = storage.db.QueryRow(query, hashedPassword, salt, userID).Err()
	if err != nil {
		return err
	}
	return nil
}

func (storage *UserSettingsStorage) IsEmailUnique(email string) (bool, error) {
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

func (storage *UserSettingsStorage) IsNicknameUnique(nickname string) (bool, error) {
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

func (storage *UserSettingsStorage) CheckPasswordByUserID(userID int64, oldPassword string) (bool, error) {
	query := `SELECT password, salt FROM users WHERE id=$1`

	rows, err := storage.db.Query(query, userID)
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
	if !rows.Next() {
		return false, nil
	}

	var password, salt string
	if err = rows.Scan(&password, &salt); err != nil {
		return false, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(password), []byte(oldPassword+salt)); err != nil {
		return false, customErrors.ErrWrongCredentials
	}

	return true, nil
}

func SanitizeEmail(email string) string {
	return sanitize.HTML(email)
}

func SanitizeNickname(nickname string) string {
	return sanitize.HTML(nickname)
}
