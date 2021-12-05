package image

import (
	"2021_2_LostPointer/internal/models"
	"bufio"
	"image"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/CapsLock-Studio/go-webpbin"
	"github.com/cenkalti/dominantcolor"
	"github.com/oliamb/cutter"
	uuid "github.com/satori/go.uuid"
	"github.com/sunshineplan/imgconv"
)

type Service interface {
	DeleteAvatar(string) error

	CreateImages(*multipart.FileHeader, string, map[int]string) (*models.ImageData, error)
}

type ImagesService struct{}

func NewImagesService() ImagesService {
	return ImagesService{}
}

//nolint:cyclop
func (service *ImagesService) CreateImages(fileHeader *multipart.FileHeader, path string, extensions map[int]string) (*models.ImageData, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer func(f multipart.File) {
		_ = f.Close()
	}(file)
	reader := bufio.NewReader(file)
	src, err := imgconv.Decode(reader)
	if err != nil {
		return nil, err
	}

	// Получаем размеры изображения
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	reader = bufio.NewReader(file)
	tmp, err := imgconv.DecodeConfig(reader)
	if err != nil {
		return nil, err
	}
	width := tmp.Width
	height := tmp.Height

	imageFilename := uuid.NewV4().String()
	var artworkColor string
	for size, extension := range extensions {
		log.Println(size, ": ", extension)
		var (
			newSrc image.Image
			out    *os.File
		)

		if height < width {
			newSrc = imgconv.Resize(src, imgconv.ResizeOption{Height: size})
		} else {
			newSrc = imgconv.Resize(src, imgconv.ResizeOption{Width: size})
		}

		img, err := cutter.Crop(newSrc, cutter.Config{Width: size, Height: size, Mode: cutter.Centered})
		if err != nil {
			return nil, err
		}

		out, err = os.Create(path + imageFilename + extension)
		if err != nil {
			return nil, err
		}
		writer := io.Writer(out)
		if err = webpbin.Encode(writer, img); err != nil {
			return nil, err
		}

		// Вычисление акцента
		if artworkColor == "" {
			artworkColor = dominantcolor.Hex(dominantcolor.Find(img))
		}
	}

	return &models.ImageData{
		Filename:     imageFilename,
		ArtworkColor: artworkColor,
	}, nil
}

func (service *ImagesService) DeleteImages(path string, filename string, extensions []string, defaultFilename string) error {
	for _, extension := range extensions {
		doesFileExist := true
		if _, err := os.Stat(path + filename + extension); os.IsNotExist(err) {
			doesFileExist = false
		}

		if filename != defaultFilename && doesFileExist {
			err := os.Remove(path + filename + extension)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
