package models

import (
	"2021_2_LostPointer/utils"
	"database/sql"
)

const Salt = "abEdfg" // Соль, которая добавляется к паролю

type User struct {
	ID int
	Username string
	Password string
}

// UserExistsLogin - используется обработчиком LoginUserHandler. Проверяет, что пользователь,
// который пытается авторизоваться есть в базе данных.
func UserExistsLogin(db *sql.DB, user User) (bool, error) {
	rows, err := db.Query(`SELECT id, username, password FROM users
			WHERE username=$1`, user.Username)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if !rows.Next() { // Из базы пришел пустой запрос => пользователя в базе данных нет
		return false, nil
	}
	var u User
	if err := rows.Scan(&u.ID, &u.Username, &u.Password); err != nil {
		return false, err
	}

	if utils.GetHash(user.Password + Salt) != u.Password { // Пароли не совпадают
		return false, nil
	}

	return true, nil
}

// IsUserUnique - используется обработчиком SignUpHandler. Проверяет что пользователь с указанным
// при регистрации username уникален.
func IsUserUnique(db *sql.DB, user User) (bool, error) {
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
func CreateUser(db *sql.DB, user User) error {
	_, err := db.Exec(`INSERT INTO users(username, password)
			VALUES($1, $2)`, user.Username, utils.GetHash(user.Password + Salt))
	if err != nil {
		return err
	}
	return nil
}
