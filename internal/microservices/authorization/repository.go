package authorization

import "2021_2_LostPointer/internal/models"

type AuthStorage interface {
	CreateSession(int64, string) error
	GetUserByCookie(string) (int64, error)
	GetUserByPassword(*models.AuthData) (int64, error)
	CreateUser(*models.RegisterData) (int64, error)

	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	GetAvatar(int64) (string, error)
}
