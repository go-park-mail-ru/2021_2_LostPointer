package repository

import (
	"2021_2_LostPointer/internal/microservices/profile/proto"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/kennygrant/sanitize"
	"golang.org/x/crypto/bcrypt"

	"2021_2_LostPointer/internal/constants"
	customErrors "2021_2_LostPointer/internal/errors"
	"2021_2_LostPointer/pkg/utils"
)

type UserSettingsStorage struct {
	db *sql.DB
}

func NewUserSettingsStorage(db *sql.DB) *UserSettingsStorage {
	return &UserSettingsStorage{db: db}
}

func (storage *UserSettingsStorage) GetSettings(userID int64) (*proto.UserSettings, error) {
	query := `SELECT email, avatar, nickname FROM users WHERE id=$1`

	rows, err := storage.db.Query(query, userID)
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
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		return nil, customErrors.ErrUserNotFound
	}

	var avatar string
	settings := &proto.UserSettings{}
	if err = rows.Scan(&settings.Email, &avatar, &settings.Nickname); err != nil {
		return nil, err
	}
	settings.BigAvatar = os.Getenv("USERS_ROOT_PREFIX") + avatar + constants.UserAvatarExtension500px
	settings.SmallAvatar = os.Getenv("USERS_ROOT_PREFIX") + avatar + constants.UserAvatarExtension150px

	return settings, nil
}

func (storage *UserSettingsStorage) UpdateEmail(userID int64, email string) error {
	query := `UPDATE users SET email=$1 WHERE id=$2`

	sanitizedEmail := sanitize.HTML(email)
	_, err := storage.db.Exec(query, strings.ToLower(sanitizedEmail), userID)
	if err != nil {
		return err
	}
	return nil
}

func (storage *UserSettingsStorage) UpdateNickname(userID int64, nickname string) error {
	query := `UPDATE users SET nickname=$1 WHERE id=$2`

	sanitizedNickname := sanitize.HTML(nickname)
	_, err := storage.db.Exec(query, sanitizedNickname, userID)
	if err != nil {
		return err
	}
	return nil
}

func (storage *UserSettingsStorage) UpdateAvatar(userID int64, fileName string) error {
	query := `UPDATE users SET avatar=$1 WHERE id=$2`

	_, err := storage.db.Exec(query, fileName, userID)
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
	_, err = storage.db.Exec(query, hashedPassword, salt, userID)
	if err != nil {
		return err
	}
	return nil
}

//nolint:rowserrcheck
func (storage *UserSettingsStorage) IsEmailUnique(email string) (bool, error) {
	query := `SELECT id FROM users WHERE lower(email)=$1`

	rows, err := storage.db.Query(query, strings.ToLower(email))
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

//nolint:rowserrcheck
func (storage *UserSettingsStorage) IsNicknameUnique(nickname string) (bool, error) {
	query := `SELECT id FROM users WHERE lower(nickname)=$1`

	rows, err := storage.db.Query(query, strings.ToLower(nickname))
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

	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()
	if !rows.Next() {
		err = rows.Err()
		if err != nil {
			return false, err
		}
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
