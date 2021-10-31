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

	var cookieValue = "awa"
	var id = 12345
	tests := []struct {
		name        string
		mock        func()
		expected    error
		expectedErr bool
	}{
		{
			name: "successes",
			mock: func() {
				mock.ExpectSet(cookieValue, id, time.Hour).SetVal("")
			},
			expectedErr: false,
		},
		{
			name: "error",
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

	var cookieValue = "awa"
	var id = 12345
	var idStr = "12345"
	tests := []struct {
		name        string
		mock        func()
		expected    error
		expectedErr bool
	}{
		{
			name: "successes",
			mock: func() {
				mock.ExpectGet(cookieValue).SetVal(idStr)
			},
			expectedErr: false,
		},
		{
			name: "error",
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

	var cookieValue = "awa"
	tests := []struct {
		name        string
		mock        func()
		expected    error
		expectedErr bool
	}{
		{
			name: "successes",
			mock: func() {
				mock.ExpectDel(cookieValue).SetVal(1)
			},
			expectedErr: false,
		},
		{
			name: "error",
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
