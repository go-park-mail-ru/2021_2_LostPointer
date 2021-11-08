package image

import (
	"github.com/oliamb/cutter"
	"io"
	"mime/multipart"
	"os"

	"github.com/chai2010/webp"
	uuid "github.com/satori/go.uuid"
	"github.com/sunshineplan/imgconv"

	"2021_2_LostPointer/internal/constants"
)

type Service interface {
	CreateImage(*multipart.FileHeader) (string, error)
	DeleteImage(string) error
}

type AvatarsService struct{}

func NewAvatarsService() AvatarsService {
	return AvatarsService{}
}

func (service *AvatarsService) CreateImage(file *multipart.FileHeader) (string, error) {
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

	fileName := uuid.NewV4().String()

	src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.BigAvatarHeight})
	avatarLarge, err := cutter.Crop(src, cutter.Config{Width: constants.BigAvatarHeight, Height: constants.BigAvatarHeight, Mode: cutter.Centered})
	if err != nil {
		return "", err
	}
	out, err := os.Create(os.Getenv("USERS_FULL_PREFIX") + fileName + constants.BigAvatarPostfix)
	if err != nil {
		return "", err
	}
	writer := io.Writer(out)
	err = webp.Encode(writer, avatarLarge, &webp.Options{Quality: 85})
	if err != nil {
		return "", err
	}

	src = imgconv.Resize(src, imgconv.ResizeOption{Height: constants.LittleAvatarHeight})
	avatarSmall, err := cutter.Crop(src, cutter.Config{Width: constants.LittleAvatarHeight, Height: constants.LittleAvatarHeight, Mode: cutter.Centered})
	if err != nil {
		return "", err
	}
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

func (service *AvatarsService) DeleteImage(filename string) error {
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
