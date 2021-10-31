package usecase

import (
	session "2021_2_LostPointer/internal/microservices/authorization/delivery"
	"2021_2_LostPointer/internal/mock"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthorizationUseCase_GetUserBySession(t *testing.T) {
	type response struct {
		userID *session.UserID
		error error
	}
	tests := []struct {
		name 		string
		mockDB 		*mock.MockUserRepository
		mockSession *mock.MockSessionRepository
		input 		*session.SessionData
		expected 	*response
		expectedErr bool
	}{
		{
			name: "Successfully returned user id",
			mockDB: &mock.MockUserRepository{},
			mockSession: &mock.MockSessionRepository{
				GetUserIdByCookieFunc: func(string) (int, error) {
					return 1, nil
				},
			},
			input: &session.SessionData{
				Cookies: "cookie_value",
			},
			expected: &response{
				userID: &session.UserID{
					UserID: 1,
				},
				error: nil,
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			r := NewAuthorizationUseCase(testCase.mockDB, testCase.mockSession)

			got, customError := r.GetUserBySession(context.Background(), testCase.input)
			if testCase.expectedErr {
				assert.NotNil(t, customError)
			} else {
				assert.Nil(t, customError)
				assert.Equal(t, testCase.expected, got)
			}
		})
	}
}
