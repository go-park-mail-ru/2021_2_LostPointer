package errors

import "github.com/pkg/errors"

var (
	ErrWrongCredentials = errors.New("Wrong credentials")
	ErrUserNotFound     = errors.New("User not found in db")
)
