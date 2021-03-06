package authorization

import (
	"2021_2_LostPointer/internal/microservices/authorization/proto"
)

//go:generate moq -out ./mock/authorization_repo_mock.go -pkg mock . AuthStorage:MockAuthStorage
type AuthStorage interface {
	CreateSession(int64, string) error
	GetUserByCookie(string) (int64, error)
	DeleteSession(string) error
	GetUserByPassword(*proto.AuthData) (int64, error)
	CreateUser(*proto.RegisterData) (int64, error)
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	GetAvatar(int64) (string, error)
}
