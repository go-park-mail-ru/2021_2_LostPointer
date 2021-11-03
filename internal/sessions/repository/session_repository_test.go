package repository

import (
	"2021_2_LostPointer/internal/constants"
	"errors"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repository := NewSessionRepository(db)

	type inputStruct struct {
		id 			int
		cookieValue string
	}
	tests := []struct {
		name        string
		mock        func()
		input 		inputStruct
		expected    error
		expectedErr bool
	}{
		{
			name: "Successfully stored session",
			mock: func() {
				mock.ExpectSet("cookie_value", 1, constants.CookieLifetime).SetVal("")
			},
			input: inputStruct{
				id: 1,
				cookieValue: "cookie_value",
			},
			expectedErr: false,
		},
		{
			name: "Error occurred in redis.Set",
			mock: func() {
				mock.ExpectSet("cookie_value", 1, constants.CookieLifetime).SetErr(errors.New("error"))
			},
			input: inputStruct{
				id: 1,
				cookieValue: "cookie_value",
			},
			expectedErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			err := repository.CreateSession(testCase.input.id, testCase.input.cookieValue)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionRepository_GetUserIdByCookie(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repository := NewSessionRepository(db)

	tests := []struct {
		name        string
		mock        func()
		input 		string
		expected   	int
		expectedErr bool
	}{
		{
			name: "Successfully returned user id by cookie value",
			mock: func() {
				mock.ExpectGet("cookie_value").SetVal("1")
			},
			input: "cookie_value",
			expected: 1,
		},
		{
			name: "Error occurred in redis.Get",
			mock: func() {
				mock.ExpectGet("cookie_value").SetErr(errors.New("error"))
			},
			input: "cookie_value",
			expectedErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			res, err := repository.GetUserIdByCookie(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, testCase.expected, res)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionRepository_DeleteSession(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repository := NewSessionRepository(db)

	tests := []struct {
		name        string
		mock        func()
		expected    error
		input 		string
		expectedErr bool
	}{
		{
			name: "Successfully deleted session",
			mock: func() {
				mock.ExpectDel("cookie_value").SetVal(1)
			},
			input: "cookie_value",
		},
		{
			name: "Error occurred in redis.Del",
			mock: func() {
				mock.ExpectDel("cookie_value").SetErr(errors.New("error"))
			},
			input: "cookie_value",
			expectedErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			err := repository.DeleteSession(testCase.input)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
