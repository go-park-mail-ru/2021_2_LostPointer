package utils

import (
	"2021_2_LostPointer/models"
	"database/sql"
)

// UserExistsLogin - используется обработчиком LoginUserHandler. Проверяет, что пользователь,
// который пытается авторизоваться есть в базе данных.
func UserExistsLogin(db *sql.DB, user models.User) (uint64, error) {
	rows, err := db.Query(`SELECT id, username, password, salt FROM users
			WHERE username=$1`, user.Username)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if !rows.Next() { // Из базы пришел пустой запрос => пользователя в базе данных нет
		return 0, nil
	}
	var u models.User
	if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Salt); err != nil {
		return 0, err
	}

	if GetHash(user.Password + u.Salt) != u.Password { // Пароли не совпадают
		return 0, nil
	}

	return u.ID, nil
}

// IsUserUnique - используется обработчиком SignUpHandler. Проверяет что пользователь с указанным
// при регистрации username уникален.
func IsUserUnique(db *sql.DB, user models.User) (bool, error) {
	rows, err := db.Query(`SELECT * FROM users WHERE username=$1`, user.Username)
	if err != nil {
		return false, err
	}
	if rows.Next() { // Пользователь с таким username зарегистрирован
		return false, nil
	}
	return true, nil
}

// CreateUser - создаем пользователя в базе
func CreateUser(db *sql.DB, user models.User) (uint64, error) {
	var lastID uint64 = 0

	salt := GetRandomString(5)

	err := db.QueryRow(`INSERT INTO users(username, password, salt)
			VALUES($1, $2, $3) RETURNING id`,
			user.Username,
			GetHash(user.Password + salt),
			salt).Scan(&lastID)
	if err != nil {
		return 0, err
	}
	return lastID, nil
}