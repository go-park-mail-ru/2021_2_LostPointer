package images

import (
	"2021_2_LostPointer/internal/utils/constants"
	"github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/sunshineplan/imgconv"
	"io"
	"mime/multipart"
	"os"
)

//go:generate moq -out ../../mock/avatar_repository_mock.go -pkg mock . AvatarRepositoryIFace:MockAvatarRepositoryIFace
type AvatarRepositoryIFace interface {
	CreateImage(*multipart.FileHeader) (string, error)
	DeleteImage(string) error
}

type AvatarRepository struct{}

func NewAvatarRepository() AvatarRepository {
	return AvatarRepository{}
}

func (imageR AvatarRepository) CreateImage(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(f multipart.File) {
		_ = f.Close()
	}(f)
	reader := io.Reader(f)
	src, err := imgconv.Decode(reader)
	if err != nil {
		return "", err
	}

	fileName := uuid.NewString()

	avatarLarge := imgconv.Resize(src, imgconv.ResizeOption{Height: constants.BigAvatarHeight})
	out, err := os.Create(os.Getenv("USERS_FULL_PREFIX") + fileName + constants.BigAvatarPostfix)
	if err != nil {
		return "", err
	}
	writer := io.Writer(out)
	err = webp.Encode(writer, avatarLarge, &webp.Options{Quality: 85})
	if err != nil {
		return "", err
	}

	avatarSmall := imgconv.Resize(src, imgconv.ResizeOption{Height: constants.LittleAvatarHeight})
	out, err = os.Create(os.Getenv("USERS_FULL_PREFIX") + fileName + constants.LittleAvatarPostfix)
	if err != nil {
		return "", err
	}
	writer = io.Writer(out)
	err = webp.Encode(writer, avatarSmall, &webp.Options{Quality: 85})
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (imageR AvatarRepository) DeleteImage(filename string) error {
	// 1) Проверяем, что файл существует
	doesFileExist := true
	if _, err := os.Stat(filename + constants.LittleAvatarPostfix); os.IsNotExist(err) {
		doesFileExist = false
	}

	// 2) Удаляем файл со старой аватаркой
	if filename != constants.AvatarDefaultFileName && doesFileExist {
		err := os.Remove(filename + constants.LittleAvatarPostfix)
		if err != nil {
			return err
		}
		err = os.Remove(filename + constants.BigAvatarPostfix)
		if err != nil {
			return err
		}
	}

	return nil
}
