package models

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
)

type User struct {
	Login string `json:"id"`
	Password string `json:"name"`
}

// Password hashing
func getHash(str string) string {
	hasher := sha1.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

// IsUserExists - check if user exists in db
func IsUserExists(db *sql.DB, user *User) (bool, error) {
	hashedPassword := getHash(user.Password)

	rows, err := db.Query(
		"SELECT EXISTS(SELECT id FROM users WHERE username = $1 AND password = $2)",
		user.Login, hashedPassword)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var flag bool
	for rows.Next() {
		if err := rows.Scan(&flag); err != nil {
			return false, err
		}
	}

	return flag, nil
}
