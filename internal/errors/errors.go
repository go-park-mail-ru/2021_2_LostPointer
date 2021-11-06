package errors

import "github.com/pkg/errors"

var (
	ErrWrongCredentials = errors.New("Wrong credentials")
	ErrTypeAssertion    = errors.New("Assertion not succeeded")
)
