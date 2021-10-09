package utils

import (
	"2021_2_LostPointer/models"
	"database/sql"
	"log"
)

const SaltLength = 5

func UserExistsLogin(db *sql.DB, user models.User) (uint64, error) {
	rows, err := db.Query(`SELECT id, email, password, salt FROM users
			WHERE email=$1`, user.Email)
	if err != nil {
		log.Fatalln(err)
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() { // Из базы пришел пустой запрос => пользователя в базе данных нет
		return 0, nil
	}
	var u models.User
	if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Salt); err != nil {
		return 0, err
	}

	if GetHash(user.Password + u.Salt) != u.Password { // Пароли не совпадают
		return 0, nil
	}

	return u.ID, nil
}

func IsUserEmailUnique(db *sql.DB, email string) (bool, error) {
	rows, err := db.Query(`SELECT id FROM users WHERE email=$1`, email)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

func IsUserNicknameUnique(db *sql.DB, nickname string) (bool, error) {
	rows, err := db.Query(`SELECT id FROM users WHERE nickname=$1`, nickname)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return false, nil
	}
	return true, nil
}

func CreateUser(db *sql.DB, user models.User, customSalt ...string) (uint64, error) {
	var lastID uint64 = 0
	var salt string
	if len(customSalt) != 0 {
		salt = customSalt[0]
	} else {
		salt = GetRandomString(SaltLength)
	}
	err := db.QueryRow(`INSERT INTO users(email, password, salt, nickname)
			VALUES($1, $2, $3, $4) RETURNING id`,
			user.Email,
			GetHash(user.Password + salt),
			salt,
			user.Nickname).Scan(&lastID)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return lastID, nil
}
