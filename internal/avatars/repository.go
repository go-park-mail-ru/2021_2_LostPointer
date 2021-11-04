package avatars

import "mime/multipart"

//go:generate moq -out ../mock/avatar_repository_mock.go -pkg mock . AvatarRepository:MockAvatarRepository
type AvatarRepository interface {
	CreateImage(*multipart.FileHeader) (string, error)
	DeleteImage(string) error
}
