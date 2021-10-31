package repository

import (
	"errors"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSessionRepository_CreateSession(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repository := NewSessionRepository(db)

	const (
		cookieValue = "awa"
		id          = 12345
	)
	tests := []struct {
		name        string
		mock        func()
		expected    error
		expectedErr bool
	}{
		{
			name: "successful session creation",
			mock: func() {
				mock.ExpectSet(cookieValue, id, time.Hour).SetVal("")
			},
		},
		{
			name: "failed session creation",
			mock: func() {
				mock.ExpectSet(cookieValue, id, time.Hour).SetErr(errors.New("awa"))
			},
			expectedErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			err := repository.CreateSession(id, cookieValue)
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

	const (
		cookieValue = "awa"
		id = 12345
		idStr = "12345"
	)
	tests := []struct {
		name        string
		mock        func()
		expected    error
		expectedErr bool
	}{
		{
			name: "successful getting user",
			mock: func() {
				mock.ExpectGet(cookieValue).SetVal(idStr)
			},
		},
		{
			name: "failed getting user",
			mock: func() {
				mock.ExpectGet(cookieValue).SetErr(errors.New("awa"))
			},
			expectedErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			res, err := repository.GetUserIdByCookie(cookieValue)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, id, res)
				assert.NoError(t, err)
			}
		})
	}
}

func TestSessionRepository_DeleteSession(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repository := NewSessionRepository(db)

	const cookieValue = "awa"
	tests := []struct {
		name        string
		mock        func()
		expected    error
		expectedErr bool
	}{
		{
			name: "successful session deletion",
			mock: func() {
				mock.ExpectDel(cookieValue).SetVal(1)
			},
		},
		{
			name: "failed session deletion",
			mock: func() {
				mock.ExpectDel(cookieValue).SetErr(errors.New("awa"))
			},
			expectedErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock()
			err := repository.DeleteSession(cookieValue)
			if testCase.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
